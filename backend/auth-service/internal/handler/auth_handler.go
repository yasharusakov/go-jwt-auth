package handler

import (
	"auth-service/internal/model"
	"auth-service/internal/service"
	"auth-service/internal/util"
	"encoding/json"
	"net/http"
	"os"
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

func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "storage error"+err.Error(), http.StatusInternalServerError)
		return
	}
	if exists.Exists {
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
	user, err := h.service.RegisterUser(r.Context(), RequestBody.Email, hashedPassword)
	if err != nil {
		http.Error(w, "registration error"+err.Error(), http.StatusInternalServerError)
		return
	}

	// generate tokens
	accessToken, refreshToken, err := util.GenerateTokens(user.Id)
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
	util.SetRefreshTokenCookie(w, refreshToken)

	// get user by id
	userById, err := h.service.GetUserByID(r.Context(), user.Id)
	if err != nil {
		http.Error(w, "error retrieving user data", http.StatusInternalServerError)
		return
	}

	// return the new access token and user data
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token": accessToken,
		"user": map[string]interface{}{
			"id":    userById.User.Id,
			"email": userById.User.Email,
		},
	})
	if err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
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
	err = bcrypt.CompareHashAndPassword([]byte(userData.User.Password), []byte(reqBody.Password))
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	// generate access and refresh tokens
	accessToken, refreshToken, err := util.GenerateTokens(userData.User.Id)
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
	util.SetRefreshTokenCookie(w, refreshToken)

	// return the new access token and user data
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token": accessToken,
		"user": map[string]interface{}{
			"id":    userData.User.Id,
			"email": userData.User.Email,
		},
	})
	if err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

// update the access token
func (h *authHandler) Refresh(w http.ResponseWriter, r *http.Request) {
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
	claims, err := util.ValidateToken(cookie.Value, []byte(refreshSecret))
	if err != nil {
		http.Error(w, "invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	// extract user ID from claims
	userID := claims.Subject

	// generate new access token
	newAccessToken, err := util.GenerateToken(userID, jwtAccessTokenExp, []byte(accessSecret))
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
			"id":    userData.User.Id,
			"email": userData.User.Email,
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
	util.RemoveRefreshTokenCookie(w)
	w.WriteHeader(http.StatusOK)
}
