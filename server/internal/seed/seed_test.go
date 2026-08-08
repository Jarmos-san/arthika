// Package seed_test contains black-box tests for the seed package.
//
// Unit tests use a mockQuerier implementing repository.Querier to isolate
// seeder logic from database access, while the integration test exercises the
// seeders against a real SQLite database.
package seed_test

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"path/filepath"
	"testing"

	"github.com/Jarmos-san/arthika/server/internal/repository"
	"github.com/Jarmos-san/arthika/server/internal/seed"
	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3" // Register SQLite3 driver.
	_ "github.com/golang-migrate/migrate/v4/source/file"      // Register file source.
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// testPlaintextPassword mirrors the hard-coded dev password in the seed
// package so tests can verify the generated hash.
const testPlaintextPassword = "dev-password"

// Sentinel errors simulating failures from the querier, defined at package
// level so tests can assert on wrapped errors.
var (
	errLookupUnavailable = errors.New("database unavailable")
	errCreateFailed      = errors.New("insert failed")
)

// mockQuerier implements repository.Querier with configurable function fields.
// Each test sets only the functions it needs; nil functions panic if called.
type mockQuerier struct {
	findUserByEmailFn      func(ctx context.Context, email string) (repository.User, error)
	createUserFn           func(ctx context.Context, arg repository.CreateUserParams) error
	countUsersFn           func(ctx context.Context) (int64, error)
	createAssetClassFn     func(ctx context.Context, arg repository.CreateAssetClassParams) error
	deleteAssetClassFn     func(ctx context.Context, id string) (string, error)
	findAssetClassByIDFn   func(ctx context.Context, id string) (repository.AssetClass, error)
	findAssetClassByNameFn func(ctx context.Context, name string) (repository.AssetClass, error)
	listAssetClassesFn     func(ctx context.Context) ([]repository.AssetClass, error)
	updateAssetClassFn     func(ctx context.Context, arg repository.UpdateAssetClassParams) (repository.AssetClass, error)
}

// FindUserByEmail delegates to m.findUserByEmailFn.
func (m *mockQuerier) FindUserByEmail(
	ctx context.Context,
	email string,
) (repository.User, error) {
	return m.findUserByEmailFn(ctx, email)
}

// CreateUser delegates to m.createUserFn.
func (m *mockQuerier) CreateUser(
	ctx context.Context,
	arg repository.CreateUserParams,
) error {
	return m.createUserFn(ctx, arg)
}

// CountUsers delegates to m.countUsersFn.
func (m *mockQuerier) CountUsers(ctx context.Context) (int64, error) {
	return m.countUsersFn(ctx)
}

// CreateAssetClass delegates to m.createAssetClassFn.
func (m *mockQuerier) CreateAssetClass(
	ctx context.Context,
	arg repository.CreateAssetClassParams,
) error {
	return m.createAssetClassFn(ctx, arg)
}

// DeleteAssetClass delegates to m.deleteAssetClassFn.
func (m *mockQuerier) DeleteAssetClass(ctx context.Context, id string) (string, error) {
	return m.deleteAssetClassFn(ctx, id)
}

// FindAssetClassByID delegates to m.findAssetClassByIDFn.
func (m *mockQuerier) FindAssetClassByID(
	ctx context.Context,
	id string,
) (repository.AssetClass, error) {
	return m.findAssetClassByIDFn(ctx, id)
}

// FindAssetClassByName delegates to m.findAssetClassByNameFn.
func (m *mockQuerier) FindAssetClassByName(
	ctx context.Context,
	name string,
) (repository.AssetClass, error) {
	return m.findAssetClassByNameFn(ctx, name)
}

// ListAssetClasses delegates to m.listAssetClassesFn.
func (m *mockQuerier) ListAssetClasses(ctx context.Context) ([]repository.AssetClass, error) {
	return m.listAssetClassesFn(ctx)
}

// UpdateAssetClass delegates to m.updateAssetClassFn.
func (m *mockQuerier) UpdateAssetClass(
	ctx context.Context,
	arg repository.UpdateAssetClassParams,
) (repository.AssetClass, error) {
	return m.updateAssetClassFn(ctx, arg)
}

// TestSeedDevUser_Creates verifies a missing dev user is created with the
// well-known ID, email and a bcrypt hash of the dev password.
func TestSeedDevUser_Creates(t *testing.T) {
	t.Parallel()

	var createdParams repository.CreateUserParams

	mock := &mockQuerier{
		findUserByEmailFn: func(_ context.Context, email string) (repository.User, error) {
			if email != seed.DevUserEmail {
				t.Errorf("expected email %s, got %s", seed.DevUserEmail, email)
			}

			return repository.User{}, sql.ErrNoRows
		},
		createUserFn: func(_ context.Context, arg repository.CreateUserParams) error {
			createdParams = arg

			return nil
		},
		countUsersFn:           nil,
		createAssetClassFn:     nil,
		deleteAssetClassFn:     nil,
		findAssetClassByIDFn:   nil,
		findAssetClassByNameFn: nil,
		listAssetClassesFn:     nil,
		updateAssetClassFn:     nil,
	}

	seeder := seed.New(mock, slog.Default())

	err := seeder.SeedDevUser(t.Context())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if createdParams.Email != seed.DevUserEmail {
		t.Errorf("expected email %s, got %s", seed.DevUserEmail, createdParams.Email)
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(createdParams.PasswordHash),
		[]byte(testPlaintextPassword),
	)
	if err != nil {
		t.Errorf("password hash does not match dev password: %v", err)
	}
}

