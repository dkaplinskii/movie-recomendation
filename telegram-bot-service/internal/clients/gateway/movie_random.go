package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"telegram-bot-service/internal/domain"
)

func (c *Client) GetRandomMovie(ctx context.Context, userID int) (*domain.Movie, error) {
	endpoint, err := url.JoinPath(c.baseURL, "movies", "random")
	if err != nil {
		return nil, fmt.Errorf("failed to build url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	q.Add("user_id", strconv.Itoa(userID))
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, domain.ErrNotFound
	default:
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var movie domain.Movie

	err = json.NewDecoder(resp.Body).Decode(&movie)
	if err != nil {
		return nil, fmt.Errorf("failed to decode json: %w", err)
	}

	return &movie, nil
}
