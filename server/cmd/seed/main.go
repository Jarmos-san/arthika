// Package `main` is the entry point of the seeding executable.
//
// It populates the database referenced by the DATABASE_URL environment variable
// (see the server configuration for defaults) with development data. Schema
// migrations are expected to have been applied first, e.g. via `task migrate:up`.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/Jarmos-san/arthika/server/internal/config"
	"github.com/Jarmos-san/arthika/server/internal/logger"
	"github.com/Jarmos-san/arthika/server/internal/repository"
	"github.com/Jarmos-san/arthika/server/internal/seed"
	_ "github.com/mattn/go-sqlite3"
)

// `main` loads the configuration and runs all registered seeders.
func main() {
	// Load the configs and the logger instance
	cfg := config.LoadConfig()
	logger := logger.New(cfg.LogLevel)

	//  Run the seeding logic
	err := run(cfg, logger)
	if err != nil {
		logger.Error("seeding failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Return a "success" message if the db seeding was completed
	logger.Info("seeding complete")
}

// run opens the database, seeds development data into it and closes the
// database before returning.
func run(cfg config.Config, logger *slog.Logger) error {
	// Attempt to open the SQLite3 database file
	database, err := sql.Open("sqlite3", cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	// Close the connection to the database after the execution goes out of scope
	defer func() {
		err := database.Close()
		if err != nil {
			logger.Error("failed to close database", slog.String("error", err.Error()))
		}
	}()

	// Check if the database is accessible
	err = database.PingContext(context.Background())
	if err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	// Prepare the SQL statements to populate the database
	queries := repository.New(database)
	seeder := seed.New(queries, logger)

	// Run the SQL statements prepared above
	err = seeder.Seed(context.Background())
	if err != nil {
		return fmt.Errorf("seed database: %w", err)
	}

	return nil
}
