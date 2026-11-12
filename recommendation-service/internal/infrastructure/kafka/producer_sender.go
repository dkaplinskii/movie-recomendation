package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type Producer struct {
	log      *zap.SugaredLogger
	producer sarama.SyncProducer
}

func NewProducer(log *zap.SugaredLogger, producer sarama.SyncProducer) *Producer {
	return &Producer{
		log:      log,
		producer: producer,
	}
}

func (p *Producer) Send(topic string, key string, payload []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(payload),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("send kafka message: %w", err)
	}

	p.log.Infow("kafka message published",
		"topic", topic,
		"key", key,
		"partition", partition,
		"offset", offset,
	)

	return nil
}
