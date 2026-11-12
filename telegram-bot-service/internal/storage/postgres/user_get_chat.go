package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const getChatIDQuery = `
		SELECT telegram_chat_id FROM users
		WHERE id = $1
	`

func (u *UserRepository) GetChatIDByUserID(ctx context.Context, userID int) (int64, error) {
	var chatID int64

	err := u.pool.QueryRow(ctx, getChatIDQuery, userID).Scan(&chatID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("no rows for userID: %d", userID)
		}
		return 0, fmt.Errorf("failed to get chatID: %w", err)
	}

	return chatID, nil
}
