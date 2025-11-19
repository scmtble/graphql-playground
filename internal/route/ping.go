package route

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

type PingHandler struct {
	logger *slog.Logger
}

func (h *PingHandler) Pattern() string {
	return "/ping"
}

func (h *PingHandler) Handler() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		h.logger.Info("Ping received")
		return ctx.SendString("pong")
	}
}

func NewPingHandler(logger *slog.Logger) *PingHandler {
	return &PingHandler{
		logger: logger,
	}
}
