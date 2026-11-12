package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const getUserIDQuery = `
		SELECT id FROM users
		WHERE telegram_id = $1
	`

func (u *UserRepository) GetUserIDByTelegramID(ctx context.Context, tgID int64) (int, error) {
	var userID int

	err := u.pool.QueryRow(ctx, getUserIDQuery, tgID).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("no rows for tgID: %d", tgID)
		}
		return 0, fmt.Errorf("failed to get userID: %w", err)
	}

	return userID, nil
}
