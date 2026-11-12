package algorithm

import "recommendation-service/internal/domain"

// Собирает уникальные ID оценённых фильмов.
func UniqueReviewedMovieIDs(reviews []domain.Review) ([]int64, map[int64]struct{}) {
	ids := make([]int64, 0, len(reviews))
	seen := make(map[int64]struct{}, len(reviews))
	for _, review := range reviews {
		if review.MovieID <= 0 {
			continue
		}
		if _, ok := seen[review.MovieID]; ok {
			continue
		}
		seen[review.MovieID] = struct{}{}
		ids = append(ids, review.MovieID)
	}
	return ids, seen
}

// Убирает из кандидатов фильмы, которые пользователь уже оценивал.
func ExcludeReviewedMovies(movies []domain.Movie, reviewed map[int64]struct{}) []domain.Movie {
	result := make([]domain.Movie, 0, len(movies))

	for _, movie := range movies {
		if _, ok := reviewed[movie.ID]; ok {
			continue
		}

		result = append(result, movie)
	}

	return result
}

// Конвертирует доменные фильмы в элементы ленты рекомендаций.
func MoviesToRecommendations(movies []domain.Movie) []domain.MovieRecommendation {
	recommendations := make([]domain.MovieRecommendation, 0, len(movies))

	for _, movie := range movies {
		recommendations = append(recommendations, domain.MovieRecommendation{
			ID:          movie.ID,
			Title:       movie.Title,
			Description: movie.Description,
			Genres:      movie.Genres,
			PosterURL:   movie.PosterURL,
			Year:        movie.Year,
			Country:     movie.Country,
			TrailerURL:  movie.TrailerURL,
		})
	}

	return recommendations
}

// Добирает в текущую ленту фильмы из пула до размера size, пропуская исключённые и уже добавленные.
func AppendUnique(current, pool []domain.Movie, excluded map[int64]struct{}, size int) []domain.Movie {
	selected := make(map[int64]struct{}, len(current))
	for _, m := range current {
		selected[m.ID] = struct{}{}
	}
	for _, m := range pool {
		if len(current) >= size {
			break
		}
		if _, ok := excluded[m.ID]; ok {
			continue
		}
		if _, ok := selected[m.ID]; ok {
			continue
		}
		current = append(current, m)
		selected[m.ID] = struct{}{}
	}

	return current
}
