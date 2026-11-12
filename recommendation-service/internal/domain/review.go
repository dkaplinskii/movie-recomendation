package domain

type Review struct {
	UserID  int64   `json:"user_id"`
	MovieID int64   `json:"movie_id"`
	Rating  float64 `json:"rating"`
}
