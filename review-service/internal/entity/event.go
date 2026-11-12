package entity

type ReviewEvent struct {
	UserID int64 `json:"user_id"`
	FilmID int64 `json:"film_id"`
	Review int32 `json:"review"`
}
