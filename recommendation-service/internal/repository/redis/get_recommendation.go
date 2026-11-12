package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"recommendation-service/internal/domain"

	"github.com/redis/go-redis/v9"
)

func (r *RecommendationRepository) GetRecommendation(ctx context.Context, userID int64) (domain.Recommendation, error) {
	key := recommendationKey(userID)

	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return domain.Recommendation{}, domain.ErrRecommendationsNotFound
		}

		return domain.Recommendation{}, fmt.Errorf("redis get recommendation feed: %w", err)
	}

	var feed domain.Recommendation

	if err := json.Unmarshal(data, &feed); err != nil {
		return domain.Recommendation{}, fmt.Errorf("unmarshal recommendation feed: %w", err)
	}

	return feed, nil
}
