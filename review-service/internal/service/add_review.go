package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/models"
)

const (
	minReview = 1
	maxReview = 5
)

func (s *Service) AddReview(ctx context.Context, input models.AddReviewInput) error {
	if input.Review < minReview || input.Review > maxReview {
		return fmt.Errorf(
			"%w: review must be between %d and %d",
			ErrInvalidReview,
			minReview,
			maxReview,
		)
	}

	event := models.ReviewCreatedEvent{
		UserID: input.UserID,
		FilmID: input.FilmID,
		Review: input.Review,
	}

	if err := s.txManager.Do(
		ctx,
		func(tx pgx.Tx) error {
			if err := s.storage.UpsertReview(ctx, tx, input); err != nil {
				return fmt.Errorf("upsert review: %w", err)
			}

			if err := s.outboxRepo.CreateEvent(ctx, tx, event); err != nil {
				return fmt.Errorf("create outbox event: %w", err)
			}

			return nil
		},
	); err != nil {
		return fmt.Errorf("save review with outbox event: %w", err)
	}

	return nil
}
