package postgres

import (
	"context"
	"fmt"

	"studentgit.kata.academy/movie-recommendation-platform-20/movie-service/internal/domain"
)

const getMoviesByFilterBaseQuery = `
	SELECT id, kinopoisk_id, title, description, genres, poster_url, year, countries, trailer_url
	FROM movies
	WHERE 1=1
	`

func (m *Storage) GetMoviesByFilter(ctx context.Context, filter domain.MovieFilter) ([]domain.Movie, error) {
	query := getMoviesByFilterBaseQuery
	args := []any{}
	counter := 0

	if filter.Country != "" {
		counter++
		query += fmt.Sprintf(" AND $%d = ANY(countries)", counter)
		args = append(args, filter.Country)
	}
	if filter.Genre != "" {
		counter++
		query += fmt.Sprintf(" AND $%d = ANY(genres)", counter)
		args = append(args, filter.Genre)
	}
	if filter.Year != 0 {
		counter++
		query += fmt.Sprintf(" AND year = $%d", counter)
		args = append(args, filter.Year)
	}

	counter++
	query += fmt.Sprintf(" ORDER BY RANDOM() LIMIT $%d;", counter)
	args = append(args, filter.Limit)

	rows, err := m.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get movies by filter: %w", err)
	}
	defer rows.Close()

	movies := make([]domain.Movie, 0, filter.Limit)

	for rows.Next() {
		movie := domain.Movie{}
		err := rows.Scan(
			&movie.ID,
			&movie.KinopoiskID,
			&movie.Title,
			&movie.Description,
			&movie.Genres,
			&movie.PosterURL,
			&movie.Year,
			&movie.Countries,
			&movie.TrailerURL,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan movies: %w", err)
		}

		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed iterations: %w", err)
	}

	return movies, nil
}
