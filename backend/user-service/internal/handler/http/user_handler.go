package httpHandler

import (
	"user-service/internal/logger"
	"user-service/internal/service"

	"github.com/gofiber/fiber/v2"
)

type Handlers interface {
	GetAllUsers(c *fiber.Ctx) error
}

type userHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) Handlers {
	return &userHandler{
		service: service,
	}
}

func (h *userHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.service.GetAllUsers(c.Context())
	if err != nil {
		logger.Log.Info().Err(err).Msg("Error retrieving all users")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error retrieving all users",
		})
	}

	return c.JSON(users)
}
