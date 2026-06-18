package handler

import (
	"context"
	"errors"

	"github.com/Jarmos-san/arthika/server/internal/api"
)

func (h *Handler) Register(
	ctx context.Context,
	_ api.RegisterRequestObject,
) (api.RegisterResponseObject, error) {
	return nil, errors.New("not implemented")
}
