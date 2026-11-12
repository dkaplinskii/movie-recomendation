package modules

import "time"

type Outbox struct {
	BatchSize    int           `env:"BATCH_SIZE" envDefault:"100"`
	PollInterval time.Duration `env:"POLL_INTERVAL" envDefault:"1s"`
}
