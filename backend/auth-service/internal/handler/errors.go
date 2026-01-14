package handler

import "errors"

var (
	ErrInvalidRequestBody = errors.New("invalid request body")
	ErrValidationFailed   = errors.New("validation failed")
)
