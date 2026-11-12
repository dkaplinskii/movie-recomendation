package telegram

import (
	"context"
	"fmt"
	"telegram-bot-service/internal/domain"
	"time"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

// В будущем для бизнес-логики
type Service interface {
	RegisterUser(ctx context.Context, user domain.User) error
}

type Bot struct {
	bot     *tele.Bot
	log     *zap.SugaredLogger
	service Service // Бизнес-логика вместо репозитория
}

func NewBot(token string, timeout time.Duration, logger *zap.SugaredLogger, service Service) (*Bot, error) {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: timeout},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("failed to create new bot: %w", err)
	}

	tgBot := &Bot{
		bot:     b,
		log:     logger,
		service: service,
	}

	return tgBot, nil
}

func (b *Bot) Start() {
	b.bot.Start()
}

func (b *Bot) Stop() {
	b.bot.Stop()
}
