# Arthika — AGENTS.md

## Quick start

```bash
go test ./server/...        # all tests (t.Parallel(), black-box)
task lint                   # golangci-lint v2, all linters
task migrate:up             # apply pending DB migrations
go run ./server/cmd/server  # start server on :8000
task generate               # re-generate sqlc code

cd client && pnpm dev        # Nuxt 4 dev server
cd client && pnpm typecheck  # vue-tsc --noEmit

pre-commit install           # install git hooks (includes crisp for commit messages)
```

## Architecture

- Monorepo: `server/` (Go) + `client/` (Nuxt 4 / Vue 3 / TypeScript)
- Entrypoint: `server/cmd/server/main.go`
- Layered: **Handler → Repository** — no service package yet (deliberate MVP choice)
- Packages: `api/`, `app/`, `auth/`, `config/`, `handler/`, `logger/`, `middleware/`, `repository/` under `server/internal/`
- API spec in `server/api/openapi.yml`, JSON:API format (`application/vnd.api+json`)
- Routing: `go-chi/chi/v5`
- Logging: `slog` JSON handler to stdout
- Config from env vars: `ADDR`, `READ_TIMEOUT`, `WRITE_TIMEOUT`, `IDLE_TIMEOUT`, `LOG_LEVEL`, `TOKEN_SECRET`, `DATABASE_URL`, `MIGRATIONS_DIR`
- Database: SQLite via `golang-migrate` (auto-applied at startup)

## Code generation

```bash
task generate                       # sqlc generate (reads sqlc.yaml, writes internal/repository/)
oapi-codegen -config oapi-codegen.yml api/openapi.yml   # writes internal/api/server.gen.go
```

Run codegen after changing `db/query/*.sql` (sqlc) or `api/openapi.yml` (oapi-codegen), then commit generated files.

## Client

- Nuxt 4 with `ssr: false` (SPA mode), `@nuxt/ui`, `vue-router`
- Package manager: `pnpm` (lockfile `pnpm-lock.yaml`)
- Commands: `pnpm dev` / `pnpm build` / `pnpm typecheck` / `pnpm generate`

## CI

Three workflows on PR to `main`:
- **Server QA** — lint (`golangci-lint`), test (`go test ./...`), build check
- **Client QA** — typecheck (`vue-tsc`), build (`nuxt build`)
- **SQL Validation** — `sqlc vet` + `sqlc diff` (ensures generated code is in sync)

## Testing

- All tests use `t.Parallel()` and black-box packages (`_test` suffix)
- Handler tests use `mockQuerier` implementing `repository.Querier`
- Pre-compute bcrypt hashes (`bcrypt.MinCost`) as package-level constants to avoid slow hashing per test
- Use `t.Context()` and `slog.DiscardHandler` / `slog.Default()`
- Integration-style handler tests use `httptest.NewRecorder` + `api.NewStrictHandler`

## Conventions

- `depguard` restricts `internal/*.go` to stdlib imports only
- Go: tab indentation; JSON/YAML/TS/Vue: 2-space (see `.editorconfig`)
- Formatters: `gci`, `gofmt`, `gofumpt`, `goimports`, `golines` (via golangci-lint)
- Migrations: `golang-migrate` CLI, naming `NNNNNN_description.{up,down}.sql`
- Commit messages linted by `crisp` pre-commit hook
- Pre-commit hooks: `trailing-whitespace`, `end-of-file-fixer`, `check-yaml`, `check-added-large-files`, `crisp`
