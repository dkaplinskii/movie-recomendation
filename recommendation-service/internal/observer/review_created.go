package observer

import (
	"context"
	"encoding/json"
	"fmt"

	"recommendation-service/internal/domain"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type ReviewCreatedProcessor interface {
	HandleReviewCreated(ctx context.Context, event domain.ReviewCreatedEvent) error
}

type ReviewCreatedHandler struct {
	log       *zap.SugaredLogger
	processor ReviewCreatedProcessor
}

func NewReviewCreatedHandler(log *zap.SugaredLogger, processor ReviewCreatedProcessor) *ReviewCreatedHandler {
	return &ReviewCreatedHandler{
		log:       log,
		processor: processor,
	}
}

func (h *ReviewCreatedHandler) Handle(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var event domain.ReviewCreatedEvent

	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return domain.ErrInvalidReviewCreatedEvent
	}

	if event.UserID <= 0 {
		return domain.ErrInvalidReviewCreatedEvent
	}

	if err := h.processor.HandleReviewCreated(ctx, event); err != nil {
		return fmt.Errorf("handle review created event: %w", err)
	}

	return nil
}
