package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

func (p *Publisher) Publish(payload []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(payload),
	}

	if _, _, err := p.producer.SendMessage(msg); err != nil {
		return fmt.Errorf("send kafka message: %w", err)
	}

	return nil
}
