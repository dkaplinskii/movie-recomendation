package service

import (
	"encoding/json"
	"fmt"
	"recommendation-service/internal/domain"
	"time"

	"github.com/google/uuid"
)

func (s *Service) NotifyRecommendationReady(userID int64) error {
	event := domain.RecommendationReadyEvent{
		EventID:   uuid.NewString(),
		UserID:    userID,
		CreatedAt: time.Now().UTC(),
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal recommendation ready event: %w", err)
	}

	key := fmt.Sprintf("%d", userID)
	if err := s.producer.Send(s.recommendationEventsTopic, key, payload); err != nil {
		return fmt.Errorf("send recommendation ready event: %w", err)
	}

	return nil
}
