package service

import (
	"context"
	"fmt"
	"recommendation-service/internal/domain"
)

func (s *Service) AddRecommendation(ctx context.Context, feed domain.Recommendation) error {
	if err := s.repo.AddRecommendation(ctx, feed); err != nil {
		return fmt.Errorf("add recommendation: %w", err)
	}

	return nil
}
