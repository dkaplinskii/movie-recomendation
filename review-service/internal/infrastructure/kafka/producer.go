package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
	"studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/config/modules"
)

func NewProducer(kafkaCfg modules.Kafka) (sarama.SyncProducer, error) {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	saramaCfg.Producer.Return.Errors = true
	saramaCfg.Producer.RequiredAcks = sarama.RequiredAcks(kafkaCfg.ProducerRequiredAcks)
	saramaCfg.Producer.Retry.Max = kafkaCfg.ProducerRetryMax
	saramaCfg.Producer.Timeout = kafkaCfg.ProducerTimeout

	producer, err := sarama.NewSyncProducer(kafkaCfg.Brokers, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("create kafka producer: %w", err)
	}
	return producer, nil
}
