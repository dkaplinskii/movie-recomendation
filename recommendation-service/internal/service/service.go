package service

import (
	"context"
	"recommendation-service/internal/algorithm"
	"recommendation-service/internal/domain"

	"go.uber.org/zap"
)

type RecommendationRepo interface {
	AddRecommendation(ctx context.Context, feed domain.Recommendation) error
	GetRecommendation(ctx context.Context, userID int64) (domain.Recommendation, error)
}

type MessageProducer interface {
	Send(topic string, key string, payload []byte) error
}

type RecommendationNotifier interface {
	NotifyRecommendationReady(userID int64) error
}

//go:generate mockgen -destination=../mocks/clients_mock.go -package=mocks recommendation-service/internal/service ReviewClient,MovieClient
type ReviewClient interface {
	GetLatestReviews(ctx context.Context, userID int64, limit int) ([]domain.Review, error)
}

type MovieClient interface {
	GetMoviesByIDs(ctx context.Context, movieIDs []int64) ([]domain.Movie, error)
	GetMoviesByFilter(ctx context.Context, filter domain.MovieFilter) ([]domain.Movie, error)
}

type Service struct {
	repo                      RecommendationRepo
	log                       *zap.SugaredLogger
	producer                  MessageProducer
	recommendationEventsTopic string
	reviewClient              ReviewClient
	movieClient               MovieClient
	algorithmCfg              algorithm.Config
}

func NewService(
	repo RecommendationRepo,
	log *zap.SugaredLogger,
	producer MessageProducer,
	recommendationEventsTopic string,
	reviewClient ReviewClient,
	movieClient MovieClient,
	algorithmCfg algorithm.Config,
) *Service {
	return &Service{
		repo:                      repo,
		log:                       log,
		producer:                  producer,
		recommendationEventsTopic: recommendationEventsTopic,
		reviewClient:              reviewClient,
		movieClient:               movieClient,
		algorithmCfg:              algorithmCfg,
	}
}
