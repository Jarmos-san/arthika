package handler_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Jarmos-san/arthika/server/internal/api"
	"github.com/Jarmos-san/arthika/server/internal/handler"
	"github.com/Jarmos-san/arthika/server/internal/repository"
	chi "github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const (
	testAssetName     = "Real Estate"
	testAssetDesc     = "Property investments"
	unknownAssetID    = "00000000-0000-0000-0000-000000000000"
	dupNameErrMessage = "asset class name already exists"
)

// errTestDB is the error returned by mock repository functions to simulate an
// internal database failure.
var errTestDB = errors.New("database unavailable")

// assetMock returns a mockQuerier with every function field set to nil. Tests
// override only the fields they exercise.
func assetMock() *mockQuerier {
	return &mockQuerier{
		createUserFn:           nil,
		findUserByEmailFn:      nil,
		countUsersFn:           nil,
		createAssetClassFn:     nil,
		deleteAssetClassFn:     nil,
		findAssetClassByIDFn:   nil,
		findAssetClassByNameFn: nil,
		listAssetClassesFn:     nil,
		updateAssetClassFn:     nil,
	}
}

// assetRow builds a repository row with a nullable description.
func assetRow(id, name, description string) repository.AssetClass {
	return repository.AssetClass{
		ID:          id,
		Name:        name,
		Description: sql.NullString{String: description, Valid: description != ""},
	}
}

// newAssetRouter mounts a fresh Handler on a chi router the same way main.go
// does.
func newAssetRouter(t *testing.T, mock *mockQuerier) http.Handler {
	t.Helper()

	hdl := handler.NewHandler(slog.Default(), mock)
	router := chi.NewRouter()
	api.HandlerFromMuxWithBaseURL(api.NewStrictHandler(hdl, nil), router, "/api")

	return router
}

// doJSONRequest performs an HTTP request against the router and returns the
// recorded response.
func doJSONRequest(
	t *testing.T,
	router http.Handler,
	method, path, body string,
) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequestWithContext(
		t.Context(),
		method,
		path,
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	return rec
}

// createAsset posts a new asset class with the default test description and
// returns the recorded response.
func createAsset(t *testing.T, router http.Handler, name string) *httptest.ResponseRecorder {
	t.Helper()

	body := fmt.Sprintf(`{"name":%q,"description":%q}`, name, testAssetDesc)

	return doJSONRequest(t, router, http.MethodPost, "/api/assets", body)
}

// decodeAsset decodes an api.AssetClass from the recorded response body.
func decodeAsset(t *testing.T, rec *httptest.ResponseRecorder) api.AssetClass {
	t.Helper()

	var asset api.AssetClass

	err := json.NewDecoder(rec.Body).Decode(&asset)
	if err != nil {
		t.Fatalf("decode response: %v", err)
	}

	return asset
}

// assetDescription returns the description with nil collapsed to "".
func assetDescription(asset api.AssetClass) string {
	if asset.Description == nil {
		return ""
	}

	return *asset.Description
}

