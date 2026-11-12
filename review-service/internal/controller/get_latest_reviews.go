package controller

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "studentgit.kata.academy/movie-recommendation-platform-20/review-service/pkg/reviewpb"
)

func (c *Controller) GetLatestReviews(ctx context.Context, req *pb.GetLatestReviewsRequest) (*pb.GetLatestReviewsResponse, error) {
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.Limit <= 0 {
		return nil, status.Error(codes.InvalidArgument, "limit must be positive")
	}

	// 1. Вызываем сервис (он возвращает []entity.Review)
	reviews, err := c.svc.GetLatestReviews(ctx, req.UserId, req.Limit)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// 2. Преобразуем entity.Review → pb.Review
	pbReviews := make([]*pb.Review, len(reviews))
	for i, r := range reviews {
		pbReviews[i] = &pb.Review{
			FilmId:    r.FilmID,
			Review:    r.Review,
			CreatedAt: r.CreatedAt.Format(time.RFC3339),
		}
	}

	return &pb.GetLatestReviewsResponse{Reviews: pbReviews}, nil
}
