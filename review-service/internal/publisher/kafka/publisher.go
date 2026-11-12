package kafka

import (
	"github.com/IBM/sarama"
)

type Publisher struct {
	producer sarama.SyncProducer
	topic    string
}

func NewPublisher(producer sarama.SyncProducer, topic string) *Publisher {
	return &Publisher{
		producer: producer,
		topic:    topic,
	}
}
