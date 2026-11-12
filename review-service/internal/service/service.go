package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/entity"
	"studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/models"
)

type Storage interface {
	GetLatestReviews(ctx context.Context, userID int64, limit int32) ([]entity.Review, error)

	UpsertReview(ctx context.Context, tx pgx.Tx, input models.AddReviewInput) error
}

type OutboxRepository interface {
	CreateEvent(ctx context.Context, tx pgx.Tx, event models.ReviewCreatedEvent) error
}

type TransactionManager interface {
	Do(ctx context.Context, callback func(pgx.Tx) error) error
}

type Service struct {
	storage    Storage
	outboxRepo OutboxRepository
	txManager  TransactionManager
}

func NewService(storage Storage, outboxRepo OutboxRepository, txManager TransactionManager) *Service {
	return &Service{
		storage:    storage,
		outboxRepo: outboxRepo,
		txManager:  txManager,
	}
}

func (s *Service) GetLatestReviews(ctx context.Context, userID int64, limit int32) ([]entity.Review, error) {
	reviews, err := s.storage.GetLatestReviews(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("get latest reviews: %w", err)
	}

	return reviews, nil
}
