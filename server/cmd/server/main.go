// Package main is the entry point of the application.
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

	"github.com/Jarmos-san/arthika/server/api"
	"github.com/Jarmos-san/arthika/server/internal/app"
	"github.com/Jarmos-san/arthika/server/internal/config"
	"github.com/Jarmos-san/arthika/server/internal/handler"
	"github.com/Jarmos-san/arthika/server/internal/logger"
	"github.com/Jarmos-san/arthika/server/internal/repository"
	"github.com/Jarmos-san/arthika/server/internal/service"
	chi "github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

const shutdownTimeout = 5 * time.Second

func main() {
	cfg := config.LoadConfig()

	logger := logger.New(cfg.LogLevel)

	router := chi.NewRouter()
	router.Use(chimw.Logger, chimw.Recoverer)

	app, err := app.New(cfg, router, logger)
	if err != nil {
		logger.Error(
			"failed to initialize application",
			slog.String("error", err.Error()),
		)

		return
	}

	queries := repository.New(app.DB)
	userService := service.NewUserService(queries, cfg.TokenSecret)
	userHandler := handler.NewHandler(userService, logger)

	h := api.HandlerFromMux(userHandler, router)

	app.Server.Handler = h

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGKILL,
	)
	defer stop()

	go func() {
		logger.Info("starting server")

		err := app.Run()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server startup failed", "error", err.Error())
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = app.Shutdown(shutdownCtx)
	if err != nil {
		logger.Error("server shutdown failed", "error", err.Error())
	}

	logger.Info("server shutdown completed gracefully")
}
