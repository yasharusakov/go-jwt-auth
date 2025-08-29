package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"server/internal/model"
	"server/internal/service"
	"server/internal/utils"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) AuthHandler {
	return &authHandler{service: service}
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// parse the request body
	var reqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// get user by email
	userData, err := h.service.GetUserByEmail(r.Context(), reqBody.Email)
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	// compare the hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(reqBody.Password))
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	// generate access and refresh tokens
	accessToken, refreshToken, err := utils.GenerateTokens(userData.ID)
	if err != nil {
		http.Error(w, "error creating tokens", http.StatusInternalServerError)
		return
	}

	// parse refresh token expiration
	expRefresh, err := time.ParseDuration(os.Getenv("JWT_REFRESH_TOKEN_EXPIRATION"))
	if err != nil {
		http.Error(w, "error parsing refresh token expiration", http.StatusInternalServerError)
		return
	}

	// save refresh token to database
	err = h.service.SaveRefreshToken(r.Context(), userData.ID, refreshToken, expRefresh)
	if err != nil {
		http.Error(w, "error saving refresh token", http.StatusInternalServerError)
		return
	}

	// set refresh_token in httponly cookie
	utils.SetRefreshTokenCookie(w, refreshToken, expRefresh)

	// return the new access token and user data
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token": accessToken,
		"user": map[string]interface{}{
			"id":    userData.ID,
			"email": userData.Email,
		},
	})
	if err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

// update the access token
func (h *authHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// load jwt secrets and expiration from environment variables
	accessSecret := os.Getenv("JWT_ACCESS_TOKEN_SECRET")
	refreshSecret := os.Getenv("JWT_REFRESH_TOKEN_SECRET")
	jwtAccessTokenExp, err := time.ParseDuration(os.Getenv("JWT_ACCESS_TOKEN_EXPIRATION"))
	if err != nil {
		http.Error(w, "invalid access token expiration duration", http.StatusInternalServerError)
		return
	}

	// check if the refresh token cookie exists
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "refresh token not found", http.StatusUnauthorized)
		return
	}

	// validate the refresh token
	claims, err := utils.ValidateToken(cookie.Value, []byte(refreshSecret))
	if err != nil {
		http.Error(w, "invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	// extract user ID from claims
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
		return
	}

	// generate new access token
	newAccessToken, err := utils.GenerateToken(userID, jwtAccessTokenExp, []byte(accessSecret))
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
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token": newAccessToken,
		"user": map[string]interface{}{
			"id":    userData.ID,
			"email": userData.Email,
		},
	})
	if err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// parse the request body
	var RequestBody model.User
	err := json.NewDecoder(r.Body).Decode(&RequestBody)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// validate user input
	if RequestBody.Email == "" || RequestBody.Password == "" {
		http.Error(w, "email or password is empty", http.StatusBadRequest)
		return
	}

	// check if user with the given email already exists
	exists, err := h.service.CheckUserExistsByEmail(r.Context(), RequestBody.Email)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "user already exists", http.StatusBadRequest)
		return
	}

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(RequestBody.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "error hashing password", http.StatusInternalServerError)
		return
	}

	// registration query
	userId, err := h.service.RegisterUser(r.Context(), RequestBody.Email, hashedPassword)
	if err != nil {
		http.Error(w, "registration error", http.StatusInternalServerError)
		return
	}

	// generate tokens
	accessToken, refreshToken, err := utils.GenerateTokens(userId)
	if err != nil {
		http.Error(w, "error creating tokens", http.StatusInternalServerError)
		return
	}

	// save refresh token to database
	expRefresh, err := time.ParseDuration(os.Getenv("JWT_REFRESH_TOKEN_EXPIRATION"))
	if err != nil {
		http.Error(w, "error parsing refresh token expiration", http.StatusInternalServerError)
		return
	}

	// save refresh token in the database
	err = h.service.SaveRefreshToken(r.Context(), userId, refreshToken, expRefresh)
	if err != nil {
		http.Error(w, "error saving refresh token", http.StatusInternalServerError)
		return
	}

	// set refresh_token in httponly cookie
	utils.SetRefreshTokenCookie(w, refreshToken, expRefresh)

	// get user by id
	userData, err := h.service.GetUserByID(r.Context(), userId)
	if err != nil {
		http.Error(w, "error retrieving user data", http.StatusInternalServerError)
		return
	}

	// return the new access token and user data
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token": accessToken,
		"user": map[string]interface{}{
			"id":    userData.ID,
			"email": userData.Email,
		},
	})
	if err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err == nil {
		_ = h.service.RemoveRefreshToken(r.Context(), cookie.Value)
	}
	utils.RemoveRefreshTokenCookie(w)
	w.WriteHeader(http.StatusOK)
}