// TestListAssets_Success verifies the rows from the repository are returned
// in order, mapping NULL descriptions to nil.
func TestListAssets_Success(t *testing.T) {
	t.Parallel()

	mock := assetMock()
	mock.listAssetClassesFn = func(_ context.Context) ([]repository.AssetClass, error) {
		return []repository.AssetClass{
			assetRow("11111111-1111-1111-1111-111111111111", "Crypto", "Digital assets"),
			assetRow("22222222-2222-2222-2222-222222222222", "Equities", ""),
		}, nil
	}

	rec := doJSONRequest(t, newAssetRouter(t, mock), http.MethodGet, "/api/assets", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var assets []api.AssetClass

	err := json.NewDecoder(rec.Body).Decode(&assets)
	if err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(assets) != 2 {
		t.Fatalf("expected 2 asset classes, got %d", len(assets))
	}

	if assets[0].Name != "Crypto" {
		t.Errorf("expected first asset named %q, got %q", "Crypto", assets[0].Name)
	}

	if got := assetDescription(assets[0]); got != "Digital assets" {
		t.Errorf("expected description %q, got %q", "Digital assets", got)
	}

	if assets[1].Description != nil {
		t.Errorf("expected NULL description mapped to nil, got %q", *assets[1].Description)
	}
}

// TestCreateAsset_Success verifies a valid request returns 201 with the
// created asset class persisted via the repository.
func TestCreateAsset_Success(t *testing.T) {
	t.Parallel()

	mock := assetMock()
	mock.findAssetClassByNameFn = func(_ context.Context, _ string) (repository.AssetClass, error) {
		return repository.AssetClass{}, sql.ErrNoRows
	}

	var got repository.CreateAssetClassParams

	mock.createAssetClassFn = func(_ context.Context, arg repository.CreateAssetClassParams) error {
		got = arg

		return nil
	}

	rec := createAsset(t, newAssetRouter(t, mock), testAssetName)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	asset := decodeAsset(t, rec)

	if asset.Id == uuid.Nil {
		t.Error("expected a non-nil id")
	}

	if asset.Name != testAssetName {
		t.Errorf("expected name %q, got %q", testAssetName, asset.Name)
	}

	if desc := assetDescription(asset); desc != testAssetDesc {
		t.Errorf("expected description %q, got %q", testAssetDesc, desc)
	}

	if got.Name != testAssetName {
		t.Errorf("expected persisted name %q, got %q", testAssetName, got.Name)
	}

	if got.ID != asset.Id.String() {
		t.Errorf("expected persisted id %s, got %s", asset.Id.String(), got.ID)
	}

	if got.Description != (sql.NullString{String: testAssetDesc, Valid: true}) {
		t.Errorf("expected persisted description %q, got %+v", testAssetDesc, got.Description)
	}
}

// TestCreateAsset_DuplicateName verifies an existing name returns 409.
func TestCreateAsset_DuplicateName(t *testing.T) {
	t.Parallel()

	mock := assetMock()
	mock.findAssetClassByNameFn = func(_ context.Context, _ string) (repository.AssetClass, error) {
		return assetRow(unknownAssetID, testAssetName, ""), nil
	}

	rec := createAsset(t, newAssetRouter(t, mock), testAssetName)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, rec.Code)
	}

	var errResp api.ErrorResponse

	err := json.NewDecoder(rec.Body).Decode(&errResp)
	if err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if errResp.Message != dupNameErrMessage {
		t.Errorf("expected message %q, got %q", dupNameErrMessage, errResp.Message)
	}
}

// TestCreateAsset_RepoError verifies a repository failure surfaces as a 500.
func TestCreateAsset_RepoError(t *testing.T) {
	t.Parallel()

	mock := assetMock()
	mock.findAssetClassByNameFn = func(_ context.Context, _ string) (repository.AssetClass, error) {
		return repository.AssetClass{}, sql.ErrNoRows
	}
	mock.createAssetClassFn = func(_ context.Context, _ repository.CreateAssetClassParams) error {
		return errTestDB
	}

	rec := createAsset(t, newAssetRouter(t, mock), testAssetName)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			rec.Code,
		)
	}
}

// TestCreateAsset_NilBody verifies a nil request body returns 422. The strict
// HTTP wrapper rejects empty bodies with 400, so the handler is exercised
// directly.
func TestCreateAsset_NilBody(t *testing.T) {
	t.Parallel()

	hdl := handler.NewHandler(slog.Default(), assetMock())

	resp, err := hdl.CreateAsset(t.Context(), api.CreateAssetRequestObject{Body: nil})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	validationResp, ok := resp.(api.CreateAsset422JSONResponse)
	if !ok {
		t.Fatalf("expected CreateAsset422JSONResponse, got %T", resp)
	}

	if len(validationResp.Errors) == 0 {
		t.Fatal("expected at least one validation error")
	}

	if validationResp.Errors[0].Field != "body" {
		t.Errorf("expected field 'body', got %s", validationResp.Errors[0].Field)
	}
}

// TestCreateAsset_BlankName verifies an empty name returns 422 before any
// repository call.
func TestCreateAsset_BlankName(t *testing.T) {
	t.Parallel()

	rec := createAsset(t, newAssetRouter(t, assetMock()), "")

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusUnprocessableEntity,
			rec.Code,
		)
	}

	var validation api.ValidationErrors

	err := json.NewDecoder(rec.Body).Decode(&validation)
	if err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(validation.Errors) == 0 {
		t.Fatal("expected at least one validation error")
	}

	if validation.Errors[0].Field != "name" {
		t.Errorf("expected field 'name', got %s", validation.Errors[0].Field)
	}
}

// TestCreateAsset_WhitespaceName verifies a whitespace-only name returns 422.
func TestCreateAsset_WhitespaceName(t *testing.T) {
	t.Parallel()

	rec := createAsset(t, newAssetRouter(t, assetMock()), "   ")

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusUnprocessableEntity,
			rec.Code,
		)
	}
}

