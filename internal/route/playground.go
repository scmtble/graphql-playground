package route

import (
	"log/slog"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
)

type PlaygroundHandler struct {
	logger *slog.Logger
}

func (h *PlaygroundHandler) Pattern() string {
	return "/playground"
}

func (h *PlaygroundHandler) Handler() fiber.Handler {
	return adaptor.HTTPHandler(playground.ApolloSandboxHandler("GraphQL playground", "/graphql"))
}

func NewPlaygroundHandler(logger *slog.Logger) *PlaygroundHandler {
	return &PlaygroundHandler{
		logger: logger,
	}
}
