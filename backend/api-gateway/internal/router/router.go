package router

import (
	"api-gateway/internal/config"
	"api-gateway/internal/middleware"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

func newProxy(target string) *httputil.ReverseProxy {
	parsedURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Error occurred while parsing URL: %v", err)
	}

	return httputil.NewSingleHostReverseProxy(parsedURL)
}

func RegisterRoutes() *http.ServeMux {
	cfg := config.GetConfig()
	mux := http.NewServeMux()

	authProxy := newProxy(cfg.ApiAuthServiceInternalURL)
	userProxy := newProxy(cfg.ApiUserServiceInternalURL)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		check := func(target string) error {
			baseURL := strings.TrimSuffix(target, "/api")

			req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/health", nil)
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

		// Check auth-service health
		if err := check(cfg.ApiAuthServiceInternalURL); err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		// Check user-service health
		if err := check(cfg.ApiUserServiceInternalURL); err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

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
