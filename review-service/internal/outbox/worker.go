package outbox

import (
	"context"
	"time"

	"studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/config/modules"
	"studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/logger"
	"studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/models"
)

type Repository interface {
	GetOutboxEvents(ctx context.Context, limit int) ([]models.OutboxEvent, error)
	DeleteOutboxEvent(ctx context.Context, id int64) error
}

type Publisher interface {
	Publish(payload []byte) error
}

type Worker struct {
	repo         Repository
	publisher    Publisher
	batchSize    int
	pollInterval time.Duration
}

func NewWorker(repo Repository, publisher Publisher, outboxCfg modules.Outbox) *Worker {
	return &Worker{
		repo:         repo,
		publisher:    publisher,
		batchSize:    outboxCfg.BatchSize,
		pollInterval: outboxCfg.PollInterval,
	}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	lg := logger.FromContext(ctx)
	defer lg.Info("outbox worker stopped")

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			if err := w.processBatch(ctx); err != nil {
				if ctx.Err() != nil {
					return
				}

				lg.Errorw(
					"process outbox batch failed",
					"error", err,
				)
			}
		}
	}
}