// TestSeedDevUser_SkipsWhenExists verifies that an existing dev user is left
// untouched without creating a duplicate.
func TestSeedDevUser_SkipsWhenExists(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		findUserByEmailFn: func(_ context.Context, email string) (repository.User, error) {
			return repository.User{
				ID:           "",
				Email:        email,
				PasswordHash: "hash",
			}, nil
		},
		countUsersFn:           nil,
		createAssetClassFn:     nil,
		createUserFn:           nil, // Panics if called; seeding must skip.
		deleteAssetClassFn:     nil,
		findAssetClassByIDFn:   nil,
		findAssetClassByNameFn: nil,
		listAssetClassesFn:     nil,
		updateAssetClassFn:     nil,
	}

	seeder := seed.New(mock, slog.Default())

	err := seeder.SeedDevUser(t.Context())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestSeedDevUser_ErrorFromLookup verifies lookup errors besides ErrNoRows are
// propagated.
func TestSeedDevUser_ErrorFromLookup(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		findUserByEmailFn: func(_ context.Context, _ string) (repository.User, error) {
			return repository.User{}, errLookupUnavailable
		},
		countUsersFn:           nil,
		createAssetClassFn:     nil,
		createUserFn:           nil,
		deleteAssetClassFn:     nil,
		findAssetClassByIDFn:   nil,
		findAssetClassByNameFn: nil,
		listAssetClassesFn:     nil,
		updateAssetClassFn:     nil,
	}

	seeder := seed.New(mock, slog.Default())

	err := seeder.SeedDevUser(t.Context())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, errLookupUnavailable) {
		t.Errorf("expected error to wrap %v, got %v", errLookupUnavailable, err)
	}
}

// TestSeedDevUser_ErrorFromCreate verifies creation failures are propagated.
func TestSeedDevUser_ErrorFromCreate(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		findUserByEmailFn: func(_ context.Context, _ string) (repository.User, error) {
			return repository.User{}, sql.ErrNoRows
		},
		createUserFn: func(_ context.Context, _ repository.CreateUserParams) error {
			return errCreateFailed
		},
		countUsersFn:           nil,
		createAssetClassFn:     nil,
		deleteAssetClassFn:     nil,
		findAssetClassByIDFn:   nil,
		findAssetClassByNameFn: nil,
		listAssetClassesFn:     nil,
		updateAssetClassFn:     nil,
	}

	seeder := seed.New(mock, slog.Default())

	err := seeder.SeedDevUser(t.Context())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, errCreateFailed) {
		t.Errorf("expected error to wrap %v, got %v", errCreateFailed, err)
	}
}

// setupSeedTestDB creates a temporary SQLite database with all migrations
// applied and returns the raw queries bound to it.
func setupSeedTestDB(t *testing.T) *repository.Queries {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "test.db")

	database, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}

	t.Cleanup(func() {
		err := database.Close()
		if err != nil {
			t.Errorf("close database: %v", err)
		}
	})

	migrator, err := migrate.New("file://../../db/migrations", "sqlite3://"+dbPath)
	if err != nil {
		t.Fatalf("create migration instance: %v", err)
	}

	err = migrator.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		t.Fatalf("apply migrations: %v", err)
	}

	sourceErr, dbErr := migrator.Close()
	if sourceErr != nil {
		t.Errorf("close migration source: %v", sourceErr)
	}

	if dbErr != nil {
		t.Errorf("close migration database: %v", dbErr)
	}

	return repository.New(database)
}

// TestSeed_RealDatabase verifies seeding works end-to-end against a real
// SQLite database with migrations applied: exactly one dev user exists with
// the well-known ID and a valid password hash.
func TestSeed_RealDatabase(t *testing.T) {
	t.Parallel()

	queries := setupSeedTestDB(t)
	seeder := seed.New(queries, slog.New(slog.DiscardHandler))

	err := seeder.Seed(t.Context())
	if err != nil {
		t.Fatalf("seed database: %v", err)
	}

	count, err := queries.CountUsers(t.Context())
	if err != nil {
		t.Fatalf("count users: %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1 user after seeding, got %d", count)
	}

	user, err := queries.FindUserByEmail(t.Context(), seed.DevUserEmail)
	if err != nil {
		t.Fatalf("find dev user: %v", err)
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(testPlaintextPassword),
	)
	if err != nil {
		t.Errorf("password hash does not match dev password: %v", err)
	}
}

// TestSeed_IsIdempotent verifies re-running the seed adds no duplicate users.
func TestSeed_IsIdempotent(t *testing.T) {
	t.Parallel()

	queries := setupSeedTestDB(t)
	seeder := seed.New(queries, slog.New(slog.DiscardHandler))

	err := seeder.Seed(t.Context())
	if err != nil {
		t.Fatalf("seed database: %v", err)
	}

	err = seeder.Seed(t.Context())
	if err != nil {
		t.Fatalf("re-seed database: %v", err)
	}

	count, err := queries.CountUsers(t.Context())
	if err != nil {
		t.Fatalf("count users: %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1 user after re-seeding, got %d", count)
	}
}
