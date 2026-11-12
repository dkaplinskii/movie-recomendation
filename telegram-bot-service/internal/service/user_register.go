package service

import (
	"context"
	"telegram-bot-service/internal/domain"
)

// Сделал пока простенький метод, чтоб не было пустого интерфейса в bot.go
func (s *BotService) RegisterUser(ctx context.Context, user domain.User) error {
	return s.repo.SaveUser(ctx, user)
}
