package kafka

import (
	"context"
	"fmt"
	"recommendation-service/internal/config"

	"recommendation-service/internal/logger"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

func InitKafkaConsumer(ctx context.Context, appName string, cfg *config.KafkaConfig) (sarama.ConsumerGroup, error) {
	uuidForClient, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("generate kafka client uuid: %w", err)
	}

	clientID := fmt.Sprintf("%s-%s", appName, uuidForClient)

	log := logger.FromContext(ctx)

	cfgMap := sarama.NewConfig()
	cfgMap.ClientID = clientID
	cfgMap.Consumer.Offsets.AutoCommit.Enable = true
	cfgMap.Consumer.Offsets.AutoCommit.Interval = cfg.ConsumerAutoCommitInterval
	cfgMap.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, cfgMap)
	if err != nil {
		return nil, fmt.Errorf("create kafka consumer: %w", err)
	}

	log.Infof("Kafka consumer successfully created: client_id=%s group_id=%s", clientID, cfg.GroupID)

	return consumer, nil
}
