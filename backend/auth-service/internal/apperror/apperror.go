package apperror

import "errors"

var (
	// Request Errors
	ErrInvalidRequestBody = errors.New("invalid request body")
	ErrValidationFailed   = errors.New("validation failed")

	// Auth Errors
	ErrInvalidEmailOrPassword       = errors.New("invalid email or password")
	ErrRefreshTokenNotFound         = errors.New("refresh token not found")
	ErrInvalidOrExpiredRefreshToken = errors.New("invalid or expired refresh token")
	ErrGeneratingAccessToken        = errors.New("error generating access token")

	// User Errors
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)
