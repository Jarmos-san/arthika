package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Jarmos-san/arthika/server/internal/api"
	"github.com/Jarmos-san/arthika/server/internal/repository"
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

// duplicateAssetClassNameMessage is the 409 error message used when an asset
// class name is already in use by another class.
const duplicateAssetClassNameMessage = "asset class name already exists"

// errAssetClassNameTaken is returned when a create or update attempt uses an
// asset class name that already exists in the database.
var errAssetClassNameTaken = errors.New("asset class name already exists")

// assetClassFromRow converts a repository row into its API representation.
// Description is nullable in the database, so a NULL description maps to a nil
// pointer; ids are always server-generated UUIDs.
func assetClassFromRow(row repository.AssetClass) api.AssetClass {
	var description *string
	if row.Description.Valid {
		description = &row.Description.String
	}

	id, _ := uuid.Parse(row.ID) // ids are server-generated UUIDs

	return api.AssetClass{
		Id:          id,
		Name:        row.Name,
		Description: description,
	}
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

// createAssetClass checks that the name is not taken, then persists the new
// asset class and returns its API representation. It returns
// errAssetClassNameTaken when a class with the same name already exists.
func (h *Handler) createAssetClass(
	ctx context.Context,
	name, description string,
) (api.AssetClass, error) {
	// Check if an asset with the provide name already exists
	_, err := h.querier.FindAssetClassByName(ctx, name)
	if err == nil {
		return api.AssetClass{}, errAssetClassNameTaken
	}

	// Raise an error if the database query failed
	if !errors.Is(err, sql.ErrNoRows) {
		return api.AssetClass{}, fmt.Errorf(
			"find asset class by name: %w",
			err,
		)
	}

	// Create a new ID to be assigned to the asset class
	assetID := uuid.New()

	// Persist the new asset class and its info in to the database
	err = h.querier.CreateAssetClass(ctx, repository.CreateAssetClassParams{
		ID:          assetID.String(),
		Name:        name,
		Description: sql.NullString{String: description, Valid: description != ""},
	})
	if err != nil {
		return api.AssetClass{}, fmt.Errorf("create asset class: %w", err)
	}

	// Prepare a string pointer to be serialised
	var desc *string

	// If the asset class description in non-empty then create a string at the pointer
	if description != "" {
		d := description
		desc = &d
	}

	// Serialise the asset class and its related information
	return api.AssetClass{
		Id:          assetID,
		Name:        name,
		Description: desc,
	}, nil
}

// CreateAsset adds a new asset class to the database.
func (h *Handler) CreateAsset(
	ctx context.Context,
	req api.CreateAssetRequestObject,
) (api.CreateAssetResponseObject, error) {
	// Raise an error if the request body is empty
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

	// Prepare the name and the description for persistence
	name := strings.TrimSpace(req.Body.Name)
	description := trimOptional(req.Body.Description)

	// Validate the user inputs
	if errs := validateAssetClassInput(name); len(errs) > 0 {
		return api.CreateAsset422JSONResponse{Errors: errs}, nil
	}

	// Persist the asset class and it's information in to the database
	asset, err := h.createAssetClass(ctx, name, description)
	if errors.Is(err, errAssetClassNameTaken) {
		return api.CreateAsset409JSONResponse{Message: duplicateAssetClassNameMessage}, nil
	}

	// Raise an error if the persistence failed
	if err != nil {
		return nil, err
	}

	// Log a message if the persistence was successful
	h.logger.InfoContext(ctx, "asset class created", "id", asset.Id)

	// Return the serialised response if the asset class was persisted
	return api.CreateAsset201JSONResponse(asset), nil
}

// ListAssets returns all asset classes from the database in name order.
func (h *Handler) ListAssets(
	ctx context.Context,
	_ api.ListAssetsRequestObject,
) (api.ListAssetsResponseObject, error) {
	rows, err := h.querier.ListAssetClasses(ctx)
	if err != nil {
		return nil, fmt.Errorf("list asset classes: %w", err)
	}

	assets := make([]api.AssetClass, 0, len(rows))
	for _, row := range rows {
		assets = append(assets, assetClassFromRow(row))
	}

	return api.ListAssets200JSONResponse(assets), nil
}

// GetAsset returns a single asset class, or 404 if the id is unknown.
func (h *Handler) GetAsset(
	ctx context.Context,
	req api.GetAssetRequestObject,
) (api.GetAssetResponseObject, error) {
	asset, err := h.querier.FindAssetClassByID(ctx, req.Id.String())
	if errors.Is(err, sql.ErrNoRows) {
		return api.GetAsset404JSONResponse{Message: assetClassNotFoundMessage}, nil
	}

	if err != nil {
		return nil, fmt.Errorf("find asset class by id: %w", err)
	}

	return api.GetAsset200JSONResponse(assetClassFromRow(asset)), nil
}

// UpdateAsset replaces the name and description of an asset class.
func (h *Handler) UpdateAsset(
	ctx context.Context,
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

	_, err := h.querier.FindAssetClassByID(ctx, req.Id.String())
	if errors.Is(err, sql.ErrNoRows) {
		return api.UpdateAsset404JSONResponse{Message: assetClassNotFoundMessage}, nil
	}

	if err != nil {
		return nil, fmt.Errorf("find asset class by id: %w", err)
	}

	existing, err := h.querier.FindAssetClassByName(ctx, name)
	if err == nil && existing.ID != req.Id.String() {
		return api.UpdateAsset409JSONResponse{Message: duplicateAssetClassNameMessage}, nil
	}

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("find asset class by name: %w", err)
	}

	updated, err := h.querier.UpdateAssetClass(ctx, repository.UpdateAssetClassParams{
		ID:          req.Id.String(),
		Name:        name,
		Description: sql.NullString{String: description, Valid: description != ""},
	})
	if err != nil {
		return nil, fmt.Errorf("update asset class: %w", err)
	}

	return api.UpdateAsset200JSONResponse(assetClassFromRow(updated)), nil
}

// DeleteAsset removes an asset class from the database.
func (h *Handler) DeleteAsset(
	ctx context.Context,
	req api.DeleteAssetRequestObject,
) (api.DeleteAssetResponseObject, error) {
	_, err := h.querier.DeleteAssetClass(ctx, req.Id.String())
	if errors.Is(err, sql.ErrNoRows) {
		return api.DeleteAsset404JSONResponse{Message: assetClassNotFoundMessage}, nil
	}

	if err != nil {
		return nil, fmt.Errorf("delete asset class: %w", err)
	}

	return api.DeleteAsset204Response{}, nil
}
