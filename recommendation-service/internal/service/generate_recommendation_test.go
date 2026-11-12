package service_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"recommendation-service/internal/algorithm"
	"recommendation-service/internal/domain"
	"recommendation-service/internal/service"

	"go.uber.org/zap"
)

func newGenerateService(reviewClient service.ReviewClient, movieClient service.MovieClient) *service.Service {
	return service.NewService(
		&fakeRecommendationRepo{},
		zap.NewNop().Sugar(),
		&fakeMessageProducer{},
		"recommendations.ready",
		reviewClient,
		movieClient,
		testAlgorithmConfig(),
	)
}

func TestGenerateRecommendationInvalidUserID(t *testing.T) {
	svc := newGenerateService(&fakeReviewClient{}, &fakeMovieClient{})

	if _, err := svc.GenerateRecommendation(context.Background(), 0); err == nil {
		t.Fatal("GenerateRecommendation() error = nil, want error for non-positive user id")
	}
}

func TestGenerateRecommendationNilReviewClient(t *testing.T) {
	svc := newGenerateService(nil, &fakeMovieClient{})

	if _, err := svc.GenerateRecommendation(context.Background(), 42); err == nil {
		t.Fatal("GenerateRecommendation() error = nil, want error for nil review client")
	}
}

func TestGenerateRecommendationReviewClientError(t *testing.T) {
	reviewClient := &fakeReviewClient{err: errors.New("review service down")}
	svc := newGenerateService(reviewClient, &fakeMovieClient{})

	if _, err := svc.GenerateRecommendation(context.Background(), 42); err == nil {
		t.Fatal("GenerateRecommendation() error = nil, want error from review client")
	}
}

func TestGenerateRecommendationMovieFilterError(t *testing.T) {
	reviewClient := &fakeReviewClient{reviews: latestReviews()}
	movieClient := &fakeMovieClient{
		byIDs:     reviewedMovies(),
		filterErr: errors.New("movie service down"),
	}
	svc := newGenerateService(reviewClient, movieClient)

	if _, err := svc.GenerateRecommendation(context.Background(), 42); err == nil {
		t.Fatal("GenerateRecommendation() error = nil, want error from movie filter")
	}
}

func TestGenerateRecommendationNoReviewsFallback(t *testing.T) {
	reviewClient := &fakeReviewClient{reviews: nil}
	movieClient := &fakeMovieClient{byFilter: candidateMovies()}
	svc := newGenerateService(reviewClient, movieClient)

	feed, err := svc.GenerateRecommendation(context.Background(), 42)
	if err != nil {
		t.Fatalf("GenerateRecommendation() error = %v", err)
	}
	if feed.UserID != 42 {
		t.Fatalf("feed user id = %d, want 42", feed.UserID)
	}
	if len(feed.Movies) != 10 {
		t.Fatalf("feed movies = %d, want 10", len(feed.Movies))
	}
}

func TestGenerateRecommendationPathTopOrder(t *testing.T) {
	reviewClient := &fakeReviewClient{reviews: latestReviews()}
	movieClient := &fakeMovieClient{
		byIDs:    reviewedMovies(),
		byFilter: candidateMovies(),
	}
	svc := newGenerateService(reviewClient, movieClient)

	feed, err := svc.GenerateRecommendation(context.Background(), 42)
	if err != nil {
		t.Fatalf("GenerateRecommendation() error = %v", err)
	}
	if len(feed.Movies) != 10 {
		t.Fatalf("feed movies = %d, want 10", len(feed.Movies))
	}

	for i, movie := range feed.Movies {
		if movie.ID != int64(i+1) {
			t.Fatalf("feed.Movies[%d].ID = %d, want %d", i, movie.ID, i+1)
		}
	}
}

func feedClients() (*fakeReviewClient, *fakeMovieClient) {
	reviewClient := &fakeReviewClient{
		reviews: []domain.Review{
			{UserID: 42, MovieID: 100, Rating: 9},
		},
	}
	candidates := make([]domain.Movie, 0, 10)
	for i := 1; i <= 10; i++ {
		candidates = append(candidates, domain.Movie{
			ID:      int64(i),
			Genres:  []string{"drama"},
			Country: "US",
			Year:    2000,
		})
	}
	movieClient := &fakeMovieClient{
		byIDs: []domain.Movie{
			{ID: 100, Genres: []string{"drama"}, Country: "US", Year: 2000},
		},
		byFilter: candidates,
	}
	return reviewClient, movieClient
}
func testAlgorithmConfig() algorithm.Config {
	return algorithm.Config{
		LatestReviewsLimit: 100,
		RecommendationSize: 10,
		MinPreferenceScore: 0,
		MaxFilterValues:    10,
		PreferredYearRange: 5,
		GenreWeight:        0.4,
		CountryWeight:      0.3,
		YearWeight:         0.2,
		RatingWeight:       0.1,
	}
}

func latestReviews() []domain.Review {
	reviews := make([]domain.Review, 0, 10)
	for k := 1; k <= 10; k++ {
		reviews = append(reviews, domain.Review{
			UserID:  42,
			MovieID: int64(1000 + k),
			Rating:  float64(11 - k),
		})
	}
	return reviews
}

func reviewedMovies() []domain.Movie {
	movies := make([]domain.Movie, 0, 10)
	for k := 1; k <= 10; k++ {
		movies = append(movies, domain.Movie{
			ID:      int64(1000 + k),
			Genres:  []string{fmt.Sprintf("g%d", k)},
			Country: "US",
			Year:    2000,
		})
	}
	return movies
}

func candidateMovies() []domain.Movie {
	movies := make([]domain.Movie, 0, 10)
	for i := 1; i <= 10; i++ {
		movies = append(movies, domain.Movie{
			ID:      int64(i),
			Genres:  []string{fmt.Sprintf("g%d", i)},
			Country: "US",
			Year:    2000,
		})
	}
	return movies
}
