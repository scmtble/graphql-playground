package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/ravilushqa/otelgqlgen"
	"github.com/scmtble/graphql-playground/graphql"
	"github.com/scmtble/graphql-playground/graphql/generated"
	"github.com/vektah/gqlparser/v2/ast"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func NewLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{},
	))
}

func NewPingHandler(logger *slog.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		logger.Info("ping received",
			slog.String("path", c.Path()),
			slog.String("method", c.Method()),
			slog.Any("ip", c.IP()),
			slog.String("host", c.Host()),
			slog.String("host_name", c.Hostname()),
			slog.Any("header", c.GetReqHeaders()),
		)
		return c.SendString("pong")
	}
}
func NewGraphqlHandler() fiber.Handler {
	srv := handler.New(generated.NewExecutableSchema(generated.Config{
		Resolvers:  &graphql.Resolver{},
		Directives: generated.DirectiveRoot{},
	}))
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(otelgqlgen.Middleware())
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return adaptor.HTTPHandler(srv)
}

func NewPlaygroundHandler() fiber.Handler {
	return adaptor.HTTPHandler(playground.Handler("GraphQL playground", "/graphql"))
}

func main() {
	run := func(lc fx.Lifecycle, logger *slog.Logger) {
		app := fiber.New(fiber.Config{})

		app.All("/ping", NewPingHandler(logger))
		app.All("/graphql", NewGraphqlHandler())
		app.All("/playground", NewPlaygroundHandler())

		go func() {
			if err := app.Listen(":8080", fiber.ListenConfig{
				DisableStartupMessage: true,
			}); err != nil {
				panic(err)
			}
		}()
		lc.Append(fx.StopHook(app.ShutdownWithContext))
	}

	fx.New(
		fx.WithLogger(func(logger *slog.Logger) fxevent.Logger {
			return &fxevent.SlogLogger{Logger: logger}
		}),
		fx.StopTimeout(time.Minute),
		fx.Provide(NewLogger),
		fx.Invoke(run),
	).Run()

}
