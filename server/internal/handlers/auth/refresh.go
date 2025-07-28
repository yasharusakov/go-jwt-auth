package auth

import (
	"encoding/json"
	"net/http"
	"os"
	"server/internal/repositories"
	"server/internal/utils"
	"strconv"
	"time"
)

// update the access token
func Refresh(w http.ResponseWriter, r *http.Request) {
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

	userData, err := repositories.GetUserByID(r.Context(), userID)
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
