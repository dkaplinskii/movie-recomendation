package ratelimiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type IPRateLimiter struct {
	client *redis.Client
	limit  int
	burst  int
}

func NewIPRateLimiter(limit, burst int, redisAddr string) *IPRateLimiter {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	return &IPRateLimiter{
		client: rdb,
		limit:  limit,
		burst:  burst,
	}
}

// Allow проверяет, можно ли пропустить запрос с данного IP.
// Возвращает true, если запрос разрешён, false — если лимит превышен.
// При ошибке Redis возвращает false и ошибку.
func (rl *IPRateLimiter) Allow(ctx context.Context, ip string) (bool, error) {
	key := fmt.Sprintf("rate:%s", ip)

	pipe := rl.client.TxPipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, time.Minute)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, fmt.Errorf("redis error: %w", err)
	}

	count := incr.Val()
	allowed := int(count) <= rl.limit+rl.burst

	if !allowed {
		return false, nil
	}
	return true, nil
}
