package auth

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"server/internal/models"
	"server/internal/repositories"
	"server/internal/utils"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// parse the request body
	var RequestBody models.User
	err := json.NewDecoder(r.Body).Decode(&RequestBody)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// get user by email
	userData, err := repositories.GetUserByEmail(r.Context(), RequestBody.Email)
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	// compare the hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(RequestBody.Password))
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
	err = repositories.SaveRefreshToken(r.Context(), userData.ID, refreshToken, expRefresh)
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
