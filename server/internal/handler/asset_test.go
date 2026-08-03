package handler_test

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Jarmos-san/arthika/server/internal/api"
	"github.com/Jarmos-san/arthika/server/internal/handler"
	chi "github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const (
	seedAssetCount = 3
	testAssetName  = "Real Estate"
	testAssetDesc  = "Property investments"
	unknownAssetID = "00000000-0000-0000-0000-000000000000"
)

// newAssetRouter mounts a fresh Handler on a chi router the same way main.go
// does. The querier is nil because the asset endpoints never touch the
// database.
func newAssetRouter(t *testing.T) http.Handler {
	t.Helper()

	hdl := handler.NewHandler(slog.Default(), nil)
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

// TestListAssets_Seeded verifies the store starts with the seeded classes in
// name order.
func TestListAssets_Seeded(t *testing.T) {
	t.Parallel()

	router := newAssetRouter(t)
	rec := doJSONRequest(t, router, http.MethodGet, "/api/assets", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var assets []api.AssetClass

	err := json.NewDecoder(rec.Body).Decode(&assets)
	if err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(assets) != seedAssetCount {
		t.Fatalf("expected %d seeded asset classes, got %d", seedAssetCount, len(assets))
	}

	expectedNames := []string{"Crypto", "Equities", "Fixed income"}
	for i, want := range expectedNames {
		if assets[i].Name != want {
			t.Errorf("expected asset %d named %q, got %q", i, want, assets[i].Name)
		}
	}
}

// TestCreateAsset_Success verifies a valid request returns 201 with the
// created asset class, and the list grows by one.
func TestCreateAsset_Success(t *testing.T) {
	t.Parallel()

	router := newAssetRouter(t)
	rec := createAsset(t, router, testAssetName)

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

	if asset.Description != testAssetDesc {
		t.Errorf("expected description %q, got %q", testAssetDesc, asset.Description)
	}

	listRec := doJSONRequest(t, router, http.MethodGet, "/api/assets", "")

	var assets []api.AssetClass

	err := json.NewDecoder(listRec.Body).Decode(&assets)
	if err != nil {
		t.Fatalf("decode list response: %v", err)
	}

	if len(assets) != seedAssetCount+1 {
		t.Errorf("expected %d asset classes after create, got %d", seedAssetCount+1, len(assets))
	}
}

// TestCreateAsset_HTTPEndpoint_EmptyBody verifies the HTTP stack rejects an
// empty request body with 400 before the request reaches the handler.
func TestCreateAsset_HTTPEndpoint_EmptyBody(t *testing.T) {
	t.Parallel()

	router := newAssetRouter(t)
	rec := doJSONRequest(t, router, http.MethodPost, "/api/assets", "")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

// TestCreateAsset_NilBody verifies a nil request body returns 422. The strict
// HTTP wrapper rejects empty bodies with 400, so the handler is exercised
// directly.
func TestCreateAsset_NilBody(t *testing.T) {
	t.Parallel()

	hdl := handler.NewHandler(slog.Default(), nil)

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

// TestCreateAsset_BlankName verifies an empty name returns 422.
func TestCreateAsset_BlankName(t *testing.T) {
	t.Parallel()

	router := newAssetRouter(t)
	rec := createAsset(t, router, "")

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

	router := newAssetRouter(t)
	rec := createAsset(t, router, "   ")

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusUnprocessableEntity,
			rec.Code,
		)
	}
}

// TestGetAsset_Found verifies fetching a created asset class returns 200 with
// the matching payload.
func TestGetAsset_Found(t *testing.T) {
	t.Parallel()

	router := newAssetRouter(t)
	created := decodeAsset(t, createAsset(t, router, testAssetName))

	rec := doJSONRequest(
		t,
		router,
		http.MethodGet,
		"/api/assets/"+created.Id.String(),
		"",
	)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	asset := decodeAsset(t, rec)

	if asset.Id != created.Id {
		t.Errorf("expected id %s, got %s", created.Id, asset.Id)
	}

	if asset.Name != testAssetName {
		t.Errorf("expected name %q, got %q", testAssetName, asset.Name)
	}
}

// TestGetAsset_NotFound verifies an unknown id returns 404.
func TestGetAsset_NotFound(t *testing.T) {
	t.Parallel()

	router := newAssetRouter(t)
	rec := doJSONRequest(
		t,
		router,
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

	router := newAssetRouter(t)
	created := decodeAsset(t, createAsset(t, router, testAssetName))

	newName := "Private Equity"
	newDesc := "Buyout and venture capital funds"
	rec := doJSONRequest(
		t,
		router,
		http.MethodPatch,
		"/api/assets/"+created.Id.String(),
		fmt.Sprintf(`{"name":%q,"description":%q}`, newName, newDesc),
	)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	asset := decodeAsset(t, rec)

	if asset.Id != created.Id {
		t.Errorf("expected id %s, got %s", created.Id, asset.Id)
	}

	if asset.Name != newName {
		t.Errorf("expected name %q, got %q", newName, asset.Name)
	}

	if asset.Description != newDesc {
		t.Errorf("expected description %q, got %q", newDesc, asset.Description)
	}
}

// TestUpdateAsset_NotFound verifies updating an unknown id returns 404.
func TestUpdateAsset_NotFound(t *testing.T) {
	t.Parallel()

	router := newAssetRouter(t)
	rec := doJSONRequest(
		t,
		router,
		http.MethodPatch,
		"/api/assets/"+unknownAssetID,
		fmt.Sprintf(`{"name":%q,"description":%q}`, testAssetName, testAssetDesc),
	)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

// TestUpdateAsset_BlankName verifies updating with an empty name returns 422.
func TestUpdateAsset_BlankName(t *testing.T) {
	t.Parallel()

	router := newAssetRouter(t)
	created := decodeAsset(t, createAsset(t, router, testAssetName))

	rec := doJSONRequest(
		t,
		router,
		http.MethodPatch,
		"/api/assets/"+created.Id.String(),
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

// TestDeleteAsset_Success verifies deleting an asset class returns 204, the
// asset is gone afterwards, and deleting it again returns 404.
func TestDeleteAsset_Success(t *testing.T) {
	t.Parallel()

	router := newAssetRouter(t)
	created := decodeAsset(t, createAsset(t, router, testAssetName))
	path := "/api/assets/" + created.Id.String()

	delRec := doJSONRequest(t, router, http.MethodDelete, path, "")
	if delRec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, delRec.Code)
	}

	getRec := doJSONRequest(t, router, http.MethodGet, path, "")
	if getRec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d after delete, got %d", http.StatusNotFound, getRec.Code)
	}

	secondDelRec := doJSONRequest(t, router, http.MethodDelete, path, "")
	if secondDelRec.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d for second delete, got %d",
			http.StatusNotFound,
			secondDelRec.Code,
		)
	}
}

// TestAssetCRUDFlow exercises the full lifecycle in one pass: create, get,
// update, then delete.
func TestAssetCRUDFlow(t *testing.T) {
	t.Parallel()

	router := newAssetRouter(t)
	created := decodeAsset(t, createAsset(t, router, testAssetName))
	path := "/api/assets/" + created.Id.String()

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
