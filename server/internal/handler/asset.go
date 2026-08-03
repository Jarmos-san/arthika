package handler

import (
	"context"
	"sort"
	"strings"
	"sync"

	"github.com/Jarmos-san/arthika/server/internal/api"
	"github.com/google/uuid"
)

// requestBodyField is the validation field name used when the request body is
// missing.
const requestBodyField = "body"

// requestBodyRequiredMessage is the validation message used when the request
// body is missing.
const requestBodyRequiredMessage = "request body is required"

// assetClassNotFoundMessage is the 404 error message used when an asset class
// id is unknown.
const assetClassNotFoundMessage = "asset class not found"

// assetClassStore is an in-memory store for asset classes.
//
// TODO(assets): Replace this type with SQL-backed repository methods
// (sqlc query in server/db/query/ + a method on repository.Querier) when
// persistence lands. Nothing outside asset.go depends on its internals.
//
//nolint:godox // Intentional marker; see plan for the persistence swap.
type assetClassStore struct {
	mu    sync.RWMutex
	items map[string]api.AssetClass
}

// Seed classes so the API returns data on first load, mirroring the
// client-side seeds in useAssetClasses.ts.
//
//nolint:gochecknoglobals // Fixed demo seed set, safe as package-level.
var seedAssetClasses = []struct{ name, description string }{
	{"Equities", "Public stocks, ETFs, and index funds"},
	{"Fixed income", "Bonds, treasuries, and cash equivalents"},
	{"Crypto", "Bitcoin, Ethereum, and other digital assets"},
}

// newAssetClassStore creates a store pre-seeded with the demo classes.
//
// TODO(assets): Drop the seeds when persistence lands and move them into a
// migration or seed script instead.
//
//nolint:godox // Intentional marker; see plan for the persistence swap.
func newAssetClassStore() *assetClassStore {
	//nolint:exhaustruct // mu is intentionally left as the zero value.
	store := &assetClassStore{
		items: make(map[string]api.AssetClass, len(seedAssetClasses)),
	}
	for _, seed := range seedAssetClasses {
		asset := api.AssetClass{
			Id:          uuid.New(),
			Name:        seed.name,
			Description: &seed.description,
		}
		store.items[asset.Id.String()] = asset
	}

	return store
}

