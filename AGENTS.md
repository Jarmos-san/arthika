# Arthika — AGENTS.md

## Quick start

```bash
go test ./...          # all tests (all use t.Parallel())
task lint              # golangci-lint v2 with all linters enabled
task migrate:up        # apply pending database migrations
go run ./cmd/server    # start the server (default :8000)
```

## Architecture

- Entrypoint: `server/cmd/server/main.go`
- Layered: **Handler → Service → (Repository — not yet implemented)**
- Packages: `app/`, `handler/`, `service/`, `domain/`, `dto/`, `config/`, `logger/` under `server/internal/`
- API format: JSON:API (`application/vnd.api+json`)
- Config from env vars only: `ADDR`, `READ_TIMEOUT`, `WRITE_TIMEOUT`, `IDLE_TIMEOUT`, `LOG_LEVEL`, `TOKEN_SECRET`
- Routing: `go-chi/chi/v5`
- Logging: `slog` JSON handler to stdout
- Database: SQLite via `golang-migrate` CLI (`task migrate:up`/`:down`/`:reset`)
- No frontend yet

## Conventions

- All tests use `t.Parallel()` and black-box test packages (`_test` suffix)
- `depguard` restricts imports in `internal/*.go` to stdlib only
- Commit messages linted by `crisp` pre-commit hook
- Go: tab indentation; JSON/YAML: 2-space (see `.editorconfig`)
- Formatters: `gci`, `gofmt`, `gofumpt`, `goimports`, `golines` (via golangci-lint)
- Migrations: `golang-migrate` CLI, naming convention `NNNNNN_description.{up,down}.sql`
