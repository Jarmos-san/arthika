// Package app provides the core application container responsible for wiring
// configuration, HTTP handlers and server lifecycle management.
//
// It acts as the composition root of the service, coordinating dependencies and
// exposing methods to start and gracefully shutdown the HTTP server.
package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Jarmos-san/arthika/server/internal/config"
	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3" // Register SQLite3 driver.
	_ "github.com/golang-migrate/migrate/v4/source/file"      // Register file source.
)

// Application represents the main Application container.
//
// It holds the runtime configuration, the HTTP server instance and the root HTTP
// handler, and the structured logger. This type is responsible for managing the
// lifecycle of the HTTP server and coordinating cross-cutting concerns such as logging.
type Application struct {
	// Config contains the application configuration values.
	Config config.Config

	// Server is the HTTP server responsible for handling incoming requests.
	Server *http.Server

	// Handler is the root HTTP handler used by the server to route requests.
	Handler http.Handler

	// DB is the SQLite database connection used by the application.
	DB *sql.DB

	// Logger is the structured logger used for emitting application-level logs.
	//
	// It is expected to be initialised by the caller and injected into the application.
	// The logger is used for lifecycle events such as startup and shutdown, as well as
	// other operational logging.
	Logger *slog.Logger
}

// New constructs and returns a new application instance.
//
// It initialises an `http.Server` using the provided configuration and handler, opens
// a SQLite database connection, and associates a structured logger for
// application-level
// logging.
//
// The returned application is ready to be started via the `Run` method.
func New(
	cfg config.Config,
	handler http.Handler,
	logger *slog.Logger,
) (*Application, error) {
	server := &http.Server{ //nolint:exhaustruct
		Addr:         cfg.Addr,
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	database, err := sql.Open("sqlite3", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	err = database.PingContext(context.Background())
	if err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	if cfg.MigrationDirectory != "" {
		err = runMigrations(cfg.DatabaseURL, cfg.MigrationDirectory)
		if err != nil {
			return nil, fmt.Errorf("run migrations: %w", err)
		}
	}

	return &Application{
		Config:  cfg,
		Server:  server,
		Handler: handler,
		DB:      database,
		Logger:  logger,
	}, nil
}

// Run starts the HTTP server and begins listening for incoming requests.
//
// It logs the server start event and blocks until the server stops. Any error returned
// by `ListenAndServer` is propagated to the caller.
func (a *Application) Run() error {
	a.Logger.Info("server started", "addr", a.Config.Addr)

	err := a.Server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("listen and serve: %w", err)
	}

	return nil
}

// runMigrations applies all pending SQL migrations using golang-migrate.
//
// It connects directly to the SQLite database using the provided URL and reads
// migration files from the specified directory. If no migrations are pending,
// ErrNoChange is returned and treated as success.
func runMigrations(databaseURL, migrationDir string) error {
	// Attempt to create an object to run the migration
	migrator, err := migrate.New("file://"+migrationDir, "sqlite3://"+databaseURL)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}

	// Defer closing the migration object post-work
	defer func() {
		// Log a warning if the migration source throws an error
		sourceErr, dbErr := migrator.Close()
		if sourceErr != nil {
			slog.Warn("migration source close failed", "error", sourceErr)
		}

		// Log a warning if there was a database issue during the migration phase
		if dbErr != nil {
			slog.Warn("migration database close failed", "error", dbErr)
		}
	}()

	// Close the migration object or return an error message
	err = migrator.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations: %w", err)
	}

	return nil
}

// Shutdown gracefully stops the HTTP server and closes the database connection.
//
// It attempts to shutdown the server using the provided context, allowing in-flight
// requests to complete before termination. If the context expires before shutdown
// completes, an error is returned.
func (a *Application) Shutdown(ctx context.Context) error {
	a.Logger.Info("server shutdown")

	err := a.Server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}

	err = a.DB.Close()
	if err != nil {
		return fmt.Errorf("close database: %w", err)
	}

	return nil
}
