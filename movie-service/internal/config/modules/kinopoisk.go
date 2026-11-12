package modules

import "time"

type Kinopoisk struct {
	APIURL  string        `env:"API_URL,required"`
	APIKey  string        `env:"API_KEY,required"`
	Timeout time.Duration `env:"TIMEOUT" envDefault:"5s"`
}
