package route

import (
	"github.com/gofiber/fiber/v3"
	"go.uber.org/fx"
)

var Module = fx.Module("route",
	fx.Provide(
		AsRoute(NewHealthHandler),
		AsRoute(NewGraphqlHandler),
		AsRoute(NewPlaygroundHandler),
	),
)

type Route interface {
	Handler() fiber.Handler
	Pattern() string
}

func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Route)),
		fx.ResultTags(`group:"routes"`),
	)
}
