package handler

import (
	"auth-service/internal/config"
	"auth-service/internal/dto"
	"auth-service/internal/logger"
	"auth-service/internal/service"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}

type authHandler struct {
	service service.AuthService
	cfg     config.Config
}

func NewAuthHandler(service service.AuthService, cfg config.Config) AuthHandler {
	return &authHandler{
		service: service,
		cfg:     cfg,
	}
}

func (h *authHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		logger.Log.Error().Err(err).Msg("failed to parse request body")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	result, err := h.service.Register(c.Context(), req.Email, req.Password)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to register user")
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    result.RefreshToken,
		Expires:  time.Now().Add(h.cfg.JWT.JWTRefreshTokenExp),
		HTTPOnly: true,
		Path:     "/",
		Secure:   h.cfg.AppEnv == "production",
		SameSite: fiber.CookieSameSiteLaxMode,
	})

	return c.Status(fiber.StatusOK).JSON(dto.AuthResponse{
		AccessToken: result.AccessToken,
		User: dto.UserResponse{
			ID:    result.UserID,
			Email: result.Email,
		},
	})
}

func (h *authHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		logger.Log.Error().Err(err).Msg("failed to parse request body")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	result, err := h.service.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to login")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    result.RefreshToken,
		Expires:  time.Now().Add(h.cfg.JWT.JWTRefreshTokenExp),
		HTTPOnly: true,
		Path:     "/",
		Secure:   h.cfg.AppEnv == "production",
		SameSite: fiber.CookieSameSiteLaxMode,
	})

	return c.Status(fiber.StatusOK).JSON(dto.AuthResponse{
		AccessToken: result.AccessToken,
		User: dto.UserResponse{
			ID:    result.UserID,
			Email: result.Email,
		},
	})
}

func (h *authHandler) Refresh(c *fiber.Ctx) error {
	cookie := c.Cookies("refresh_token")
	if cookie == "" {
		logger.Log.Warn().Msg("refresh token not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "refresh token not found",
		})
	}

	result, err := h.service.Refresh(c.Context(), cookie)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to refresh tokens")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    result.RefreshToken,
		Expires:  time.Now().Add(h.cfg.JWT.JWTRefreshTokenExp),
		HTTPOnly: true,
		Path:     "/",
		Secure:   h.cfg.AppEnv == "production",
		SameSite: fiber.CookieSameSiteLaxMode,
	})

	return c.Status(fiber.StatusOK).JSON(dto.AuthResponse{
		AccessToken: result.AccessToken,
		User: dto.UserResponse{
			ID:    result.UserID,
			Email: result.Email,
		},
	})
}

func (h *authHandler) Logout(c *fiber.Ctx) error {
	cookie := c.Cookies("refresh_token")

	if cookie != "" {
		_ = h.service.Logout(c.Context(), cookie)
	}

	logger.Log.Info().Msg("user logged out")

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HTTPOnly: true,
		Secure:   h.cfg.AppEnv == "production",
		SameSite: fiber.CookieSameSiteLaxMode,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully logged out",
	})
}
