package handler_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Jarmos-san/arthika/server/internal/handler"
)

func TestGetSpecJSON_Success(t *testing.T) {
	t.Parallel()

	h := handler.NewDocsHandler(slog.Default())
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/openapi.json", nil)
	rec := httptest.NewRecorder()

	h.GetSpecJSON(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	ct := rec.Header().Get("Content-Type")
	if ct != "application/vnd.oai.openapi+json" {
		t.Fatalf("unexpected content-type: %s", ct)
	}

	if rec.Body.Len() == 0 {
		t.Fatal("expected non-empty spec body")
	}
}

func TestGetDocsPage_Success(t *testing.T) {
	t.Parallel()

	h := handler.NewDocsHandler(slog.Default())
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/docs", nil)
	rec := httptest.NewRecorder()

	h.GetDocsPage(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	ct := rec.Header().Get("Content-Type")
	if ct != "text/html; charset=utf-8" {
		t.Fatalf("unexpected content-type: %s", ct)
	}

	if !strings.Contains(rec.Body.String(), "Redoc") {
		t.Fatal("expected HTML to contain Redoc reference")
	}
}
