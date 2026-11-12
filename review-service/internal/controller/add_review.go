package controller

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/models"
	"studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/service"
	"studentgit.kata.academy/movie-recommendation-platform-20/review-service/pkg/reviewpb"
)

func (h *Controller) AddReview(ctx context.Context, req *reviewpb.AddReviewRequest) (*reviewpb.AddReviewResponse, error) {
	input := models.AddReviewInput{
		UserID: req.UserId,
		FilmID: req.FilmId,
		Review: req.Review,
	}

	err := h.svc.AddReview(ctx, input)
	if err != nil {
		if errors.Is(err, service.ErrInvalidReview) {
			return &reviewpb.AddReviewResponse{
				Success: false,
			}, status.Error(codes.InvalidArgument, err.Error())
		}

		return &reviewpb.AddReviewResponse{
			Success: false,
		}, status.Error(codes.Internal, "failed to add review")
	}

	return &reviewpb.AddReviewResponse{
		Success: true,
	}, nil
}
