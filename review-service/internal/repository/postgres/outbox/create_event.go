package outbox

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/models"
)

const insertOutboxEventQuery = `
INSERT INTO outbox_events (user_id, film_id, review)
VALUES ($1, $2, $3)
`

func (r *PostgresOutboxRepository) CreateEvent(ctx context.Context, tx pgx.Tx, event models.ReviewCreatedEvent) error {
	if _, err := tx.Exec(
		ctx,
		insertOutboxEventQuery,
		event.UserID,
		event.FilmID,
		event.Review,
	); err != nil {
		return fmt.Errorf("insert outbox event: %w", err)
	}

	return nil
}
