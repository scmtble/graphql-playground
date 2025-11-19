package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/scmtble/graphql-playground/internal/route"
	"github.com/scmtble/graphql-playground/internal/server"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{},
	))

	fx.New(
		fx.WithLogger(func(logger *slog.Logger) fxevent.Logger {
			return &fxevent.SlogLogger{Logger: logger}
		}),
		fx.Supply(logger),
		route.Module,
		server.Module,
		fx.StopTimeout(time.Minute),
	).Run()
}
