// Package `main` is the entry point of the application.
//
// It is responsible for initializing configuration, setting up HTTP routing,
// constructing the application container and managing the server lifecycle including
// graceful shutdown on system signals.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Jarmos-san/arthika/server/internal/api"
	"github.com/Jarmos-san/arthika/server/internal/app"
	"github.com/Jarmos-san/arthika/server/internal/config"
	"github.com/Jarmos-san/arthika/server/internal/handler"
	"github.com/Jarmos-san/arthika/server/internal/logger"
	"github.com/Jarmos-san/arthika/server/internal/repository"
	chi "github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

// shutdownTimeout defines the maximum duration allowed for gracefully shutting
// down the HTTP server.
//
// During shutdown, the server stops accepting new connections and waits for in-flight
// requests to complete. If this timeout is exceeded, the shutdown process is aborted
// and any remaining connections may be terminated.
const shutdownTimeout = 5 * time.Second

// `main` initialises and runs the HTTP server.
//
// It performs the following steps:
//   - Loads configurations from environment variables.
//   - Sets up HTTP routes.
//   - Constructs the application container.
//   - Starts the server in a separate goroutine.
//   - Listens for OS signals to trigger graceful shutdown.
//
// The server is gracefully shutdown when an interrupt or termination signal is
// received, allowing in-flight requests to complete wihin a timeout period.
func main() {
	// Load application configuration from environment variables.
	cfg := config.LoadConfig()

	logger := logger.New(cfg.LogLevel)

	// Initialise the Chi router.
	router := chi.NewRouter()
	router.Use(chimw.Logger, chimw.Recoverer)

	// Construct the app container (opens DB, runs migrations).
	app, err := app.New(cfg, router, logger)
	if err != nil {
		logger.Error(
			"failed to initialize application",
			slog.String("error", err.Error()),
		)

		return
	}

	queries := repository.New(app.DB)
	h := handler.NewHandler(logger, queries)
	strictHandler := api.NewStrictHandler(h, nil)

	docsHandler := handler.NewDocsHandler(logger)
	router.Get("/docs", docsHandler.GetDocsPage)
	router.Get("/openapi.json", docsHandler.GetSpecJSON)

	api.HandlerFromMuxWithBaseURL(strictHandler, router, "/api")

	// Create a context that is cancelled on interrupt or kill signals.
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGKILL,
	)
	defer stop()

	// Start the HTTP server in a separate goroutine. This allows the `main` goroutine
	// to listen for shutdown signals
	go func() {
		logger.Info("starting server")

		err := app.Run()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server startup failed", "error", err.Error())
		}
	}()

	// Block until a shutdown signal is received.
	<-ctx.Done()

	// Create a timeout-bound context for graceful shutdown.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Attempt a graceful shutdown of the server.
	err = app.Shutdown(shutdownCtx)
	if err != nil {
		logger.Error("server shutdown failed", "error", err.Error())
	}

	logger.Info("server shutdown completed gracefully")
}
