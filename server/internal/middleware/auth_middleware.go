package middleware

import (
	"net/http"
	"os"
	"server/internal/utils"
	"strings"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessSecret := os.Getenv("JWT_ACCESS_TOKEN_SECRET")

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "access token not found", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ValidateToken(tokenStr, []byte(accessSecret))
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		r.Header.Set("X-User-ID", claims.Subject)
		next(w, r)
	}
}
