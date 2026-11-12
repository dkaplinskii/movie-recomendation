package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/models"
)

const upsertReviewQuery = `
INSERT INTO reviews (user_id, film_id, review)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, film_id)
DO UPDATE SET review = EXCLUDED.review
`

func (r *PostgresRepository) UpsertReview(ctx context.Context, tx pgx.Tx, input models.AddReviewInput) error {
	if _, err := tx.Exec(
		ctx,
		upsertReviewQuery,
		input.UserID,
		input.FilmID,
		input.Review,
	); err != nil {
		return fmt.Errorf("upsert review: %w", err)
	}

	return nil
}
