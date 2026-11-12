package postgres

import (
	"context"
	"fmt"

	"studentgit.kata.academy/movie-recommendation-platform-20/movie-service/internal/domain"
)

const saveMovieQuery = `
	INSERT INTO movies (kinopoisk_id, title, description, genres, poster_url, year, countries, trailer_url)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8)
	ON CONFLICT(kinopoisk_id) DO UPDATE
	SET title = EXCLUDED.title,
		description = EXCLUDED.description,
		poster_url = EXCLUDED.poster_url,
		trailer_url = EXCLUDED.trailer_url
`

func (m *Storage) SaveMovie(ctx context.Context, movie domain.Movie) error {
	_, err := m.pool.Exec(ctx, saveMovieQuery,
		movie.KinopoiskID,
		movie.Title,
		movie.Description,
		movie.Genres,
		movie.PosterURL,
		movie.Year,
		movie.Countries,
		movie.TrailerURL,
	)
	if err != nil {
		return fmt.Errorf("failed to save movie: %w", err)
	}

	return nil
}
