package domain

import "time"

type ReviewCreatedEvent struct {
	EventID   string    `json:"event_id"`
	UserID    int64     `json:"user_id"`
	MovieID   int64     `json:"movie_id"`
	Rating    int       `json:"rating"`
	CreatedAt time.Time `json:"created_at"`
}
