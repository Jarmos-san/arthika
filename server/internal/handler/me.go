package handler

import (
	"context"
	"fmt"

	"github.com/Jarmos-san/arthika/server/internal/api"
	"github.com/Jarmos-san/arthika/server/internal/auth"
	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
)

// GetCurrentUser returns the authenticated user's ID and email.
//
// The auth middleware validates the session cookie (or Authorization header) and
// injects the user's claims into the context. This handler is a thin pass-
// through that reads those claims and returns them as a UserResponse.
//
// Returns:
//   - GetCurrentUser200JSONResponse with the user's ID and email.
//   - An error if the user ID cannot be parsed (mapped to 500 by the strict
//     server middleware).
func (h *Handler) GetCurrentUser(
	ctx context.Context,
	_ api.GetCurrentUserRequestObject,
) (api.GetCurrentUserResponseObject, error) {
	userID := auth.UserIDFromContext(ctx)
	email := auth.EmailFromContext(ctx)

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to parse user ID from context", "error", err)

		return nil, fmt.Errorf("parse user ID: %w", err)
	}

	return api.GetCurrentUser200JSONResponse{
		Id:    userUUID,
		Email: types.Email(email),
	}, nil
}
