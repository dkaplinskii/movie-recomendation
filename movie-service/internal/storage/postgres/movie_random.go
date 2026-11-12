package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"studentgit.kata.academy/movie-recommendation-platform-20/movie-service/internal/domain"
)

const getRandomMovieQuery = `
	SELECT id, kinopoisk_id, title, description, genres, poster_url, year, countries, trailer_url
	FROM movies
	ORDER BY RANDOM()
	LIMIT 1;
`

func (m *Storage) GetRandomMovie(ctx context.Context) (domain.Movie, error) {
	var movie domain.Movie

	err := m.pool.QueryRow(ctx, getRandomMovieQuery).Scan(
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
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Movie{}, fmt.Errorf("table is empty")
		}
		return domain.Movie{}, fmt.Errorf("failed to get random movie: %w", err)
	}

	return movie, nil
}
