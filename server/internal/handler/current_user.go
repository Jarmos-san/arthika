package handler

import (
	"context"
	"fmt"

	"github.com/Jarmos-san/arthika/server/internal/api"
	"github.com/Jarmos-san/arthika/server/internal/auth"
	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
)

// CurrentUser returns the authenticated user's ID and email.
//
// The auth middleware validates the HttpOnly cookie and injects the user's ID
// and email into the request context. This handler reads those values and
// returns them — no database query is required.
//
// Returns:
//   - CurrentUser200JSONResponse with the user's ID and email (200 OK).
//   - An unwrapped error if the user ID from context cannot be parsed (500).
func (h *Handler) CurrentUser(
	ctx context.Context,
	_ api.CurrentUserRequestObject,
) (api.CurrentUserResponseObject, error) {
	userID := auth.UserIDFromContext(ctx)
	email := auth.EmailFromContext(ctx)

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to parse user ID as UUID", "error", err)

		return nil, fmt.Errorf("parse user ID: %w", err)
	}

	return api.CurrentUser200JSONResponse{
		Id:    userUUID,
		Email: types.Email(email),
	}, nil
}
