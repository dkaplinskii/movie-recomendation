package domain

type User struct {
	UserID         int   `json:"user_id"`
	TelegramID     int64 `json:"telegram_id"`
	TelegramChatID int64 `json:"telegram_chat_id"`
}