// TestGetAsset_Found verifies fetching an existing asset class returns 200
// with the matching payload.
func TestGetAsset_Found(t *testing.T) {
	t.Parallel()

	mock := assetMock()
	mock.findAssetClassByIDFn = func(_ context.Context, id string) (repository.AssetClass, error) {
		return assetRow(id, testAssetName, testAssetDesc), nil
	}

	rec := doJSONRequest(
		t,
		newAssetRouter(t, mock),
		http.MethodGet,
		"/api/assets/"+unknownAssetID,
		"",
	)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	asset := decodeAsset(t, rec)

	if asset.Id.String() != unknownAssetID {
		t.Errorf("expected id %s, got %s", unknownAssetID, asset.Id.String())
	}

	if asset.Name != testAssetName {
		t.Errorf("expected name %q, got %q", testAssetName, asset.Name)
	}

	if desc := assetDescription(asset); desc != testAssetDesc {
		t.Errorf("expected description %q, got %q", testAssetDesc, desc)
	}
}

// TestGetAsset_NotFound verifies an unknown id returns 404.
func TestGetAsset_NotFound(t *testing.T) {
	t.Parallel()

	mock := assetMock()
	mock.findAssetClassByIDFn = func(_ context.Context, _ string) (repository.AssetClass, error) {
		return repository.AssetClass{}, sql.ErrNoRows
	}

	rec := doJSONRequest(
		t,
		newAssetRouter(t, mock),
		http.MethodGet,
		"/api/assets/"+unknownAssetID,
		"",
	)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

// TestUpdateAsset_Success verifies updating an asset class changes the name
// and description while keeping the id.
func TestUpdateAsset_Success(t *testing.T) {
	t.Parallel()

	mock := assetMock()
	mock.findAssetClassByIDFn = func(_ context.Context, id string) (repository.AssetClass, error) {
		return assetRow(id, testAssetName, testAssetDesc), nil
	}
	mock.findAssetClassByNameFn = func(_ context.Context, _ string) (repository.AssetClass, error) {
		return repository.AssetClass{}, sql.ErrNoRows
	}
	mock.updateAssetClassFn = func(
		_ context.Context,
		arg repository.UpdateAssetClassParams,
	) (repository.AssetClass, error) {
		return assetRow(arg.ID, arg.Name, arg.Description.String), nil
	}

	newName := "Private Equity"
	newDesc := "Buyout and venture capital funds"
	rec := doJSONRequest(
		t,
		newAssetRouter(t, mock),
		http.MethodPatch,
		"/api/assets/"+unknownAssetID,
		fmt.Sprintf(`{"name":%q,"description":%q}`, newName, newDesc),
	)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	asset := decodeAsset(t, rec)

	if asset.Id.String() != unknownAssetID {
		t.Errorf("expected id %s, got %s", unknownAssetID, asset.Id.String())
	}

	if asset.Name != newName {
		t.Errorf("expected name %q, got %q", newName, asset.Name)
	}

	if desc := assetDescription(asset); desc != newDesc {
		t.Errorf("expected description %q, got %q", newDesc, desc)
	}
}

// TestUpdateAsset_NotFound verifies updating an unknown id returns 404.
func TestUpdateAsset_NotFound(t *testing.T) {
	t.Parallel()

	mock := assetMock()
	mock.findAssetClassByIDFn = func(_ context.Context, _ string) (repository.AssetClass, error) {
		return repository.AssetClass{}, sql.ErrNoRows
	}

	rec := doJSONRequest(
		t,
		newAssetRouter(t, mock),
		http.MethodPatch,
		"/api/assets/"+unknownAssetID,
		fmt.Sprintf(`{"name":%q,"description":%q}`, testAssetName, testAssetDesc),
	)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

// TestUpdateAsset_DuplicateName verifies renaming to an existing name owned
// by a different class returns 409.
func TestUpdateAsset_DuplicateName(t *testing.T) {
	t.Parallel()

	mock := assetMock()
	mock.findAssetClassByIDFn = func(_ context.Context, id string) (repository.AssetClass, error) {
		return assetRow(id, testAssetName, testAssetDesc), nil
	}
	mock.findAssetClassByNameFn = func(_ context.Context, _ string) (repository.AssetClass, error) {
		return assetRow("33333333-3333-3333-3333-333333333333", testAssetName, ""), nil
	}

	rec := doJSONRequest(
		t,
		newAssetRouter(t, mock),
		http.MethodPatch,
		"/api/assets/"+unknownAssetID,
		fmt.Sprintf(`{"name":%q,"description":%q}`, testAssetName, testAssetDesc),
	)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, rec.Code)
	}

	var errResp api.ErrorResponse

	err := json.NewDecoder(rec.Body).Decode(&errResp)
	if err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if errResp.Message != dupNameErrMessage {
		t.Errorf("expected message %q, got %q", dupNameErrMessage, errResp.Message)
	}
}

