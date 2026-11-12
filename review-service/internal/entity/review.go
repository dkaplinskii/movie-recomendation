package entity

import "time"

type Review struct {
	FilmID    int64     `json:"film_id"`
	Review    int32     `json:"review"`
	CreatedAt time.Time `json:"created_at"`
}
