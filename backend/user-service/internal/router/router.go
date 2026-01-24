package router

import (
	"context"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"net/http"
	"time"
	httpHandler "user-service/internal/handler/http"
	"user-service/internal/logger"
)

func RegisterRoutes(handlers httpHandler.Handlers, db *gorm.DB) http.Handler {
	m := mux.NewRouter()

	m.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		logger.Log.Info().Msg("Health check passed")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

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
			logger.Log.Error().Err(err).Msg("database is not ready")
			http.Error(w, "database is not ready", http.StatusServiceUnavailable)
			return
		}

		logger.Log.Info().Msg("Ready check passed")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("READY"))
	})

	m.HandleFunc("/api/user/get-all", handlers.GetAllUsers).Methods("GET")

	return m
}