// create adds a new asset class and returns it with a generated UUID.
func (s *assetClassStore) create(name, description string) api.AssetClass {
	asset := api.AssetClass{
		Id:          uuid.New(),
		Name:        name,
		Description: &description,
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[asset.Id.String()] = asset

	return asset
}

// list returns all asset classes sorted by name for a stable UI order.
func (s *assetClassStore) list() []api.AssetClass {
	s.mu.RLock()
	defer s.mu.RUnlock()

	assets := make([]api.AssetClass, 0, len(s.items))
	for _, asset := range s.items {
		assets = append(assets, asset)
	}

	sort.Slice(assets, func(i, j int) bool { return assets[i].Name < assets[j].Name })

	return assets
}

// get returns the asset class with the given id, if it exists.
func (s *assetClassStore) get(assetID string) (api.AssetClass, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	asset, ok := s.items[assetID]

	return asset, ok
}

// update replaces the name and description of the asset class with the given
// id. The second return value reports whether the id existed.
func (s *assetClassStore) update(assetID, name, description string) (
	api.AssetClass,
	bool,
) {
	s.mu.Lock()
	defer s.mu.Unlock()

	asset, ok := s.items[assetID]
	if !ok {
		//nolint:exhaustruct // Zero value; unused by callers when ok is false.
		return api.AssetClass{}, false
	}

	asset.Name = name
	asset.Description = &description
	s.items[assetID] = asset

	return asset, true
}

// delete removes the asset class with the given id. The return value reports
// whether the id existed.
func (s *assetClassStore) delete(assetID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[assetID]; !ok {
		return false
	}

	delete(s.items, assetID)

	return true
}

// trimOptional returns the trimmed value of an optional string field, or an
// empty string when the field is absent.
func trimOptional(s *string) string {
	if s == nil {
		return ""
	}

	return strings.TrimSpace(*s)
}

// validateAssetClassInput mirrors the client-side rule in
// validateAssetClass (client/app/utils/validators.ts): name is required
// and must be non-blank.
func validateAssetClassInput(name string) []api.ValidationError {
	if strings.TrimSpace(name) == "" {
		return []api.ValidationError{
			{
				Field:   "name",
				Message: "name is required",
			},
		}
	}

	return nil
}

// CreateAsset adds a new asset class to the in-memory store.
//
// TODO(assets): Persist via a repository method once the DB layer exists.
//
//nolint:godox // Intentional marker; see plan for the persistence swap.
func (h *Handler) CreateAsset(
	ctx context.Context,
	req api.CreateAssetRequestObject,
) (api.CreateAssetResponseObject, error) {
	if req.Body == nil {
		return api.CreateAsset422JSONResponse{
			Errors: []api.ValidationError{
				{
					Field:   requestBodyField,
					Message: requestBodyRequiredMessage,
				},
			},
		}, nil
	}

	name := strings.TrimSpace(req.Body.Name)
	description := trimOptional(req.Body.Description)

	if errs := validateAssetClassInput(name); len(errs) > 0 {
		return api.CreateAsset422JSONResponse{Errors: errs}, nil
	}

	asset := h.assetClasses.create(name, description)
	h.logger.InfoContext(ctx, "asset class created", "id", asset.Id)

	return api.CreateAsset201JSONResponse(asset), nil
}

// ListAssets returns all asset classes from the in-memory store.
func (h *Handler) ListAssets(
	_ context.Context,
	_ api.ListAssetsRequestObject,
) (api.ListAssetsResponseObject, error) {
	return api.ListAssets200JSONResponse(h.assetClasses.list()), nil
}

// GetAsset returns a single asset class, or 404 if the id is unknown.
func (h *Handler) GetAsset(
	_ context.Context,
	req api.GetAssetRequestObject,
) (api.GetAssetResponseObject, error) {
	asset, ok := h.assetClasses.get(req.Id.String())
	if !ok {
		return api.GetAsset404JSONResponse{Message: assetClassNotFoundMessage}, nil
	}

	return api.GetAsset200JSONResponse(asset), nil
}

// UpdateAsset replaces the name and description of an asset class.
//
// TODO(assets): Persist via a repository method once the DB layer exists.
//
//nolint:godox // Intentional marker; see plan for the persistence swap.
func (h *Handler) UpdateAsset(
	_ context.Context,
	req api.UpdateAssetRequestObject,
) (api.UpdateAssetResponseObject, error) {
	if req.Body == nil {
		return api.UpdateAsset422JSONResponse{
			Errors: []api.ValidationError{
				{
					Field:   requestBodyField,
					Message: requestBodyRequiredMessage,
				},
			},
		}, nil
	}

	name := strings.TrimSpace(req.Body.Name)
	description := trimOptional(req.Body.Description)

	if errs := validateAssetClassInput(name); len(errs) > 0 {
		return api.UpdateAsset422JSONResponse{Errors: errs}, nil
	}

	asset, ok := h.assetClasses.update(req.Id.String(), name, description)
	if !ok {
		return api.UpdateAsset404JSONResponse{Message: assetClassNotFoundMessage}, nil
	}

	return api.UpdateAsset200JSONResponse(asset), nil
}

// DeleteAsset removes an asset class from the in-memory store.
//
// TODO(assets): Decide whether deletion should be blocked (409) or cascade
// once holdings reference asset classes; persist via a repository method.
//
//nolint:godox // Intentional marker; see plan for the persistence swap.
func (h *Handler) DeleteAsset(
	_ context.Context,
	req api.DeleteAssetRequestObject,
) (api.DeleteAssetResponseObject, error) {
	if !h.assetClasses.delete(req.Id.String()) {
		return api.DeleteAsset404JSONResponse{Message: assetClassNotFoundMessage}, nil
	}

	return api.DeleteAsset204Response{}, nil
}
