package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"recommendation-service/internal/domain"
	"recommendation-service/internal/service"

	"go.uber.org/zap"
)

type fakeReviewClient struct {
	reviews []domain.Review
	err     error
}

func (c *fakeReviewClient) GetLatestReviews(ctx context.Context, userID int64, limit int) ([]domain.Review, error) {
	return c.reviews, c.err
}

type fakeMovieClient struct {
	byIDs     []domain.Movie
	byFilter  []domain.Movie
	byIDsErr  error
	filterErr error
}

func (c *fakeMovieClient) GetMoviesByIDs(ctx context.Context, movieIDs []int64) ([]domain.Movie, error) {
	return c.byIDs, c.byIDsErr
}

func (c *fakeMovieClient) GetMoviesByFilter(ctx context.Context, filter domain.MovieFilter) ([]domain.Movie, error) {
	return c.byFilter, c.filterErr
}
func TestHandleReviewCreatedNewFeed(t *testing.T) {
	ctx := context.Background()
	repo := &fakeRecommendationRepo{
		getErr: domain.ErrRecommendationsNotFound,
	}
	producer := &fakeMessageProducer{}
	reviewClient, movieClient := feedClients()
	svc := service.NewService(repo, zap.NewNop().Sugar(), producer, "recommendations.ready", reviewClient, movieClient, testAlgorithmConfig())

	err := svc.HandleReviewCreated(ctx, domain.ReviewCreatedEvent{
		EventID: "event-1",
		UserID:  42,
		MovieID: 100,
		Rating:  9,
	})
	if err != nil {
		t.Fatalf("HandleReviewCreated() error = %v", err)
	}

	if repo.addCalls != 1 {
		t.Fatalf("repo.AddRecommendation calls = %d, want 1", repo.addCalls)
	}
	if repo.addedFeed.UserID != 42 {
		t.Fatalf("saved recommendation user id = %d, want 42", repo.addedFeed.UserID)
	}
	if len(repo.addedFeed.Movies) != 10 {
		t.Fatalf("saved movies count = %d, want 10", len(repo.addedFeed.Movies))
	}

	if producer.sendCalls != 1 {
		t.Fatalf("producer.Send calls = %d, want 1", producer.sendCalls)
	}
	if producer.topic != "recommendations.ready" {
		t.Fatalf("producer topic = %q, want %q", producer.topic, "recommendations.ready")
	}
	if producer.key != "42" {
		t.Fatalf("producer key = %q, want %q", producer.key, "42")
	}

	var event domain.RecommendationReadyEvent
	if err := json.Unmarshal(producer.payload, &event); err != nil {
		t.Fatalf("unmarshal producer payload: %v", err)
	}
	if event.UserID != 42 {
		t.Fatalf("notification user id = %d, want 42", event.UserID)
	}
}

func TestHandleReviewCreatedSameFeed(t *testing.T) {
	ctx := context.Background()
	repo := &fakeRecommendationRepo{
		feed: generatedFeed(42),
	}
	producer := &fakeMessageProducer{}
	reviewClient, movieClient := feedClients()
	svc := service.NewService(repo, zap.NewNop().Sugar(), producer, "recommendations.ready", reviewClient, movieClient, testAlgorithmConfig())
	err := svc.HandleReviewCreated(ctx, domain.ReviewCreatedEvent{
		EventID: "event-1",
		UserID:  42,
		MovieID: 100,
		Rating:  9,
	})
	if err != nil {
		t.Fatalf("HandleReviewCreated() error = %v", err)
	}

	if repo.addCalls != 0 {
		t.Fatalf("repo.AddRecommendation calls = %d, want 0", repo.addCalls)
	}
	if producer.sendCalls != 0 {
		t.Fatalf("producer.Send calls = %d, want 0", producer.sendCalls)
	}
}

func TestHandleReviewCreatedRepoError(t *testing.T) {
	ctx := context.Background()
	repo := &fakeRecommendationRepo{
		getErr: errors.New("redis is down"),
	}
	producer := &fakeMessageProducer{}
	reviewClient, movieClient := feedClients()
	svc := service.NewService(repo, zap.NewNop().Sugar(), producer, "recommendations.ready", reviewClient, movieClient, testAlgorithmConfig())
	err := svc.HandleReviewCreated(ctx, domain.ReviewCreatedEvent{
		EventID: "event-1",
		UserID:  42,
		MovieID: 100,
		Rating:  9,
	})
	if err == nil {
		t.Fatal("HandleReviewCreated() error = nil, want error")
	}

	if repo.addCalls != 0 {
		t.Fatalf("repo.AddRecommendation calls = %d, want 0", repo.addCalls)
	}
	if producer.sendCalls != 0 {
		t.Fatalf("producer.Send calls = %d, want 0", producer.sendCalls)
	}
}

type fakeRecommendationRepo struct {
	feed      domain.Recommendation
	getErr    error
	addErr    error
	addedFeed domain.Recommendation
	addCalls  int
}

func (r *fakeRecommendationRepo) GetRecommendation(ctx context.Context, userID int64) (domain.Recommendation, error) {
	if r.getErr != nil {
		return domain.Recommendation{}, r.getErr
	}

	return r.feed, nil
}

func (r *fakeRecommendationRepo) AddRecommendation(ctx context.Context, feed domain.Recommendation) error {
	r.addCalls++
	r.addedFeed = feed

	return r.addErr
}

type fakeMessageProducer struct {
	topic     string
	key       string
	payload   []byte
	sendCalls int
	sendErr   error
}

func (p *fakeMessageProducer) Send(topic string, key string, payload []byte) error {
	p.sendCalls++
	p.topic = topic
	p.key = key
	p.payload = payload

	return p.sendErr
}

func generatedFeed(userID int64) domain.Recommendation {
	movies := make([]domain.MovieRecommendation, 0, 10)
	for i := 1; i <= 10; i++ {
		movies = append(movies, domain.MovieRecommendation{
			ID: int64(i),
		})
	}

	return domain.Recommendation{
		UserID: userID,
		Movies: movies,
	}
}
