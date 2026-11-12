package kinopoisk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"studentgit.kata.academy/movie-recommendation-platform-20/movie-service/internal/errors"
	movieMapper "studentgit.kata.academy/movie-recommendation-platform-20/movie-service/internal/mapper/movie"
)

func (c *KinopoiskClient) GetFilmByID(ctx context.Context, id int64) (movieMapper.FilmResponseDTO, error) {
	endpoint := fmt.Sprintf(c.baseURL+getFilmPath, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return movieMapper.FilmResponseDTO{}, fmt.Errorf("create film request: %w", err)
	}
	req.Header.Set("X-API-KEY", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return movieMapper.FilmResponseDTO{}, fmt.Errorf("do film request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return movieMapper.FilmResponseDTO{}, errors.ErrUnauthorized

		case http.StatusNotFound:
			return movieMapper.FilmResponseDTO{}, errors.ErrNotFound

		default:
			return movieMapper.FilmResponseDTO{}, fmt.Errorf("unexpected status: %d", resp.StatusCode)
		}
	}

	var film movieMapper.FilmResponseDTO

	if err := json.NewDecoder(resp.Body).Decode(&film); err != nil {
		return movieMapper.FilmResponseDTO{}, fmt.Errorf("decode film response: %w", err)
	}

	return film, nil
}
