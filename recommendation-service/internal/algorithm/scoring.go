package algorithm

import (
	"recommendation-service/internal/domain"
	"sort"
)

type scoredMovie struct {
	Movie domain.Movie
	Score float64
}

// Считает среднюю оценку пользователя по всем его отзывам.
func AverageReviewRating(reviews []domain.Review) float64 {
	if len(reviews) == 0 {
		return 0
	}

	var sum float64
	for _, review := range reviews {
		sum += review.Rating
	}

	return sum / float64(len(reviews))
}

// Возвращает средний жанровый рейтинг фильма по тем его жанрам, что есть в профиле.
func averageGenreScore(genres []string, genreScores map[string]float64) float64 {
	if len(genres) == 0 {
		return 0
	}

	var (
		sum   float64
		count int
	)

	for _, genre := range genres {
		score, ok := genreScores[genre]
		if !ok {
			continue
		}
		sum += score
		count++
	}

	if count == 0 {
		return 0
	}

	return sum / float64(count)
}

// Усредняет рейтинги годов в окне ±yearRange вокруг года фильма.
func windowedYearScore(year int, yearScores map[int]float64, yearRange int) float64 {
	if year <= 0 || len(yearScores) == 0 {
		return 0
	}

	var (
		sum   float64
		count int
	)

	for reviewedYear, score := range yearScores {
		if reviewedYear < year-yearRange || reviewedYear > year+yearRange {
			continue
		}
		sum += score
		count++
	}

	if count == 0 {
		return 0
	}

	return sum / float64(count)
}

// Проставляет каждому фильму-кандидату итоговый рейтинг по профилю пользователя.
func ScoreCandidateMovies(movies []domain.Movie, profile Profile, userAverageRating float64, cfg Config) []scoredMovie {
	scored := make([]scoredMovie, 0, len(movies))
	for _, movie := range movies {
		scored = append(scored, scoredMovie{
			Movie: movie,
			Score: calculateMovieScore(movie, profile, userAverageRating, cfg),
		})
	}
	return scored
}

// Считает итоговый рейтинг фильма по формуле из ТЗ: w1·жанр + w2·страна + w3·год + w4·оценка пользователя.
func calculateMovieScore(movie domain.Movie, profile Profile, userAverageRating float64, cfg Config) float64 {
	genreScore := averageGenreScore(movie.Genres, profile.GenreScores)
	countryScore := profile.CountryScores[movie.Country]
	yearScore := windowedYearScore(movie.Year, profile.YearScores, cfg.PreferredYearRange)

	return cfg.GenreWeight*genreScore +
		cfg.CountryWeight*countryScore +
		cfg.YearWeight*yearScore +
		cfg.RatingWeight*userAverageRating
}

// Сортирует кандидатов по score по убыванию.
func TopScoredMovies(scored []scoredMovie, limit int) []domain.Movie {
	sort.Slice(scored, func(i, j int) bool {
		if scored[i].Score != scored[j].Score {
			return scored[i].Score > scored[j].Score
		}
		return scored[i].Movie.ID < scored[j].Movie.ID
	})

	if len(scored) > limit {
		scored = scored[:limit]
	}

	result := make([]domain.Movie, 0, len(scored))
	for _, item := range scored {
		result = append(result, item.Movie)
	}

	return result
}
