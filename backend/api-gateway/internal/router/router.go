package router

import (
	"api-gateway/internal/cache"
	"api-gateway/internal/config"
	"api-gateway/internal/logger"
	"api-gateway/internal/middleware"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

func SetupRoutes(app *fiber.App, redisCache cache.RedisCache, cfg config.Config) {
	app.Get("/health", func(c *fiber.Ctx) error {
		logger.Log.Info().Msg("Health check passed")
		return c.SendString("OK")
	})

	app.Get("/ready", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
		defer cancel()

		// ==================================
		// TODO: Rewrite to Fiber
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
		// ==================================

		// Check auth-service readiness
		if err := check(cfg.ApiAuthServiceInternalURL); err != nil {
			logger.Log.Warn().Err(err).Msg("auth-service is not ready")
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}

		// Check user-service readiness
		if err := check(cfg.ApiUserServiceInternalURL); err != nil {
			logger.Log.Warn().Err(err).Msg("user-service is not ready")
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}

		// Ping redis
		if err := redisCache.Ping(ctx); err != nil {
			logger.Log.Warn().Err(err).Msg("redis is not ready")
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}

		logger.Log.Info().Msg("Ready check passed")
		return c.SendString("OK")
	})

	auth := app.Group("/api/auth",
		middleware.CORS(cfg),
		middleware.RateLimit(redisCache),
	)
	auth.All("/*", proxyTo(cfg.ApiAuthServiceInternalURL))

	user := app.Group("/api/user",
		middleware.CORS(cfg),
		middleware.Auth(cfg),
	)
	user.All("/*", proxyTo(cfg.ApiUserServiceInternalURL))
}

func proxyTo(target string) fiber.Handler {
	// Remove /api from target if exists, since the path already contains /api/...
	target = strings.TrimSuffix(target, "/api")

	return func(c *fiber.Ctx) error {
		url := target + c.OriginalURL()
		return proxy.Do(c, url)
	}
}
