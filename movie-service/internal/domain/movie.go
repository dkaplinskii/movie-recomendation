package domain

type Movie struct {
	ID          int64    `json:"id"`
	KinopoiskID int64    `json:"kinopoisk_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Genres      []string `json:"genres"`
	PosterURL   string   `json:"poster_url"`
	Year        int      `json:"year"`
	Countries   []string `json:"countries"`
	TrailerURL  string   `json:"trailer_url"`
}

type MovieFilter struct {
	Country string `json:"country"`
	Genre   string `json:"genre"`
	Year    int    `json:"year"`
	Limit   int    `json:"limit"`
}
