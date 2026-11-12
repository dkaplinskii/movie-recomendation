package kinopoisk

import (
	"net/http"

	"studentgit.kata.academy/movie-recommendation-platform-20/movie-service/internal/config/modules"
)

const (
	getFilmPath   = "/api/v2.2/films/%d"
	getVideosPath = "/api/v2.2/films/%d/videos"
)

type KinopoiskClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewKinopoiskClient(cfg modules.Kinopoisk) (*KinopoiskClient, error) {
	client := &http.Client{Timeout: cfg.Timeout}

	return &KinopoiskClient{
		baseURL:    cfg.APIURL,
		apiKey:     cfg.APIKey,
		httpClient: client,
	}, nil
}
