// Package service implements the business logic layer of the application.
//
// Services orchestrate domain operations and coordinate with repositories
// to fulfill use cases.
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Jarmos-san/arthika/server/api"
	"github.com/Jarmos-san/arthika/server/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Service implements user-related business logic.
type Service struct {
	q           repository.Querier
	tokenSecret string
}

// NewUserService constructs a Service with its required dependencies.
func NewUserService(q repository.Querier, tokenSecret string) *Service {
	return &Service{q: q, tokenSecret: tokenSecret}
}

// GetUser returns a static user stub.
func (s *Service) GetUser() api.User {
	return api.User{
		Name: "John Doe",
	}
}

// CreateUser registers a new user and returns the created user.
func (s *Service) CreateUser(
	ctx context.Context,
	username, email, password string,
) (api.User, error) {
	userID := uuid.NewString()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return api.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	params := repository.CreateUserParams{
		ID:           userID,
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
	}

	err = s.q.CreateUser(ctx, params)
	if err != nil {
		return api.User{}, fmt.Errorf("failed to save user: %w", err)
	}

	return api.User{
		Name: username,
	}, nil
}

// Login authenticates a user and returns a signed JWT token response.
func (s *Service) Login(
	ctx context.Context,
	email, password string,
) (api.TokenResponse, error) {
	user, err := s.q.FindUserByEmail(ctx, email)
	if err != nil {
		return api.TokenResponse{}, fmt.Errorf("invalid credentials: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash), []byte(password),
	)
	if err != nil {
		return api.TokenResponse{}, fmt.Errorf("invalid credentials: %w", err)
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub": user.ID,
		"iat": now.Unix(),
		"exp": now.Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedToken, err := token.SignedString([]byte(s.tokenSecret))
	if err != nil {
		return api.TokenResponse{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return api.TokenResponse{
		AccessToken: signedToken,
		TokenType:   "Bearer",
		ExpiresIn:   float32(time.Hour.Seconds()),
	}, nil
}
