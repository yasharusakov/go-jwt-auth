package router

import (
	grpcClient "auth-service/internal/client/grpc"
	"auth-service/internal/handler"
	"auth-service/internal/logger"
	"context"
	"net/http"
	"time"

	_ "auth-service/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"
)

func RegisterRoutes(handlers handler.AuthHandler, db *gorm.DB, grpcUserClient grpcClient.UserService) http.Handler {
	m := mux.NewRouter()

	m.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		logger.Log.Info().Msg("Health check passed")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	m.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		sqlDB, err := db.WithContext(ctx).DB()
		if err != nil {
			logger.Log.Error().Err(err).Msg("database is not ready")
			http.Error(w, "database is not ready", http.StatusServiceUnavailable)
			return
		}

		if err := sqlDB.Ping(); err != nil {
			logger.Log.Error().Err(err).Msg("database ping failed")
			http.Error(w, "database is not ready", http.StatusServiceUnavailable)
			return
		}

		if err := grpcUserClient.Ping(ctx); err != nil {
			logger.Log.Error().Err(err).Msg("gRPC user client ping failed")
			http.Error(w, "gRPC user client not ready", http.StatusServiceUnavailable)
			return
		}

		logger.Log.Info().Msg("Ready check passed")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("READY"))
	}).Methods("GET")

	m.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	m.HandleFunc("/api/auth/register", handlers.Register).Methods("POST")
	m.HandleFunc("/api/auth/login", handlers.Login).Methods("POST")
	m.HandleFunc("/api/auth/refresh", handlers.Refresh).Methods("GET")
	m.HandleFunc("/api/auth/logout", handlers.Logout).Methods("POST")

	return m
}
