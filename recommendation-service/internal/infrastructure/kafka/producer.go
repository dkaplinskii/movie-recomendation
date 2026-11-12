package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"

	"recommendation-service/internal/config"

	"recommendation-service/internal/logger"
)

func InitKafkaProducer(ctx context.Context, appName string, cfg *config.KafkaConfig) (sarama.SyncProducer, error) {
	log := logger.FromContext(ctx)

	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	saramaCfg.Producer.Return.Errors = true
	saramaCfg.Producer.RequiredAcks = sarama.RequiredAcks(cfg.ProducerRequiredAcks)
	saramaCfg.Producer.Retry.Max = cfg.ProducerRetryMax
	saramaCfg.Producer.Timeout = cfg.ProducerTimeout
	saramaCfg.ClientID = appName

	producer, err := sarama.NewSyncProducer(cfg.Brokers, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("create kafka sync producer: %w", err)
	}

	log.Infow("kafka producer created",
		"brokers", cfg.Brokers,
		"topic", cfg.RecommendationEventsTopic,
	)

	return producer, nil
}
