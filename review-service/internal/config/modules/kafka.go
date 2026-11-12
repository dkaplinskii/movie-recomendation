package modules

import "time"

type Kafka struct {
	Brokers              []string      `env:"BROKERS" envSeparator:"," envDefault:"kafka:9092"`
	ReviewEventsTopic    string        `env:"REVIEW_EVENTS_TOPIC" envDefault:"reviews.created"`
	ProducerRequiredAcks int           `env:"PRODUCER_REQUIRED_ACKS" envDefault:"1"`
	ProducerRetryMax     int           `env:"PRODUCER_RETRY_MAX" envDefault:"3"`
	ProducerTimeout      time.Duration `env:"PRODUCER_TIMEOUT" envDefault:"10s"`
}
