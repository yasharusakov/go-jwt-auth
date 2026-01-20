package router

import (
	"api-gateway/internal/config"
	"api-gateway/internal/middleware"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
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
