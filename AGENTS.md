# Arthika — AGENTS.md

## Quick start

```bash
task setup                      # install all deps (server + client)
task server:test                # go test ./...
task server:lint                # golangci-lint v2
task migrate:up                 # apply pending DB migrations
task oapi:gen                   # regenerate Go stubs + TS types from OpenAPI spec
go run ./server/cmd/server      # start server on :8000
cd client && pnpm dev           # Nuxt dev server
```

Run a single test:

```bash
go test ./server/internal/handler/ -run TestRegister_Success
```

## Architecture

- Monorepo: `server/` (Go) + `client/` (Nuxt 4 / Vue 3 / TypeScript)
- Entrypoint: `server/cmd/server/main.go`
- Layers: **Handler → Repository** (no service package — deliberate MVP choice)
- Server packages under `server/internal/`: `api/`, `app/`, `auth/`, `config/`, `handler/`, `logger/`, `middleware/`, `repository/`
- Client app layout: `client/app/` with `pages/`, `components/`, `stores/`, `utils/`, `middleware/`, `layouts/`, `assets/`, `openapi/`, `types/`
- API spec: `openapi.yml` (project root) — JSON:API format (`application/vnd.api+json`)
- Routing: `go-chi/chi/v5` with strict server interface (`oapi-codegen`)
- Logging: `slog` JSON handler to stdout
- Database: SQLite via `golang-migrate` (auto-applied at startup)
- Config from env vars: `ADDR`, `READ_TIMEOUT`, `WRITE_TIMEOUT`, `IDLE_TIMEOUT`, `LOG_LEVEL`, `TOKEN_SECRET`, `DATABASE_URL`, `MIGRATIONS_DIR`
- Docs served at `/docs` (Redoc UI) and `/openapi.json`

## Code generation

Three generators, run in this order when the OpenAPI spec changes:

```bash
cd server && sqlc generate                            # db/query/*.sql → internal/repository/
cd server && oapi-codegen -config oapi-codegen.yml ../openapi.yml  # → internal/api/server.gen.go
cd client && pnpm run openapi:generate                # kubb → client/app/openapi/
```

Shortcut: `task oapi:gen` runs all three in order.

**`task server:generate` is broken** — it references a non-existent `oapi-codegen` subtask. Use `task oapi:gen` instead.

Run codegen after changing `db/query/*.sql` (sqlc) or `openapi.yml` (oapi-codegen + kubb), then commit generated files.

## Client

- Nuxt 4 (`compatibilityVersion: 4`), SSR enabled (`ssr: true`)
- UI: `reka-ui` (headless components via `reka-ui/nuxt` module)
- Styling: Tailwind CSS v4 via `@tailwindcss/vite`
- Package manager: `pnpm` (lockfile `pnpm-lock.yaml`)
- Linting: `oxlint` (`pnpm lint`) — config in `client/oxlint.config.ts`
- Formatting: `oxfmt` (`pnpm fmt` / `pnpm fmt:check`)
- TypeScript API types: kubb (`pnpm run openapi:generate` → `client/app/openapi/`)
- Typecheck requires `.nuxt/` — run `pnpm install` or `task client:install` first
- `postinstall` hook runs `nuxt prepare` which generates `.nuxt/`
- Generated Kubb output (`app/openapi/`) is excluded from oxlint filename-case rules via override

## Client types

- `client/app/openapi/` — Kubb-generated types from `openapi.yml` (do not edit manually)
- `client/app/types/` — hand-written types, mirrors `app/` directory structure:
  - `types/components/` — component prop types (e.g. `Input.ts`), PascalCase filenames
  - `types/utils/` — utility types (e.g. `validators.ts`)
- Types use default exports; imports use `~/types/...` paths

## Testing

### Server

- All tests use `t.Parallel()` and black-box packages (`_test` suffix)
- Handler tests share a `mockQuerier` struct implementing `repository.Querier` — set only the functions each test needs, nil functions panic if called
- Pre-compute bcrypt hashes (`bcrypt.MinCost`) as package-level constants (see `login_test.go` `testPasswordHash`) to avoid slow hashing per test
- Use `t.Context()` and `slog.Default()` (or `slog.DiscardHandler`)
- Integration-style handler tests use `httptest.NewRecorder` + `api.NewStrictHandler`
- Handler implements `api.StrictServerInterface` (compile-time check in `handler.go`)

### Client

- Two vitest projects: `unit` (node env, `tests/unittests/` + `tests/store/`) and `nuxt` (nuxt env, `tests/nuxt/`)
- Run all client tests: `pnpm test` or `task client:test`
- Run a single test file: `pnpm vitest run tests/unittests/validators.test.ts`

## CI

Three workflows on PR to `main`:
- **Server QA** — lint (`golangci-lint`), test (`go test ./...`), build check (`go build -o /dev/null ./...`)
- **Client QA** — lint (`oxlint`), format check (`oxfmt --check`), typecheck (`vue-tsc`), test (`vitest`), build (`nuxt build`)
- **SQL Validation** — `sqlc vet` + `sqlc diff` (ensures generated Go code matches SQL)

## Conventions

- `depguard` restricts `internal/*.go` to stdlib imports only
- `wsl_v5` is disabled in `.golangci.yml` (was previously enabled)
- Go: tab indentation; JSON/YAML/TS/Vue/CSS: 2-space (see `.editorconfig`)
- Formatters: `gci`, `gofmt`, `gofumpt`, `goimports`, `golines` (via golangci-lint)
- Migrations: `golang-migrate` CLI, naming `NNNNNN_description.{up,down}.sql`
- Commit messages linted by `crisp` pre-commit hook
- Pre-commit hooks: `trailing-whitespace`, `end-of-file-fixer`, `check-yaml`, `check-added-large-files`, `oxfmt-check` (client), `oxlint` (client), `crisp`
