package service

import (
	"context"
	"fmt"
	"recommendation-service/internal/algorithm"
	"recommendation-service/internal/domain"
)

func (s *Service) GenerateRecommendation(ctx context.Context, userID int64) (domain.Recommendation, error) {
	if userID <= 0 {
		return domain.Recommendation{}, fmt.Errorf("user id must be positive")
	}

	if s.reviewClient == nil {
		return domain.Recommendation{}, fmt.Errorf("review client is not configured")
	}
	if s.movieClient == nil {
		return domain.Recommendation{}, fmt.Errorf("movie client is not configured")
	}

	// Берём последние оценки пользователя для построения профиля предпочтений.
	reviews, err := s.reviewClient.GetLatestReviews(ctx, userID, s.algorithmCfg.LatestReviewsLimit)
	if err != nil {
		return domain.Recommendation{}, fmt.Errorf("get latest reviews: %w", err)
	}
	userAverageRating := algorithm.AverageReviewRating(reviews)

	// У пользователя ещё нет оценок — отдаём случайную ленту.
	if len(reviews) == 0 {
		return s.generateFallbackRecommendation(ctx, userID, nil)
	}

	// Собираем уникальные ID оценённых фильмов.
	movieIDs, reviewedIDs := algorithm.UniqueReviewedMovieIDs(reviews)
	if len(movieIDs) == 0 {
		return s.generateFallbackRecommendation(ctx, userID, nil)
	}

	// Подтягиваем сами фильмы (жанры, страну, год) для профиля.
	reviewedMovies, err := s.movieClient.GetMoviesByIDs(ctx, movieIDs)
	if err != nil {
		return domain.Recommendation{}, fmt.Errorf("get movies by ids: %w", err)
	}

	//movie-сервис не вернул данные — построить профиль невозможно, отдаём случайную ленту.
	if len(reviewedMovies) == 0 {
		return s.generateFallbackRecommendation(ctx, userID, reviewedIDs)
	}

	// Строим профиль пользователя: средние оценки по жанрам, странам и годам.
	profile := algorithm.BuildUserProfile(reviews, reviewedMovies)

	// Превращаем профиль в фильтр: только предпочитаемые жанры/страны и интересующий диапазон лет.
	filter := algorithm.BuildMovieFilter(profile, s.algorithmCfg)

	// Берём пул кандидатов из movie-service по этому фильтру.
	candidateMovies, err := s.movieClient.GetMoviesByFilter(ctx, filter)
	if err != nil {
		return domain.Recommendation{}, fmt.Errorf("get movies by filter: %w", err)
	}

	// Исключаем уже оценённые фильмы.
	candidateMovies = algorithm.ExcludeReviewedMovies(candidateMovies, reviewedIDs)

	// После исключения кандидатов не осталось — отдаём случайную ленту вместо пустой.
	if len(candidateMovies) == 0 {
		return s.generateFallbackRecommendation(ctx, userID, reviewedIDs)
	}

	scoredMovies := algorithm.ScoreCandidateMovies(candidateMovies, profile, userAverageRating, s.algorithmCfg)

	// Сортируем по рейтингу и берём топ фильмов.
	topMovies := algorithm.TopScoredMovies(scoredMovies, s.algorithmCfg.RecommendationSize)

	// Если в ленте меньше 10 фильмов, то добираем фильмами, которые пользователь ещё не оценивал.
	if len(topMovies) < s.algorithmCfg.RecommendationSize {
		topMovies = s.topUpMovies(ctx, topMovies, reviewedIDs, s.algorithmCfg.RecommendationSize)
	}

	return domain.Recommendation{
		UserID: userID,
		Movies: algorithm.MoviesToRecommendations(topMovies),
	}, nil
}

func (s *Service) generateFallbackRecommendation(ctx context.Context, userID int64, reviewed map[int64]struct{}) (domain.Recommendation, error) {
	pool, err := s.defaultMovies(ctx)
	if err != nil {
		return domain.Recommendation{}, err
	}
	movies := algorithm.AppendUnique(nil, pool, reviewed, s.algorithmCfg.RecommendationSize)
	return domain.Recommendation{
		UserID: userID,
		Movies: algorithm.MoviesToRecommendations(movies),
	}, nil
}

func (s *Service) defaultMovies(ctx context.Context) ([]domain.Movie, error) {
	movies, err := s.movieClient.GetMoviesByFilter(ctx, domain.MovieFilter{
		Limit: s.algorithmCfg.CandidatePoolSize,
	})
	if err != nil {
		return nil, fmt.Errorf("get default movies: %w", err)
	}
	return movies, nil
}

// Добивает неполную ленту дефолтными фильмами из movie-service; при ошибке логирует и возвращает что есть.
func (s *Service) topUpMovies(ctx context.Context, current []domain.Movie, reviewed map[int64]struct{}, size int) []domain.Movie {
	pool, err := s.defaultMovies(ctx)
	if err != nil {
		s.log.Warnw("top-up failed", "error", err)
		return current
	}
	return algorithm.AppendUnique(current, pool, reviewed, size)
}
