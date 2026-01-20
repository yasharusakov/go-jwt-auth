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

	// Do not show to client
	ErrGeneratingAccessToken = errors.New("error generating access token")
	ErrGeneratingTokens      = errors.New("error generating tokens")

	// User Errors
	ErrUserAlreadyExists = errors.New("user already exists")

	// Do not show to client (only in refresh)
	ErrUserNotFound = errors.New("user not found")
)
