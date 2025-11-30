package route

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

type HealthHandler struct {
	logger *slog.Logger
}

func (h *HealthHandler) Pattern() string {
	return "/health"
}

func (h *HealthHandler) Handler() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		return ctx.SendString("Healthy")
	}
}

func NewHealthHandler(logger *slog.Logger) *HealthHandler {
	return &HealthHandler{
		logger: logger,
	}
}
