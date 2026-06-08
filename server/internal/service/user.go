// Package service implements the business logic layer of the application.
//
// Services orchestrate domain operations and coordinate with repositories
// to fulfill use cases.
package service

import (
	"context"
	"fmt"

	"github.com/Jarmos-san/arthika/server/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user entity returned by the service layer.
type User struct {
	Name string `json:"name"`
}

// UserService defines the business operations available for users.
type UserService interface {
	GetUser() (User, error)
	CreateUser(
		ctx context.Context,
		username, email, password string,
	) (repository.User, error)
}

// Service is the concrete implementation of UserService.
type Service struct {
	q repository.Querier
}

// NewUserService creates a new Service backed by the given query interface.
func NewUserService(q repository.Querier) *Service {
	return &Service{q: q}
}

// GetUser returns a static user.
//
// This is a temporary stub until the service is fully wired.
func (s *Service) GetUser() (User, error) {
	return User{
		Name: "John Doe",
	}, nil
}

// CreateUser registers a new user and returns the created user DTO.
func (s *Service) CreateUser(
	ctx context.Context,
	username, email, password string,
) (repository.User, error) {
	userID := uuid.NewString()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return repository.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	params := repository.CreateUserParams{
		ID:           userID,
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
	}

	err = s.q.CreateUser(ctx, params)
	if err != nil {
		return repository.User{}, fmt.Errorf("failed to save user: %w", err)
	}

	return repository.User{
		ID:           userID,
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
	}, nil
}
