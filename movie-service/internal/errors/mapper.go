package errors

import "errors"

var (
	ErrInvalidTitle = errors.New("invalid movie title")
	ErrInvalidYear  = errors.New("invalid movie year")
)
