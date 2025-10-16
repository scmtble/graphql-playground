package main

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/scmtble/graphql-playground/graphql"
	"github.com/scmtble/graphql-playground/graphql/generated"
	"github.com/vektah/gqlparser/v2/ast"
)

func main() {
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

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/playground", playground.ApolloSandboxHandler("GraphQL playground", "/graphql"))
	http.Handle("/graphql", srv)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
