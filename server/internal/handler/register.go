package handler

import (
	"context"
	"errors"

	"github.com/Jarmos-san/arthika/server/internal/api"
)

// ErrNotImplemented is returned when an endpoint has not been implemented yet.
var ErrNotImplemented = errors.New("not implemented")

// Register handles user registration.
func (h *Handler) Register(
	_ context.Context,
	_ api.RegisterRequestObject,
) (api.RegisterResponseObject, error) {
	return nil, ErrNotImplemented
}
