package handler

import (
	"log/slog"

	"github.com/Jarmos-san/arthika/server/internal/api"
	"github.com/Jarmos-san/arthika/server/internal/repository"
)

var _ api.StrictServerInterface = (*Handler)(nil)

// Handler serves HTTP requests for the API.
type Handler struct {
	logger  *slog.Logger
	querier repository.Querier

	// assetClasses is the in-memory asset class store.
	//nolint:godox // Intentional marker; see plan for the persistence swap.
	// TODO(assets): Replace with a SQL-backed repository when persistence
	// lands. The handlers only depend on assetClassStore, so the swap is
	// localised to asset.go and newAssetClassStore.
	assetClasses *assetClassStore
}

// NewHandler creates a new Handler.
func NewHandler(logger *slog.Logger, querier repository.Querier) *Handler {
	return &Handler{
		logger:       logger,
		querier:      querier,
		assetClasses: newAssetClassStore(),
	}
}
