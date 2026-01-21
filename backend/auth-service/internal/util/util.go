package util

import (
	"net/http"
	"time"
)

// TODO: Modify the cookie path "const RefreshTokenCookiePath = "/api/auth"

func SetRefreshTokenCookie(w http.ResponseWriter, refreshToken string, exp time.Duration, isProd bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(exp),
		HttpOnly: true,
		Path:     "/",
		Secure:   isProd,
		SameSite: http.SameSiteLaxMode,
	})
}

func RemoveRefreshTokenCookie(w http.ResponseWriter, isProd bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   isProd,
		SameSite: http.SameSiteLaxMode,
	})
}
