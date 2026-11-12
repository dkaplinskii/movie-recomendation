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

type MovieFilter struct {
	Genres    []string `json:"genres"`
	Countries []string `json:"countries"`
	YearFrom  int      `json:"year_from"`
	YearTo    int      `json:"year_to"`
	Limit     int      `json:"limit"`
}
