package movie

type FilmResponseDTO struct {
	KinopoiskID  int64        `json:"kinopoiskId"`
	NameRu       string       `json:"nameRu"`
	NameEn       string       `json:"nameEn"`
	NameOriginal string       `json:"nameOriginal"`
	Description  string       `json:"description"`
	Countries    []CountryDTO `json:"countries"`
	Genres       []GenreDTO   `json:"genres"`
	Year         string       `json:"year"`
	PosterURL    string       `json:"posterUrl"`
}

type CountryDTO struct {
	Country string `json:"country"`
}

type GenreDTO struct {
	Genre string `json:"genre"`
}

type VideosResponseDTO struct {
	Items []VideoItemDTO `json:"items"`
}

type VideoItemDTO struct {
	URL  string `json:"url"`
	Name string `json:"name"`
	Site string `json:"site"`
}
