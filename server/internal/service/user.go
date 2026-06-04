// Package service implements the business logic layer of the application.
//
// Services orchestrate domain operations and coordinate with repositories
// to fulfill use cases.
package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Jarmos-san/arthika/server/internal/domain"
	"github.com/Jarmos-san/arthika/server/internal/dto"
	"github.com/Jarmos-san/arthika/server/internal/repository"
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
		name, email, password string,
	) (dto.CreateUser, error)
}

// Service is the concrete implementation of UserService.
type Service struct {
	Repo repository.UserRepository
}

// NewUserService creates a new Service backed by the given repository.
func NewUserService(repo repository.UserRepository) *Service {
	return &Service{Repo: repo}
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
	name, email, password string,
) (dto.CreateUser, error) {
	userID := uuid.NewString()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return dto.CreateUser{}, fmt.Errorf("failed to hash password: %w", err)
	}

	user := domain.User{
		ID:           userID,
		Username:     name,
		Email:        email,
		PasswordHash: string(hash),
	}

	err = s.Repo.Create(ctx, user)
	if err != nil {
		return dto.CreateUser{}, fmt.Errorf("failed to save user: %w", err)
	}

	return dto.CreateUser{
		ID:           userID,
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
	}, nil
}
