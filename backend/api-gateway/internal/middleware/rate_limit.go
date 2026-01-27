package middleware

import (
	"api-gateway/internal/cache"
	"api-gateway/internal/util"
	"net/http"
)

type RateLimitMiddleware interface {
	RateLimit(next http.HandlerFunc) http.HandlerFunc
}

type rateLimitMiddleware struct {
	cache.RedisCache
}

func NewRateLimitMiddleware(c cache.RedisCache) RateLimitMiddleware {
	return &rateLimitMiddleware{
		RedisCache: c,
	}
}

func (rlm *rateLimitMiddleware) RateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := util.GetClientIP(r)

		result, err := rlm.CheckRateLimit(r.Context(), clientIP)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		if !result {
			w.WriteHeader(http.StatusTooManyRequests)
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	}
}
