package router

import (
	grpcClient "auth-service/internal/client/grpc"
	"auth-service/internal/handler"
	"auth-service/internal/logger"
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(handlers handler.AuthHandler, db *pgxpool.Pool, grpcUserClient grpcClient.UserService) http.Handler {
	m := mux.NewRouter()

	m.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	m.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := db.Ping(ctx); err != nil {
			logger.Log.Error().
				Err(err).
				Msg("Database ping failed")
			http.Error(w, "database not ready", http.StatusServiceUnavailable)
			return
		}

		if err := grpcUserClient.Ping(ctx); err != nil {
			logger.Log.Error().
				Err(err).
				Msg("gRPC user client ping failed")
			http.Error(w, "gRPC user client not ready", http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("READY"))
	}).Methods("GET")

	m.HandleFunc("/api/auth/register", handlers.Register).Methods("POST")
	m.HandleFunc("/api/auth/login", handlers.Login).Methods("POST")
	m.HandleFunc("/api/auth/refresh", handlers.Refresh).Methods("GET")
	m.HandleFunc("/api/auth/logout", handlers.Logout).Methods("POST")

	return m
}
