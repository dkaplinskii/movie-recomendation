package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"studentgit.kata.academy/movie-recommendation-platform-20/movie-service/internal/domain"
)

const getMovieByIDQuery = `
	SELECT id, kinopoisk_id, title, description, genres, poster_url, year, countries, trailer_url
	FROM movies
	WHERE id = $1
`

func (m *Storage) GetMovieByID(ctx context.Context, id int64) (domain.Movie, error) {
	var movie domain.Movie

	err := m.pool.QueryRow(ctx, getMovieByIDQuery, id).Scan(
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
			return domain.Movie{}, fmt.Errorf("no rows for id: %d", id)
		}
		return domain.Movie{}, fmt.Errorf("failed to get movie: %w", err)
	}

	return movie, nil
}
