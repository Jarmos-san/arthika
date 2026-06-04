// Package repository provides data access interfaces and implementations.
//
// It follows the Repository pattern to abstract storage from business logic.
package repository

import (
	"context"
	"fmt"

	"github.com/Jarmos-san/arthika/server/internal/db"
)

// UserRepository defines persistence operations for users.
type UserRepository interface {
	Create(ctx context.Context, user db.User) error
	FindByEmail(ctx context.Context, email string) (db.User, error)
}

// UserRepo is the concrete SQLite-backed implementation of UserRepository.
type UserRepo struct {
	q *db.Queries
}

// NewUserRepository creates a new UserRepo backed by a SQLite database.
func NewUserRepository(dbtx db.DBTX) *UserRepo {
	return &UserRepo{q: db.New(dbtx)}
}

// Create inserts a new user into the database.
//
// It returns an error if the insert fails (e.g. duplicate username or email).
func (r *UserRepo) Create(ctx context.Context, user db.User) error {
	err := r.q.CreateUser(ctx, db.CreateUserParams(user))
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	return nil
}

// FindByEmail retrieves a user by their email address.
//
// It returns the user or an error if the email is not found.
func (r *UserRepo) FindByEmail(ctx context.Context, email string) (db.User, error) {
	user, err := r.q.FindUserByEmail(ctx, email)
	if err != nil {
		return db.User{}, fmt.Errorf("find user by email: %w", err)
	}

	return user, nil
}
