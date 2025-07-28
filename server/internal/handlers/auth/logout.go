package auth

import (
	"net/http"
	"server/internal/utils"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	utils.RemoveRefreshTokenCookie(w, r)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}
