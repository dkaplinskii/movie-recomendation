package domain

type MovieRecommendation struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Genres      []string `json:"genres"`
	PosterURL   string   `json:"poster_url"`
	Year        int      `json:"year"`
	Country     string   `json:"country"`
	TrailerURL  string   `json:"trailer_url"`
}

type Recommendation struct {
	UserID int64                 `json:"user_id"`
	Movies []MovieRecommendation `json:"movies"`
}
