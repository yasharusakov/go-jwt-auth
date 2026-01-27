package app

import (
	"api-gateway/internal/cache"
	"api-gateway/internal/config"
	"api-gateway/internal/logger"
	"api-gateway/internal/router"
	"api-gateway/internal/server"
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	cfg := config.GetConfig()
	logger.Init(cfg.AppEnv)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	redisCache := cache.NewRedisCache(cfg.RedisConfig)
	defer redisCache.Close()

	handlers := router.RegisterRoutes(redisCache, cfg)

	srv := &server.HttpServer{}
	serverErrors := make(chan error, 1)

	go func() {
		logger.Log.Info().
			Str("port", cfg.ApiGatewayInternalPort).
			Msg("starting HTTP server")

		serverErrors <- srv.Run(cfg.ApiGatewayInternalPort, handlers)
	}()

	select {
	case <-ctx.Done():
		logger.Log.Info().Msg("shutting down HTTP server")
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Fatal().
				Err(err).
				Msg("server error")
		}
		return
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error().
			Err(err).
			Msg("server shutdown error")
	} else {
		logger.Log.Info().Msg("server stopped gracefully")
	}

}
