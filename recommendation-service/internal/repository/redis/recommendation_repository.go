package redis

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RecommendationRepository struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRecommendationRepository(client *redis.Client, ttl time.Duration) *RecommendationRepository {
	return &RecommendationRepository{
		client: client,
		ttl:    ttl,
	}
}

func recommendationKey(userID int64) string {
	return fmt.Sprintf("recommendations:user:%d", userID)
}
