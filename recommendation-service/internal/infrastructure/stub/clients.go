package stub

import (
	"context"

	"recommendation-service/internal/domain"
	"recommendation-service/internal/mocks"
	"recommendation-service/internal/service"

	"go.uber.org/mock/gomock"
)

type nopT struct{}

func (nopT) Errorf(string, ...any) {}
func (nopT) Fatalf(string, ...any) {}

func NewReviewClient() service.ReviewClient {
	ctrl := gomock.NewController(nopT{})
	mock := mocks.NewMockReviewClient(ctrl)

	mock.EXPECT().
		GetLatestReviews(gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes().
		DoAndReturn(func(_ context.Context, userID int64, _ int) ([]domain.Review, error) {
			return defaultReviews(userID), nil
		})

	return mock
}

func NewMovieClient() service.MovieClient {
	ctrl := gomock.NewController(nopT{})
	mock := mocks.NewMockMovieClient(ctrl)

	mock.EXPECT().
		GetMoviesByIDs(gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(defaultReviewedMovies(), nil)

	mock.EXPECT().
		GetMoviesByFilter(gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(defaultCandidateMovies(), nil)

	return mock
}
