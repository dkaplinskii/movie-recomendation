package kinopoisk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"studentgit.kata.academy/movie-recommendation-platform-20/movie-service/internal/errors"
	movieMapper "studentgit.kata.academy/movie-recommendation-platform-20/movie-service/internal/mapper/movie"
)

func (c *KinopoiskClient) GetVideosByID(ctx context.Context, id int64) (movieMapper.VideosResponseDTO, error) {
	endpoint := fmt.Sprintf(c.baseURL+getVideosPath, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return movieMapper.VideosResponseDTO{}, fmt.Errorf("create videos request: %w", err)
	}
	req.Header.Set("X-API-KEY", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return movieMapper.VideosResponseDTO{}, fmt.Errorf("do videos request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return movieMapper.VideosResponseDTO{}, errors.ErrUnauthorized

		case http.StatusNotFound:
			return movieMapper.VideosResponseDTO{}, nil

		default:
			return movieMapper.VideosResponseDTO{}, fmt.Errorf("unexpected status: %d", resp.StatusCode)
		}
	}

	var videos movieMapper.VideosResponseDTO

	if err := json.NewDecoder(resp.Body).Decode(&videos); err != nil {
		return movieMapper.VideosResponseDTO{}, fmt.Errorf("decode videos response: %w", err)
	}

	return videos, nil
}
