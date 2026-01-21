package router

import (
	"context"
	"net/http"
	"time"
	httpHandler "user-service/internal/handler/http"
	"user-service/internal/logger"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(handlers httpHandler.Handlers, db *pgxpool.Pool) http.Handler {
	m := mux.NewRouter()

	m.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	m.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := db.Ping(ctx); err != nil {
			logger.Log.Error().Err(err).Msg("database is not ready")
			http.Error(w, "database is not ready", http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("READY"))
	})

	m.HandleFunc("/api/user/get-all", handlers.GetAllUsers).Methods("GET")

	return m
}
