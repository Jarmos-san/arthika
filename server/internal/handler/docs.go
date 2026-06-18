package handler

import (
	"log/slog"
	"net/http"

	"github.com/Jarmos-san/arthika/server/internal/api"
)

// DocsHandler serves the OpenAPI specification and a Redoc documentation UI.
type DocsHandler struct {
	logger *slog.Logger
}

// NewDocsHandler creates a new DocsHandler.
func NewDocsHandler(logger *slog.Logger) *DocsHandler {
	return &DocsHandler{logger: logger}
}

// GetSpecJSON responds with the raw OpenAPI specification as JSON.
func (h *DocsHandler) GetSpecJSON(writer http.ResponseWriter, _ *http.Request) {
	spec, err := api.GetSpecJSON()
	if err != nil {
		h.logger.Error("failed to get spec JSON", slog.Any("error", err))
		http.Error(
			writer,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)

		return
	}

	writer.Header().Set("Content-Type", "application/vnd.oai.openapi+json")
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(spec)
}

// docsHTML is the Redoc standalone HTML page.
// It fetches the spec from /openapi.json and renders it client-side.
const docsHTML = `<!doctype html>
<html>
<head>
	<title>Arthika API Docs</title>
	<meta charset="utf-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
</head>
<body>
	<div id="redoc-container"></div>
	<script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"></script>
	<script>
		Redoc.init('/openapi.json', {}, document.getElementById('redoc-container'));
	</script>
</body>
</html>`

// GetDocsPage renders the Redoc HTML documentation UI.
func (h *DocsHandler) GetDocsPage(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte(docsHTML))
}
