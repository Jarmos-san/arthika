# Arthika — AGENTS.md

## Quick start

```bash
go test ./...          # all tests (all use t.Parallel())
task lint              # golangci-lint v2 with all linters enabled
go run ./cmd/server    # start the server (default :8000)
```

## Architecture

- Entrypoint: `server/cmd/server/main.go` (not `cmd/api/main.go` as ARCHITECTURE.md states)
- Layered: **Handler → Service → (Repository — not yet implemented)**
- API format: JSON:API (`application/vnd.api+json`)
- Config from env vars only: `ADDR`, `READ_TIMEOUT`, `WRITE_TIMEOUT`, `IDLE_TIMEOUT`, `LOG_LEVEL`, `TOKEN_SECRET`
- Logging: `slog` JSON handler to stdout
- No database layer yet; no frontend yet (ARCHITECTURE.md is partially aspirational)

## Conventions

- All tests use `t.Parallel()` and black-box test packages (`_test` suffix)
- `depguard` restricts imports in `internal/*.go` to stdlib only
- Commit messages linted by `crisp` pre-commit hook
- Go: tab indentation; JSON/YAML: 2-space (see `.editorconfig`)
- Formatters: `gci`, `gofmt`, `gofumpt`, `goimports`, `golines` (via golangci-lint)
