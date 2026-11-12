package controller

import (
	"studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/service"
	"studentgit.kata.academy/movie-recommendation-platform-20/review-service/pkg/reviewpb"
)

type Controller struct {
	reviewpb.UnimplementedReviewServiceServer
	svc *service.Service
}

func NewController(service *service.Service) *Controller {
	return &Controller{svc: service}
}
