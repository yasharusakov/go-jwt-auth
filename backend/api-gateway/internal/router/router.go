package router

import (
	"api-gateway/internal/config"
	"api-gateway/internal/logger"
	"api-gateway/internal/middleware"
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

func newProxy(target string) *httputil.ReverseProxy {
	parsedURL, err := url.Parse(target)
	if err != nil {
		logger.Log.Fatal().
			Err(err).
			Str("url", target).
			Msg("error occurred while parsing URL")
	}

	return httputil.NewSingleHostReverseProxy(parsedURL)
}

func RegisterRoutes(cfg config.Config) *http.ServeMux {
	mux := http.NewServeMux()

	authProxy := newProxy(cfg.ApiAuthServiceInternalURL)
	userProxy := newProxy(cfg.ApiUserServiceInternalURL)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		logger.Log.Info().Msg("Health check passed")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		check := func(target string) error {
			baseURL := strings.TrimSuffix(target, "/api")

			req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/ready", nil)
			if err != nil {
				return fmt.Errorf("failed to create request for %s: %w", baseURL, err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return fmt.Errorf("%s is not ready: %w", baseURL, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("%s is not ready: status %d", baseURL, resp.StatusCode)
			}

			return nil
		}

		// Check auth service readiness
		if err := check(cfg.ApiAuthServiceInternalURL); err != nil {
			logger.Log.Warn().
				Err(err).
				Msg("auth-service is not ready")
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		// Check user service readiness
		if err := check(cfg.ApiUserServiceInternalURL); err != nil {
			logger.Log.Warn().
				Err(err).
				Msg("user-service is not ready")
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		logger.Log.Info().Msg("Ready check passed")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("READY"))
	})

	mux.Handle("/api/auth/", middleware.CORSMiddleware(
		http.StripPrefix("/api/", authProxy).ServeHTTP,
	))
	mux.Handle("/api/user/", middleware.CORSMiddleware(
		middleware.AuthMiddleware(
			http.StripPrefix("/api/", userProxy).ServeHTTP,
		),
	))

	return mux
}
