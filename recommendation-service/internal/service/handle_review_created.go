package service

import (
	"context"
	"errors"
	"fmt"
	"recommendation-service/internal/domain"
)

func (s *Service) HandleReviewCreated(ctx context.Context, event domain.ReviewCreatedEvent) error {
	userID := event.UserID

	// Получаем предыдущую ленту рекомендаций из Redis.
	oldFeed, err := s.repo.GetRecommendation(ctx, userID)
	hadPrevious := true
	if err != nil {
		if errors.Is(err, domain.ErrRecommendationsNotFound) {
			hadPrevious = false
			oldFeed = domain.Recommendation{}
		} else {
			return fmt.Errorf("get previous recommendation: %w", err)
		}
	}

	// Генерируем новую ленту рекомендаций на основе последних оценок пользователя.
	newFeed, err := s.GenerateRecommendation(ctx, userID)
	if err != nil {
		return fmt.Errorf("generate recommendation: %w", err)
	}

	// Если предыдущая лента была и новая лента полностью совпадает с ней, то нет смысла заново сохранять данные в Redis и отправлять в Kafka.
	if hadPrevious && feedsEqual(oldFeed, newFeed) {
		return nil
	}

	// Сохраняем новую ленту рекомендаций в Redis.
	if err := s.AddRecommendation(ctx, newFeed); err != nil {
		return fmt.Errorf("add recommendation: %w", err)
	}

	//Отправляем событие, что рекомендации были обновлены.
	if err := s.NotifyRecommendationReady(userID); err != nil {
		return fmt.Errorf("notify recommendations ready: %w", err)
	}

	s.log.Infow("review created event handled",
		"event_id", event.EventID,
		"user_id", event.UserID,
		"movie_id", event.MovieID,
		"rating", event.Rating,
	)

	return nil
}

func feedsEqual(a, b domain.Recommendation) bool {
	if len(a.Movies) != len(b.Movies) {
		return false
	}

	for i := range a.Movies {
		if a.Movies[i].ID != b.Movies[i].ID {
			return false
		}
	}

	return true
}
