// Package seed provides programmatic database seeding for development.
//
// Seeding is kept separate from schema migrations: migrations manage the
// database structure while this package populates the database with
// development-friendly data. Every seeder is idempotent so the executable can
// be re-run safely against an existing database.
//
// New seeders should be added as methods on Seeder and invoked from Seed().
package seed

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Jarmos-san/arthika/server/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	// DevUserEmail is the login email of the development user.
	DevUserEmail = "dev@example.com"

	// devPassword is the plaintext password for the development user.
	//
	// It is hard-coded for local development only and will be made configurable
	// in the future.
	devPassword = "dev-password"
)

// Seeder seeds development data into the database.
//
// It wraps the repository.Querier through which all database access happens.
type Seeder struct {
	queries repository.Querier
	logger  *slog.Logger
}

// New constructs a Seeder ready to populate the database.
//
// The provided querier is used for all database operations and the logger for
// reporting what was seeded.
func New(queries repository.Querier, logger *slog.Logger) *Seeder {
	return &Seeder{
		queries: queries,
		logger:  logger,
	}
}

// Seed runs all registered seeders in order.
//
// Future seeders (e.g., demo transactions) should be invoked from this method
// alongside the existing SeedDevUser and SeedAssetClasses calls.
func (s *Seeder) Seed(ctx context.Context) error {
	// Seed the database with a "dev user" for the development environment
	err := s.SeedDevUser(ctx)
	if err != nil {
		return fmt.Errorf("seed dev user: %w", err)
	}

	// Seed the database with the default asset classes
	err = s.SeedAssetClasses(ctx)
	if err != nil {
		return fmt.Errorf("seed asset classes: %w", err)
	}

	return nil
}

// SeedDevUser creates the development user if one with DevUserEmail does not
// already exist.
//
// It is idempotent: an existing user is left untouched and logged as skipped.
func (s *Seeder) SeedDevUser(ctx context.Context) error {
	// Make a fetch query in the database with the "dev user" email address
	_, err := s.queries.FindUserByEmail(ctx, DevUserEmail)

	// Check if the "dev user" already exists in the database, return an error otherwise
	switch {
	case err == nil:
		s.logger.Info(
			"dev user already exists, skipping",
			"email",
			DevUserEmail,
		)

		return nil
	case !errors.Is(err, sql.ErrNoRows):
		return fmt.Errorf("check for existing dev user: %w", err)
	}

	// Create a hashed password for the "dev user"
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(devPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf("hash dev password: %w", err)
	}

	// Add the "dev user" to the database
	err = s.queries.CreateUser(ctx, repository.CreateUserParams{
		ID:           uuid.NewString(),
		Email:        DevUserEmail,
		PasswordHash: string(passwordHash),
	})
	if err != nil {
		return fmt.Errorf("create dev user: %w", err)
	}

	// Log the "success" message
	s.logger.Info("dev user created", "email", DevUserEmail)

	return nil
}

// SeedAssetClasses creates the default asset classes if they do not already
// exist.
//
// It is idempotent: an existing class is left untouched and logged as skipped.
func (s *Seeder) SeedAssetClasses(ctx context.Context) error {
	assetClasses := []struct {
		name,
		description string
	}{
		{
			"Equity",
			"Ownership stakes in companies, e.g. stocks and shares.",
		},
		{
			"Debt Bonds",
			"Fixed-income securities representing loans to issuers.",
		},
		{
			"Commodities",
			"Raw materials and primary goods, e.g. gold, oil, wheat.",
		},
		{
			"Mutual Funds",
			"Pooled investment vehicles managed by professionals.",
		},
	}

	for _, assetClass := range assetClasses {
		// Check if the asset class already exists in the database
		_, err := s.queries.FindAssetClassByName(
			ctx,
			assetClass.name,
		)

		switch {
		case err == nil:
			s.logger.Info(
				"asset class already exists, skipping",
				"name",
				assetClass.name,
			)

			continue
		case !errors.Is(err, sql.ErrNoRows):
			return fmt.Errorf(
				"check for existing asset class %q: %w",
				assetClass.name,
				err,
			)
		}

		// Add the asset class to the database
		err = s.queries.CreateAssetClass(ctx, repository.CreateAssetClassParams{
			ID:          uuid.NewString(),
			Name:        assetClass.name,
			Description: sql.NullString{String: assetClass.description, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("create asset class %q: %w", assetClass.name, err)
		}

		// Log the "success" message
		s.logger.Info("asset class created", "name", assetClass.name)
	}

	return nil
}
