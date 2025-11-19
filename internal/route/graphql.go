package route

import (
	"log/slog"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/ravilushqa/otelgqlgen"
	"github.com/scmtble/graphql-playground/internal/graphql"
	"github.com/scmtble/graphql-playground/internal/graphql/generated"
	"github.com/vektah/gqlparser/v2/ast"
)

type GraphqlHandler struct {
	logger *slog.Logger
}

func (h *GraphqlHandler) Pattern() string {
	return "/graphql"
}

func (h *GraphqlHandler) Handler() fiber.Handler {
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

func NewGraphqlHandler(logger *slog.Logger) *GraphqlHandler {
	return &GraphqlHandler{
		logger: logger,
	}
}
