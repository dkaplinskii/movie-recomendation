package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type ConsumerRunner struct {
	log     *zap.SugaredLogger
	group   sarama.ConsumerGroup
	topic   string
	handler sarama.ConsumerGroupHandler
}

func NewConsumerRunner(log *zap.SugaredLogger, group sarama.ConsumerGroup, topic string, handler sarama.ConsumerGroupHandler) *ConsumerRunner {
	return &ConsumerRunner{
		log:     log,
		group:   group,
		topic:   topic,
		handler: handler,
	}
}

func (r *ConsumerRunner) Run(ctx context.Context) error {
	r.log.Infow("kafka consumer runner started", "topic", r.topic)

	for {
		if err := r.group.Consume(ctx, []string{r.topic}, r.handler); err != nil {
			return fmt.Errorf("consume kafka topic %q: %w", r.topic, err)
		}

		if ctx.Err() != nil {
			r.log.Info("kafka consumer runner stopped by context")
			return nil
		}

		r.log.Infow("kafka consumer session ended, restarting consume loop", "topic", r.topic)
	}
}

func (r *ConsumerRunner) Close() error {
	if r.group == nil {
		return nil
	}

	if err := r.group.Close(); err != nil {
		return fmt.Errorf("close kafka consumer group: %w", err)
	}

	r.log.Info("kafka consumer group closed")

	return nil
}
