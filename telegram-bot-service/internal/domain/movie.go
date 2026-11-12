package domain

type Movie struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Genres      []string `json:"genres"`
	PosterURL   string   `json:"poster_url"`
	Year        int      `json:"year"`
	Country     string   `json:"country"`
	TrailerURL  string   `json:"trailer_url"`
}

type Review struct {
	MovieID int64 `json:"movie_id"`
	UserID  int   `json:"user_id"`
	Rating  int   `json:"rating"`
}
