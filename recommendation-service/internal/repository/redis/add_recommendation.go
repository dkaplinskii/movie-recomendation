package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"recommendation-service/internal/domain"
)

func (r *RecommendationRepository) AddRecommendation(ctx context.Context, feed domain.Recommendation) error {
	data, err := json.Marshal(feed)
	if err != nil {
		return fmt.Errorf("marshal recommendation feed: %w", err)
	}

	key := recommendationKey(feed.UserID)

	if err := r.client.Set(ctx, key, data, r.ttl).Err(); err != nil {
		return fmt.Errorf("redis set recommendation feed: %w", err)
	}

	return nil
}
