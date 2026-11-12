package outbox

import (
	"context"
	"encoding/json"
	"fmt"

	"studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/models"
)

func (w *Worker) processBatch(ctx context.Context) error {
	events, err := w.repo.GetOutboxEvents(ctx, w.batchSize)
	if err != nil {
		return fmt.Errorf("get outbox events: %w", err)
	}

	for _, event := range events {
		reviewEvent := models.ReviewCreatedEvent{
			UserID: event.UserID,
			FilmID: event.FilmID,
			Review: event.Review,
		}

		payload, err := json.Marshal(reviewEvent)
		if err != nil {
			return fmt.Errorf("marshal outbox event %d: %w", event.EventID, err)
		}

		if err := w.publisher.Publish(payload); err != nil {
			return fmt.Errorf("publish outbox event %d: %w", event.EventID, err)
		}

		if err := w.repo.DeleteOutboxEvent(ctx, event.EventID); err != nil {
			return fmt.Errorf("delete outbox event %d: %w", event.EventID, err)
		}
	}

	return nil
}
