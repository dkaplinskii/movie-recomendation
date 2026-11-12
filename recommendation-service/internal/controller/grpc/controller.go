package grpc

import (
	recommendationv1 "recommendation-service/pkg/recommendation/v1"

	"go.uber.org/zap"
)

type Controller struct {
	recommendationv1.UnimplementedRecommendationServiceServer
	service RecommendationService
	log     *zap.SugaredLogger
}

func NewController(service RecommendationService, log *zap.SugaredLogger) *Controller {
	return &Controller{
		service: service,
		log:     log,
	}
}
