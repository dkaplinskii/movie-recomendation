package service

import (
	"context"
	"telegram-bot-service/internal/domain"
)

type UserRepository interface {
	SaveUser(ctx context.Context, user domain.User) error
	GetUserIDByTelegramID(ctx context.Context, tgID int64) (int, error)
	GetChatIDByUserID(ctx context.Context, userID int) (int64, error)
}

type Client interface {
	GetRandomMovie(ctx context.Context, userID int) (*domain.Movie, error)
	GetRecommendations(ctx context.Context, userID int) ([]domain.Movie, error)
	AddReview(ctx context.Context, review domain.Review) error
}

type BotService struct {
	repo   UserRepository
	client Client
}

func NewBotService(repo UserRepository, client Client) *BotService {
	return &BotService{repo: repo, client: client}
}
