package apperror

import "github.com/gofiber/fiber/v2"

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Client Errors
func BadRequest(msg string) *AppError {
	return &AppError{Code: fiber.StatusBadRequest, Message: msg}
}

func Unauthorized(msg string) *AppError {
	return &AppError{Code: fiber.StatusUnauthorized, Message: msg}
}

func NotFound(msg string) *AppError {
	return &AppError{Code: fiber.StatusNotFound, Message: msg}
}

func Conflict(msg string) *AppError {
	return &AppError{Code: fiber.StatusConflict, Message: msg}
}

// Server Errors
func Internal(err error) *AppError {
	return &AppError{
		Code:    fiber.StatusInternalServerError,
		Message: "internal server error",
		Err:     err,
	}
}
