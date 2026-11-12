package outbox

import (
	"context"
	"fmt"

	"studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/models"
)

const getOutboxEventQuery = `
SELECT event_id, user_id, film_id, review
FROM outbox_events
ORDER BY event_id
LIMIT $1
`

func (r *PostgresOutboxRepository) GetOutboxEvents(ctx context.Context, limit int) ([]models.OutboxEvent, error) {
	var events []models.OutboxEvent

	rows, err := r.pool.Query(ctx, getOutboxEventQuery, limit)
	if err != nil {
		return nil, fmt.Errorf("query outbox events: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		event := models.OutboxEvent{}
		if err := rows.Scan(&event.EventID, &event.UserID, &event.FilmID, &event.Review); err != nil {
			return nil, fmt.Errorf("scan outbox event: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate outbox events: %w", err)
	}

	return events, nil
}
