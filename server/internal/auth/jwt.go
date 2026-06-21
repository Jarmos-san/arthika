// Package auth provides JWT token generation, validation, and context
// helpers shared across handlers and middleware.
package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// errInvalidSigningMethod is returned when a token uses an unexpected signing
// algorithm.
var errInvalidSigningMethod = errors.New("unexpected signing method")

// Claims holds the custom claims embedded in JWTs issued by the application.
type Claims struct {
	jwt.RegisteredClaims

	Email string `json:"email"`
}

// tokenExpiry is the lifetime of issued JWTs.
const tokenExpiry = 24 * time.Hour

// GenerateToken creates a signed JWT for the given user.
func GenerateToken(userID, email, secret string) (string, error) {
	// The current timestamp (required when generating the JWT)
	now := time.Now()

	// An object representing the JWT claims
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Arthika API",
			Subject:   userID,
			Audience:  jwt.ClaimStrings{},
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenExpiry)),
			NotBefore: nil,
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.NewString(),
		},
		Email: email,
	}

	// Generate the JWT using the claims created above and the HMAC-SHA256 algorithmj
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the generated token using a secret token passed to the software during
	// initial setup
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signedToken, nil
}

// ValidateToken parses and validates a signed JWT, returning the extracted claims.
func ValidateToken(tokenString, secret string) (*Claims, error) {
	// Create an object representing the JWT claims
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "",
			Subject:   "",
			Audience:  jwt.ClaimStrings{},
			ExpiresAt: nil,
			NotBefore: nil,
			IssuedAt:  nil,
			ID:        "",
		},
		Email: "",
	}

	// Parse the JWT and validate its integrity or throw an error if otherwise
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		// Validate the signing algorithm used to generate the JWT or throw an error if
		// it is not valid HMAC-SHA26 algorithm
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errInvalidSigningMethod
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("validate token: %w", err)
	}

	return claims, nil
}

// contextKey is an unexported type used to store auth claims in context.
type contextKey struct{}

// claimsContext holds auth data stored in the request context.
type claimsContext struct {
	userID string
	email  string
}

// NewContext embeds the user ID and email into the provided context.
func NewContext(ctx context.Context, userID, email string) context.Context {
	return context.WithValue(
		ctx,
		contextKey{},
		claimsContext{userID: userID, email: email},
	)
}

// UserIDFromContext extracts the authenticated user's ID from the context.
// Returns empty string if the context was not set by the auth middleware.
func UserIDFromContext(ctx context.Context) string {
	if c, ok := ctx.Value(contextKey{}).(claimsContext); ok {
		return c.userID
	}

	return ""
}

// EmailFromContext extracts the authenticated user's email from the context.
// Returns empty string if the context was not set by the auth middleware.
func EmailFromContext(ctx context.Context) string {
	if c, ok := ctx.Value(contextKey{}).(claimsContext); ok {
		return c.email
	}

	return ""
}
