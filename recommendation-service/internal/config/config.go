package config

import (
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	App                     AppConfig
	Server                  ServerConfig
	Redis                   RedisConfig
	Kafka                   KafkaConfig
	RecommendationCache     RecommendationCacheConfig
	RecommendationAlgorithm RecommendationAlgorithmConfig
}

type AppConfig struct {
	AppName     string `env:"APP_NAME" envDefault:"recommendation-service"`
	Environment string `env:"APP_ENV" envDefault:"local"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"info"`
}

type ServerConfig struct {
	GRPCPort        int           `env:"GRPC_PORT" envDefault:"50051"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"10s"`
}

type RedisConfig struct {
	RedisHost           string        `env:"REDIS_HOST" envDefault:"redis"`
	RedisPort           int           `env:"REDIS_PORT" envDefault:"6379"`
	RedisPassword       string        `env:"REDIS_PASSWORD"`
	RedisDB             int           `env:"REDIS_DB" envDefault:"0"`
	RedisStartupTimeout time.Duration `env:"REDIS_STARTUP_TIMEOUT" envDefault:"10s"`
}

type RecommendationCacheConfig struct {
	RecommendationCacheTTL time.Duration `env:"RECOMMENDATION_CACHE_TTL" envDefault:"24h"`
}

type KafkaConfig struct {
	Brokers                   []string `env:"KAFKA_BROKERS" envSeparator:"," envDefault:"kafka:9092"`
	GroupID                   string   `env:"KAFKA_GROUP_ID" envDefault:"recommendation-service"`
	ReviewEventsTopic         string   `env:"KAFKA_REVIEW_EVENTS_TOPIC" envDefault:"reviews.created"`
	RecommendationEventsTopic string   `env:"KAFKA_RECOMMENDATION_EVENTS_TOPIC" envDefault:"recommendations.ready"`

	ConsumerAutoCommitInterval time.Duration `env:"KAFKA_CONSUMER_AUTO_COMMIT_INTERVAL" envDefault:"1s"`

	ProducerRequiredAcks int           `env:"KAFKA_PRODUCER_REQUIRED_ACKS" envDefault:"1"`
	ProducerRetryMax     int           `env:"KAFKA_PRODUCER_RETRY_MAX" envDefault:"3"`
	ProducerTimeout      time.Duration `env:"KAFKA_PRODUCER_TIMEOUT" envDefault:"10s"`
}

type RecommendationAlgorithmConfig struct {
	LatestReviewsLimit int     `env:"RECOMMENDATION_LATEST_REVIEWS_LIMIT" envDefault:"100"`
	RecommendationSize int     `env:"RECOMMENDATION_SIZE" envDefault:"10"`
	MinPreferenceScore float64 `env:"RECOMMENDATION_MIN_PREFERENCE_SCORE" envDefault:"7.0"`
	MaxFilterValues    int     `env:"RECOMMENDATION_MAX_FILTER_VALUES" envDefault:"5"`
	PreferredYearRange int     `env:"RECOMMENDATION_PREFERRED_YEAR_RANGE" envDefault:"5"`
	GenreWeight        float64 `env:"RECOMMENDATION_GENRE_WEIGHT" envDefault:"0.4"`
	CountryWeight      float64 `env:"RECOMMENDATION_COUNTRY_WEIGHT" envDefault:"0.3"`
	YearWeight         float64 `env:"RECOMMENDATION_YEAR_WEIGHT" envDefault:"0.2"`
	RatingWeight       float64 `env:"RECOMMENDATION_RATING_WEIGHT" envDefault:"0.1"`
	MinYearSampleSize  int     `env:"RECOMMENDATION_MIN_YEAR_SAMPLE_SIZE" envDefault:"2"`
	CandidatePoolSize  int     `env:"RECOMMENDATION_CANDIDATE_POOL_SIZE" envDefault:"200"`
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("load env config: %w", err)
	}

	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("parse env config: %w", err)
	}

	return &cfg, nil
}
