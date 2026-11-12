package service

import (
	"context"
	"errors"
	"fmt"
	"recommendation-service/internal/domain"
)

func (s *Service) GetMovieRecommendation(ctx context.Context, userID int64) (domain.Recommendation, error) {
	recommendation, err := s.repo.GetRecommendation(ctx, userID)

	if err == nil {
		return recommendation, nil
	}

	if !errors.Is(err, domain.ErrRecommendationsNotFound) {
		return domain.Recommendation{}, fmt.Errorf("get recommendation: %w", err)
	}

	recommendation, err = s.GenerateRecommendation(ctx, userID)
	if err != nil {
		return domain.Recommendation{}, fmt.Errorf("generate recommendation: %w", err)
	}

	if err := s.repo.AddRecommendation(ctx, recommendation); err != nil {
		return domain.Recommendation{}, fmt.Errorf("save generated recommendation: %w", err)
	}

	return recommendation, nil
}
