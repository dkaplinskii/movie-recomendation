package domain

import "time"

type RecommendationReadyEvent struct {
	EventID   string    `json:"event_id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
