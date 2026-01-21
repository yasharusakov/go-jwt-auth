package middleware

import (
	"api-gateway/internal/config"
	"api-gateway/internal/logger"
	"api-gateway/internal/util"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	accessSecret := config.GetConfig().JWTAccessTokenSecret
	return func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Log.Warn().
				Str("path", r.URL.Path).
				Str("ip", r.RemoteAddr).
				Msg("access token not found")

			http.Error(w, "access token not found", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		_, err := util.ValidateToken(tokenStr, []byte(accessSecret))
		if err != nil {
			logger.Log.Warn().
				Err(err).
				Str("path", r.URL.Path).
				Str("ip", r.RemoteAddr).
				Msg("invalid or expired token")
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
