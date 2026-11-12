package service

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/entity"
	"studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/models"
)

type storageStub struct {
	reviews []entity.Review
	err     error
}

func (s *storageStub) GetLatestReviews(_ context.Context, _ int64, _ int32) ([]entity.Review, error) {
	return s.reviews, s.err
}

func (s *storageStub) UpsertReview(_ context.Context, _ pgx.Tx, _ models.AddReviewInput) error {
	return nil
}

func TestGetLatestReviews(t *testing.T) {
	testReviews := []entity.Review{
		{
			FilmID:    1,
			Review:    5,
			CreatedAt: time.Now(),
		},
		{
			FilmID:    2,
			Review:    4,
			CreatedAt: time.Now(),
		},
	}

	storage := &storageStub{
		reviews: testReviews,
	}

	svc := NewService(storage, nil, nil)

	reviews, err := svc.GetLatestReviews(
		context.Background(),
		1,
		5,
	)

	require.NoError(t, err)
	require.Len(t, reviews, len(testReviews))
	require.Equal(t, testReviews[0].FilmID, reviews[0].FilmID)
	require.Equal(t, testReviews[0].Review, reviews[0].Review)
}
