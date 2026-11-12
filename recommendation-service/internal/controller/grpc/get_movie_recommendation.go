package grpc

import (
	"context"
	"recommendation-service/internal/domain"

	recommendationv1 "recommendation-service/pkg/recommendation/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RecommendationService interface {
	GetMovieRecommendation(ctx context.Context, userID int64) (domain.Recommendation, error)
}

func (c *Controller) GetMovieRecommendation(ctx context.Context, req *recommendationv1.GetMovieRecommendationRequest) (*recommendationv1.GetMovieRecommendationResponse, error) {
	recommendation, err := c.service.GetMovieRecommendation(ctx, req.GetUserId())
	if err != nil {
		c.log.Errorw("failed to get movie recommendations",
			"user_id", req.GetUserId(),
			"error", err,
		)
		return nil, status.Error(codes.Internal, "failed to get movie recommendations")
	}

	movies := make([]*recommendationv1.MovieRecommendation, 0, len(recommendation.Movies))
	for _, m := range recommendation.Movies {
		movies = append(movies, &recommendationv1.MovieRecommendation{
			Id:          m.ID,
			Title:       m.Title,
			Description: m.Description,
			Genres:      m.Genres,
			PosterUrl:   m.PosterURL,
			Year:        int32(m.Year),
			Country:     m.Country,
			TrailerUrl:  m.TrailerURL,
		})
	}

	return &recommendationv1.GetMovieRecommendationResponse{
		Movies: movies,
	}, nil
}
