package algorithm

import (
	"recommendation-service/internal/domain"
	"sort"
)

type scoredString struct {
	Value string
	Score float64
}

const noYearFilter = 0

// Находит наиболее любимый год (с достаточной выборкой и оценкой выше порога) и возвращает диапазон ±yearRange вокруг него; если такого года нет — возвращает нулевой диапазон.
func preferredYearRange(yearScores map[int]float64, yearCounts map[int]int, minScore float64, yearRange, minSampleSize int) (int, int) {
	var (
		bestYear  int
		bestScore float64
		found     bool
	)
	for year, score := range yearScores {
		if year <= 0 || score < minScore {
			continue
		}
		if yearCounts[year] < minSampleSize {
			continue
		}
		if !found || score > bestScore || (score == bestScore && year > bestYear) {
			bestYear = year
			bestScore = score
			found = true
		}
	}
	if !found {
		return noYearFilter, noYearFilter
	}
	return bestYear - yearRange, bestYear + yearRange
}

// Отбирает значения с оценкой не ниже порога.
func topScores(scores map[string]float64, minScore float64, limit int) []string {
	items := make([]scoredString, 0, len(scores))

	for value, score := range scores {
		if score < minScore {
			continue
		}

		items = append(items, scoredString{
			Value: value,
			Score: score,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Score != items[j].Score {
			return items[i].Score > items[j].Score
		}
		return items[i].Value < items[j].Value
	})
	if len(items) > limit {
		items = items[:limit]
	}

	result := make([]string, 0, len(items))
	for _, item := range items {
		result = append(result, item.Value)
	}

	return result
}

// Превращает профиль пользователя в фильтр для movie-service: предпочитаемые жанры и страны, диапазон лет и размер пула кандидатов.
func BuildMovieFilter(profile Profile, cfg Config) domain.MovieFilter {
	yearFrom, yearTo := preferredYearRange(
		profile.YearScores,
		profile.YearCounts,
		cfg.MinPreferenceScore,
		cfg.PreferredYearRange,
		cfg.MinYearSampleSize,
	)

	return domain.MovieFilter{
		Genres: topScores(
			profile.GenreScores,
			cfg.MinPreferenceScore,
			cfg.MaxFilterValues,
		),
		Countries: topScores(
			profile.CountryScores,
			cfg.MinPreferenceScore,
			cfg.MaxFilterValues,
		),
		YearFrom: yearFrom,
		YearTo:   yearTo,
		Limit:    cfg.CandidatePoolSize,
	}
}
