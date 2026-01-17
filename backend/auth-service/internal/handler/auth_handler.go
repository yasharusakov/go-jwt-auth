package handler

import (
	"auth-service/internal/apperror"
	"auth-service/internal/config"
	"auth-service/internal/dto"
	"auth-service/internal/httpresponse"
	"auth-service/internal/service"
	"auth-service/internal/util"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type AuthHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	service   service.AuthService
	validator *validator.Validate
	cfg       config.Config
}

func NewAuthHandler(service service.AuthService, cfg config.Config) AuthHandler {
	return &authHandler{
		service:   service,
		validator: validator.New(),
		cfg:       cfg,
	}
}

func (h *authHandler) decodeAndValidate(r *http.Request, dst any) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return fmt.Errorf("%w: %w", apperror.ErrInvalidRequestBody, err)
	}
	if err := h.validator.Struct(dst); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return fmt.Errorf("%w: %w", apperror.ErrValidationFailed, validationErrors)
		}
		return apperror.ErrValidationFailed
	}
	return nil
}

func (h *authHandler) respondWithError(err error, w http.ResponseWriter) {
	switch {
	case errors.Is(err, apperror.ErrUserAlreadyExists):
		httpresponse.WriteError(w, "user already exists", http.StatusConflict)
	case errors.Is(err, apperror.ErrInvalidEmailOrPassword):
		httpresponse.WriteError(w, "invalid email or password", http.StatusUnauthorized)
	case errors.Is(err, apperror.ErrRefreshTokenNotFound):
		httpresponse.WriteError(w, "refresh token not found", http.StatusUnauthorized)
	case errors.Is(err, apperror.ErrInvalidOrExpiredRefreshToken):
		httpresponse.WriteError(w, "invalid or expired refresh token", http.StatusUnauthorized)
	case errors.Is(err, apperror.ErrUserNotFound):
		httpresponse.WriteError(w, "user not found", http.StatusUnauthorized)
	case errors.Is(err, apperror.ErrValidationFailed), errors.Is(err, apperror.ErrInvalidRequestBody):
		httpresponse.WriteError(w, err.Error(), http.StatusBadRequest)
	default:
		log.Printf("INTERNAL SERVER ERROR: %v", err)
		httpresponse.WriteError(w, "internal server error", http.StatusInternalServerError)
	}
}

func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := h.decodeAndValidate(r, &req); err != nil {
		h.respondWithError(err, w)
		return
	}

	result, err := h.service.Register(r.Context(), req.Email, req.Password)

	if err != nil {
		h.respondWithError(err, w)
		return
	}

	util.SetRefreshTokenCookie(w, result.RefreshToken, h.cfg.JWT.JWTRefreshTokenExp, h.cfg.AppEnv == "production")

	httpresponse.WriteJSON(w, http.StatusOK, dto.AuthResponse{
		AccessToken: result.AccessToken,
		User:        result.User,
	})
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := h.decodeAndValidate(r, &req); err != nil {
		h.respondWithError(err, w)
		return
	}

	result, err := h.service.Login(r.Context(), req.Email, req.Password)

	if err != nil {
		h.respondWithError(err, w)
		return
	}

	util.SetRefreshTokenCookie(w, result.RefreshToken, h.cfg.JWT.JWTRefreshTokenExp, h.cfg.AppEnv == "production")

	httpresponse.WriteJSON(w, http.StatusOK, dto.AuthResponse{
		AccessToken: result.AccessToken,
		User:        result.User,
	})
}

func (h *authHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		httpresponse.WriteError(w, "refresh token not found", http.StatusUnauthorized)
		return
	}

	result, err := h.service.Refresh(r.Context(), cookie.Value)

	if err != nil {
		h.respondWithError(err, w)
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, dto.AuthResponse{
		AccessToken: result.AccessToken,
		User:        result.User,
	})
}

func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err == nil {
		_ = h.service.Logout(r.Context(), cookie.Value)
	}

	util.RemoveRefreshTokenCookie(w, h.cfg.AppEnv == "production")
	w.WriteHeader(http.StatusOK)
}
