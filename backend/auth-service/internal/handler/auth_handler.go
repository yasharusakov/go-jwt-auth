package handler

import (
	"auth-service/internal/apperror"
	"auth-service/internal/config"
	"auth-service/internal/dto"
	"auth-service/internal/httpresponse"
	"auth-service/internal/logger"
	"auth-service/internal/service"
	"auth-service/internal/util"
	"encoding/json"
	"errors"
	"fmt"
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
		logger.Log.Error().Err(err).Msg("internal server error")
		httpresponse.WriteError(w, "internal server error", http.StatusInternalServerError)
	}
}

// Register godoc
// @Summary      Register a new user
// @Description  Creates a new user account with the provided email and password.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterRequest true "Email и Password"
// @Success      200 {object} dto.AuthResponse "Successful registration"
// @Failure      400 {object} map[string]string "Validation error"
// @Failure      401 {object} map[string]string "Invalid email or password"
// @Failure      409 {object} map[string]string "User already exists"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /register [post]
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

	logger.Log.Info().
		Str("id", result.UserID).
		Str("email", result.Email).
		Msg("user registered")

	util.SetRefreshTokenCookie(w, result.RefreshToken, h.cfg.JWT.JWTRefreshTokenExp, h.cfg.AppEnv == "production")

	httpresponse.WriteJSON(w, http.StatusOK, dto.AuthResponse{
		AccessToken: result.AccessToken,
		User: dto.UserResponse{
			ID:    result.UserID,
			Email: result.Email,
		},
	})
}

// Login godoc
// @Summary      User login
// @Description  Authenticates a user with email and password.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Credentials"
// @Success      200 {object} dto.AuthResponse "Successful login"
// @Failure      400 {object} map[string]string "Validation error"
// @Failure      401 {object} map[string]string "Invalid email or password"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /login [post]
func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := h.decodeAndValidate(r, &req); err != nil {
		h.respondWithError(err, w)
		return
	}

	result, err := h.service.Login(r.Context(), req.Email, req.Password)

	if err != nil {
		logger.Log.Warn().
			Str("email", req.Email).
			Msg("login failed")
		h.respondWithError(err, w)
		return
	}

	logger.Log.Info().
		Str("id", result.UserID).
		Str("email", result.Email).
		Msg("user logged in")

	util.SetRefreshTokenCookie(w, result.RefreshToken, h.cfg.JWT.JWTRefreshTokenExp, h.cfg.AppEnv == "production")

	httpresponse.WriteJSON(w, http.StatusOK, dto.AuthResponse{
		AccessToken: result.AccessToken,
		User: dto.UserResponse{
			ID:    result.UserID,
			Email: result.Email,
		},
	})
}

// Refresh godoc
// @Summary      Update tokens
// @Description  Generates new access and refresh tokens using a valid refresh token from cookies.
// @Tags         auth
// @Produce      json
// @Success      200 {object} dto.AuthResponse "New tokens"
// @Failure      401 {object} map[string]string "Refresh token invalid or expired"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /refresh [post]
func (h *authHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		logger.Log.Warn().
			Msg("refresh token not found")
		httpresponse.WriteError(w, "refresh token not found", http.StatusUnauthorized)
		return
	}

	result, err := h.service.Refresh(r.Context(), cookie.Value)

	if err != nil {
		logger.Log.Warn().
			Msg("refresh token invalid or expired")
		h.respondWithError(err, w)
		return
	}

	util.SetRefreshTokenCookie(w, result.RefreshToken, h.cfg.JWT.JWTRefreshTokenExp, h.cfg.AppEnv == "production")

	httpresponse.WriteJSON(w, http.StatusOK, dto.AuthResponse{
		AccessToken: result.AccessToken,
		User: dto.UserResponse{
			ID:    result.UserID,
			Email: result.Email,
		},
	})
}

// Logout godoc
// @Summary      Logout user
// @Description  Removes the refresh token cookies and invalidates the refresh token from PostgreSQL.
// @Tags         auth
// @Produce      json
// @Success      200 {object} map[string]string "Successful logout"
// @Router       /logout [post]
func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")

	if err == nil {
		_ = h.service.Logout(r.Context(), cookie.Value)
	}

	logger.Log.Info().Msg("user logged out")
	util.RemoveRefreshTokenCookie(w, h.cfg.AppEnv == "production")
	w.WriteHeader(http.StatusOK)
}
