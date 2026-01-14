package handler

import (
	"auth-service/internal/config"
	"auth-service/internal/service"
	"auth-service/internal/util"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	decodeAndValidate(r *http.Request, dst any) error
}

type authHandler struct {
	service   service.AuthService
	validator *validator.Validate
}

func NewAuthHandler(service service.AuthService) AuthHandler {
	return &authHandler{
		service:   service,
		validator: validator.New(),
	}
}

func (h *authHandler) decodeAndValidate(r *http.Request, dst any) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidRequestBody, err)
	}
	if err := h.validator.Struct(dst); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return fmt.Errorf("%w: %w", ErrValidationFailed, validationErrors)
		}
		return ErrValidationFailed
	}
	return nil
}

func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	// parse the request body
	var req RegisterRequest
	if err := h.decodeAndValidate(r, &req); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO!: Layers (handler -> service -> repository)
	// TODO!: Move business logic to service layer
	// TODO!: make cfg global and do not call it every time in handler

	// check if user with the given email already exists
	exists, err := h.service.CheckUserExistsByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "storage error"+err.Error(), http.StatusInternalServerError)
		return
	}
	if exists.Exists {
		http.Error(w, "user already exists", http.StatusBadRequest)
		return
	}

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "error hashing password", http.StatusInternalServerError)
		return
	}

	// registration query
	user, err := h.service.RegisterUser(r.Context(), req.Email, hashedPassword)
	if err != nil {
		http.Error(w, "registration error"+err.Error(), http.StatusInternalServerError)
		return
	}

	cfg := config.GetConfig()
	// generate tokens
	accessToken, refreshToken, err := util.GenerateTokens(user.Id, cfg.JWT)
	if err != nil {
		http.Error(w, "error creating tokens"+err.Error(), http.StatusInternalServerError)
		return
	}

	// save refresh token in the storage
	err = h.service.SaveRefreshToken(r.Context(), user.Id, refreshToken)
	if err != nil {
		http.Error(w, "error saving refresh token"+err.Error(), http.StatusInternalServerError)
		return
	}

	// set refresh_token in httponly cookie
	util.SetRefreshTokenCookie(w, refreshToken, cfg.JWT.JWTRefreshTokenExp, cfg.AppEnv == "production")

	// get user by id
	userById, err := h.service.GetUserByID(r.Context(), user.Id)
	if err != nil {
		http.Error(w, "error retrieving user data", http.StatusInternalServerError)
		return
	}

	// return the new access token and user data
	writeJSON(w, http.StatusOK, AuthResponse{
		AccessToken: accessToken,
		User: UserResponse{
			ID:    userById.User.Id,
			Email: userById.User.Email,
		},
	})
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	// parse the request body
	var req LoginRequest
	if err := h.decodeAndValidate(r, &req); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get user by email
	userData, err := h.service.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	// compare the hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(userData.User.Password), []byte(req.Password))
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	cfg := config.GetConfig()

	// generate access and refresh tokens
	accessToken, refreshToken, err := util.GenerateTokens(userData.User.Id, cfg.JWT)
	if err != nil {
		http.Error(w, "error creating tokens", http.StatusInternalServerError)
		return
	}

	// save refresh token to storage
	err = h.service.SaveRefreshToken(r.Context(), userData.User.Id, refreshToken)
	if err != nil {
		http.Error(w, "error saving refresh token", http.StatusInternalServerError)
		return
	}

	// set refresh_token in httponly cookie
	util.SetRefreshTokenCookie(w, refreshToken, cfg.JWT.JWTRefreshTokenExp, cfg.AppEnv == "production")

	// return the new access token and user data
	writeJSON(w, http.StatusOK, AuthResponse{
		AccessToken: accessToken,
		User: UserResponse{
			ID:    userData.User.Id,
			Email: userData.User.Email,
		},
	})
}

// update the access token
func (h *authHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	// load jwt secrets and expiration from environment variables
	cfg := config.GetConfig()

	// check if the refresh token cookie exists
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "refresh token not found", http.StatusUnauthorized)
		return
	}

	// validate the refresh token
	claims, err := util.ValidateToken(cookie.Value, []byte(cfg.JWT.JWTRefreshTokenSecret))
	if err != nil {
		http.Error(w, "invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	// extract user ID from claims
	userID := claims.Subject

	// generate new access token
	newAccessToken, err := util.GenerateToken(userID, cfg.JWT.JWTAccessTokenExp, []byte(cfg.JWT.JWTAccessTokenSecret))
	if err != nil {
		http.Error(w, "error generating new access token", http.StatusInternalServerError)
		return
	}

	userData, err := h.service.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// return the new access token and user data
	writeJSON(w, http.StatusOK, AuthResponse{
		AccessToken: newAccessToken,
		User: UserResponse{
			ID:    userData.User.Id,
			Email: userData.User.Email,
		},
	})
}

func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err == nil {
		_ = h.service.RemoveRefreshToken(r.Context(), cookie.Value)
	}

	cfg := config.GetConfig()

	util.RemoveRefreshTokenCookie(w, cfg.AppEnv == "production")
	w.WriteHeader(http.StatusOK)
}
