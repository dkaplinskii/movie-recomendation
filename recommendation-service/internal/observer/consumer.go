package observer

import (
	"context"
	"errors"

	"recommendation-service/internal/domain"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type MessageHandler interface {
	Handle(ctx context.Context, msg *sarama.ConsumerMessage) error
}

type Consumer struct {
	log     *zap.SugaredLogger
	handler MessageHandler
}

func NewConsumer(log *zap.SugaredLogger, handler MessageHandler) *Consumer {
	return &Consumer{
		log:     log,
		handler: handler,
	}
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if err := c.handler.Handle(session.Context(), msg); err != nil {
			c.log.Errorw("failed to handle kafka message",
				"topic", msg.Topic,
				"partition", msg.Partition,
				"offset", msg.Offset,
				"error", err,
			)

			if errors.Is(err, domain.ErrInvalidReviewCreatedEvent) {
				session.MarkMessage(msg, "")
			}

			continue
		}

		session.MarkMessage(msg, "")
	}

	return nil
}
