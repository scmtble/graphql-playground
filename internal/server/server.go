package server

import (
	"github.com/gofiber/fiber/v3"
	"github.com/scmtble/graphql-playground/internal/route"
	"go.uber.org/fx"
)

var Module = fx.Module("server",
	fx.Provide(
		fx.Annotate(
			NewServer,
			fx.ParamTags(`group:"routes"`),
		),
	),
	fx.Invoke(RunServer),
)

func NewServer(routes []route.Route) (*fiber.App, error) {
	app := fiber.New(fiber.Config{})

	for _, r := range routes {
		app.All(r.Pattern(), r.Handler())
	}

	return app, nil
}

func RunServer(lc fx.Lifecycle, app *fiber.App) {
	lc.Append(fx.StartHook(func() {
		go func() {
			if err := app.Listen(":8080", fiber.ListenConfig{
				DisableStartupMessage: true,
			}); err != nil {
				panic(err)
			}
		}()
	}))
	lc.Append(fx.StopHook(app.ShutdownWithContext))
}
