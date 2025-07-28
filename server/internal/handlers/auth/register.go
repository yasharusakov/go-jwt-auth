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

func Register(w http.ResponseWriter, r *http.Request) {
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

	// validate user input
	if RequestBody.Email == "" || RequestBody.Password == "" {
		http.Error(w, "email or password is empty", http.StatusBadRequest)
		return
	}

	// check if user with the given email already exists
	exists, err := repositories.CheckUserExistsByEmail(r.Context(), RequestBody.Email)
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
	userId, err := repositories.RegisterUser(r.Context(), RequestBody.Email, hashedPassword)
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
	err = repositories.SaveRefreshToken(r.Context(), userId, refreshToken, expRefresh)
	if err != nil {
		http.Error(w, "error saving refresh token", http.StatusInternalServerError)
		return
	}

	// set refresh_token in httponly cookie
	utils.SetRefreshTokenCookie(w, refreshToken, expRefresh)

	// get user by id
	userData, err := repositories.GetUserByID(r.Context(), userId)
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
