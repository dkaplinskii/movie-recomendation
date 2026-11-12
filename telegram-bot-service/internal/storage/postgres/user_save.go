package postgres

import (
	"context"
	"fmt"
	"telegram-bot-service/internal/domain"
)

const saveUserQuery = `
		INSERT INTO users (telegram_id, telegram_chat_id)
		VALUES($1, $2)
		ON CONFLICT (telegram_id) DO UPDATE 
		SET telegram_chat_id = EXCLUDED.telegram_chat_id
	`

func (u *UserRepository) SaveUser(ctx context.Context, user domain.User) error {
	_, err := u.pool.Exec(ctx, saveUserQuery, user.TelegramID, user.TelegramChatID)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}
