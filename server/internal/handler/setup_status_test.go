package handler_test

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jarmos-san/arthika/server/internal/api"
	"github.com/Jarmos-san/arthika/server/internal/handler"
)

// TestSystemStatus_NeedsSetup verifies that an empty database returns
// needsSetup=true.
func TestSystemStatus_NeedsSetup(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn:      nil,
		findUserByEmailFn: nil,
		countUsersFn: func(_ context.Context) (int64, error) {
			return 0, nil
		},
		createAssetClassFn:     nil,
		deleteAssetClassFn:     nil,
		findAssetClassByIDFn:   nil,
		findAssetClassByNameFn: nil,
		listAssetClassesFn:     nil,
		updateAssetClassFn:     nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock)

	resp, err := hdl.SystemStatus(t.Context(), api.SystemStatusRequestObject{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	statusResp, ok := resp.(api.SystemStatus200JSONResponse)
	if !ok {
		t.Fatalf("expected SystemStatus200JSONResponse, got %T", resp)
	}

	if !statusResp.NeedsSetup {
		t.Error("expected needsSetup=true, got false")
	}
}

// TestSystemStatus_NoSetupNeeded verifies that a database with existing users
// returns needsSetup=false.
func TestSystemStatus_NoSetupNeeded(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn:      nil,
		findUserByEmailFn: nil,
		countUsersFn: func(_ context.Context) (int64, error) {
			return 1, nil
		},
		createAssetClassFn:     nil,
		deleteAssetClassFn:     nil,
		findAssetClassByIDFn:   nil,
		findAssetClassByNameFn: nil,
		listAssetClassesFn:     nil,
		updateAssetClassFn:     nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock)

	resp, err := hdl.SystemStatus(t.Context(), api.SystemStatusRequestObject{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	statusResp, ok := resp.(api.SystemStatus200JSONResponse)
	if !ok {
		t.Fatalf("expected SystemStatus200JSONResponse, got %T", resp)
	}

	if statusResp.NeedsSetup {
		t.Error("expected needsSetup=false, got true")
	}
}

// TestSystemStatus_HTTPEndpoint verifies the full HTTP stack returns 200 with
// the correct JSON body.
func TestSystemStatus_HTTPEndpoint(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn:      nil,
		findUserByEmailFn: nil,
		countUsersFn: func(_ context.Context) (int64, error) {
			return 0, nil
		},
		createAssetClassFn:     nil,
		deleteAssetClassFn:     nil,
		findAssetClassByIDFn:   nil,
		findAssetClassByNameFn: nil,
		listAssetClassesFn:     nil,
		updateAssetClassFn:     nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock)
	strictHandler := api.NewStrictHandler(hdl, nil)

	req := httptest.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		"/api/setup/status",
		nil,
	)

	rec := httptest.NewRecorder()
	strictHandler.SystemStatus(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	expected := `{"needsSetup":true}` + "\n"
	if rec.Body.String() != expected {
		t.Errorf("expected body %q, got %q", expected, rec.Body.String())
	}
}
