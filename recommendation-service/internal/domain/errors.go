package domain

import "errors"

var (
	ErrRecommendationsNotFound   = errors.New("recommendations not found")
	ErrInvalidReviewCreatedEvent = errors.New("invalid review created event")
)
