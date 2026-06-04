// Package repository provides data access interfaces and implementations.
//
// It follows the Repository pattern to abstract storage from business logic.
package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Jarmos-san/arthika/server/internal/domain"
)

// UserRepository defines persistence operations for users.
type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
}

// UserRepo is the concrete SQLite-backed implementation of UserRepository.
type UserRepo struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepo backed by a SQLite database.
func NewUserRepository(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

// Create inserts a new user into the database.
//
// It returns an error if the insert fails (e.g. duplicate username or email).
func (r *UserRepo) Create(ctx context.Context, user domain.User) error {
	query := `INSERT INTO users (id, username, email, password_hash) VALUES (?, ?, ?, ?)`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
	)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	return nil
}

// FindByEmail retrieves a user by their email address.
//
// It returns the user or an error if the email is not found.
func (r *UserRepo) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	query := `SELECT id, username, email, password_hash FROM users WHERE email = ?`

	row := r.db.QueryRowContext(ctx, query, email)

	var user domain.User

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	if err != nil {
		return domain.User{}, fmt.Errorf("find user by email: %w", err)
	}

	return user, nil
}
