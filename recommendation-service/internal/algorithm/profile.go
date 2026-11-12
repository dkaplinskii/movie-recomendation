package algorithm

import "recommendation-service/internal/domain"

type Profile struct {
	GenreScores   map[string]float64
	CountryScores map[string]float64
	YearScores    map[int]float64
	YearCounts    map[int]int
}

type scoreAccumulator struct {
	sum   float64
	count int
}

func (a *scoreAccumulator) addScore(score float64) {
	a.sum += score
	a.count++
}

func (a scoreAccumulator) averageScore() float64 {
	if a.count == 0 {
		return 0
	}

	return a.sum / float64(a.count)
}

// Строит профиль предпочтений пользователя: средние оценки по жанрам, странам и годам на основе его оценок и данных о соответствующих фильмах.
func BuildUserProfile(reviews []domain.Review, movies []domain.Movie) Profile {
	reviewByMovieID := make(map[int64]domain.Review, len(reviews))

	for _, review := range reviews {
		reviewByMovieID[review.MovieID] = review
	}

	genreAcc := make(map[string]scoreAccumulator)
	countryAcc := make(map[string]scoreAccumulator)
	yearAcc := make(map[int]scoreAccumulator)

	for _, movie := range movies {
		review, ok := reviewByMovieID[movie.ID]
		if !ok {
			continue
		}

		for _, genre := range movie.Genres {
			acc := genreAcc[genre]
			acc.addScore(review.Rating)
			genreAcc[genre] = acc
		}

		if movie.Country != "" {
			acc := countryAcc[movie.Country]
			acc.addScore(review.Rating)
			countryAcc[movie.Country] = acc
		}

		if movie.Year > 0 {
			acc := yearAcc[movie.Year]
			acc.addScore(review.Rating)
			yearAcc[movie.Year] = acc
		}
	}

	return Profile{
		GenreScores:   averageStringScores(genreAcc),
		CountryScores: averageStringScores(countryAcc),
		YearScores:    averageIntScores(yearAcc),
		YearCounts:    countIntScores(yearAcc),
	}
}

// Извлекает количество оценок по каждому году (размер выборки).
func countIntScores(acc map[int]scoreAccumulator) map[int]int {
	result := make(map[int]int, len(acc))
	for key, value := range acc {
		result[key] = value.count
	}
	return result
}

// Сворачивает аккумуляторы со строковыми ключами (жанры, страны) в средние оценки.
func averageStringScores(acc map[string]scoreAccumulator) map[string]float64 {
	result := make(map[string]float64, len(acc))

	for key, value := range acc {
		result[key] = value.averageScore()
	}

	return result
}

// Сворачивает аккумуляторы с числовыми ключами (годы) в средние оценки.
func averageIntScores(acc map[int]scoreAccumulator) map[int]float64 {
	result := make(map[int]float64, len(acc))

	for key, value := range acc {
		result[key] = value.averageScore()
	}

	return result
}