// TestUpdateAsset_BlankName verifies updating with an empty name returns 422.
func TestUpdateAsset_BlankName(t *testing.T) {
	t.Parallel()

	rec := doJSONRequest(
		t,
		newAssetRouter(t, assetMock()),
		http.MethodPatch,
		"/api/assets/"+unknownAssetID,
		fmt.Sprintf(`{"name":%q,"description":%q}`, "", testAssetDesc),
	)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusUnprocessableEntity,
			rec.Code,
		)
	}
}

// TestDeleteAsset_Success verifies deleting an asset class returns 204.
func TestDeleteAsset_Success(t *testing.T) {
	t.Parallel()

	mock := assetMock()
	mock.deleteAssetClassFn = func(_ context.Context, id string) (string, error) {
		return id, nil
	}

	rec := doJSONRequest(
		t,
		newAssetRouter(t, mock),
		http.MethodDelete,
		"/api/assets/"+unknownAssetID,
		"",
	)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}
}

// TestDeleteAsset_NotFound verifies deleting an unknown id returns 404.
func TestDeleteAsset_NotFound(t *testing.T) {
	t.Parallel()

	mock := assetMock()
	mock.deleteAssetClassFn = func(_ context.Context, _ string) (string, error) {
		return "", sql.ErrNoRows
	}

	rec := doJSONRequest(
		t,
		newAssetRouter(t, mock),
		http.MethodDelete,
		"/api/assets/"+unknownAssetID,
		"",
	)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

// crudFlowMock returns a mockQuerier wired for a full lifecycle: create
// stores the row, get and update read it back, delete succeeds. The stored
// row is exposed so tests can assert against the persisted state.
func crudFlowMock() (*mockQuerier, *repository.AssetClass) {
	mock := assetMock()

	var stored repository.AssetClass

	mock.findAssetClassByNameFn = func(_ context.Context, _ string) (repository.AssetClass, error) {
		return repository.AssetClass{}, sql.ErrNoRows
	}
	mock.createAssetClassFn = func(_ context.Context, arg repository.CreateAssetClassParams) error {
		stored = repository.AssetClass(arg)

		return nil
	}
	mock.findAssetClassByIDFn = func(_ context.Context, id string) (repository.AssetClass, error) {
		if id != stored.ID {
			return repository.AssetClass{}, sql.ErrNoRows
		}

		return stored, nil
	}
	mock.updateAssetClassFn = func(
		_ context.Context,
		arg repository.UpdateAssetClassParams,
	) (repository.AssetClass, error) {
		stored = repository.AssetClass{
			ID:          arg.ID,
			Name:        arg.Name,
			Description: arg.Description,
		}

		return stored, nil
	}
	mock.deleteAssetClassFn = func(_ context.Context, id string) (string, error) {
		return id, nil
	}

	return mock, &stored
}

// TestAssetCRUDFlow exercises the full lifecycle in one pass: create, get,
// update, then delete.
func TestAssetCRUDFlow(t *testing.T) {
	t.Parallel()

	mock, stored := crudFlowMock()
	router := newAssetRouter(t, mock)
	created := decodeAsset(t, createAsset(t, router, testAssetName))
	path := "/api/assets/" + created.Id.String()

	if stored.ID != created.Id.String() {
		t.Errorf("expected persisted id %s, got %s", created.Id.String(), stored.ID)
	}

	getRec := doJSONRequest(t, router, http.MethodGet, path, "")
	if getRec.Code != http.StatusOK {
		t.Fatalf("expected status %d on get, got %d", http.StatusOK, getRec.Code)
	}

	updateRec := doJSONRequest(
		t,
		router,
		http.MethodPatch,
		path,
		fmt.Sprintf(`{"name":%q,"description":%q}`, testAssetName, "Updated "+testAssetDesc),
	)
	if updateRec.Code != http.StatusOK {
		t.Fatalf("expected status %d on update, got %d", http.StatusOK, updateRec.Code)
	}

	delRec := doJSONRequest(t, router, http.MethodDelete, path, "")
	if delRec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d on delete, got %d", http.StatusNoContent, delRec.Code)
	}
}
