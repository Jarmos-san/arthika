# ARCHITECTURE.md

## Overview

This project is a privacy-first, self-hostable web application for managing
financial investments, tracking portfolios, and maintaining research on
potential investments. The system is designed to minimize data exposure while
providing users with full control over their data, including optional anonymous
public sharing.

---

## Core Principles

### 1. Privacy First

- No unnecessary personally identifiable information (PII)
- No third-party tracking or analytics
- All data remains under user control
- Public sharing is strictly opt-in and anonymized

### 2. Self-Hostable by Design

- Minimal external dependencies
- Simple deployment via Docker
- Works in low-resource environments

### 3. Simplicity Over Complexity

- Prefer clear and maintainable abstractions
- Avoid premature optimization and over-engineering
- Incrementally introduce complexity only when justified

### 4. Modular Architecture

- Separation of concerns across layers
- Clear domain boundaries

---

## High-Level Architecture

### Tech Stack

- **Backend:** Go + Chi router
- **Frontend:** Nuxt.js (Vue 3) — *not yet implemented*
- **Database:** SQLite

### System Layout

```
Frontend (Nuxt.js)   — not yet implemented
        ↓
Backend API (Go + Chi)
        ↓
Database (SQLite)
```

---

## Backend Architecture (Go)

### Layered Design

```
Handler → Service → Repository → Database

> **Status:** All layers are implemented. Repository is sqlc-generated from hand-written SQL. Database is SQLite with auto-applied migrations via golang-migrate.
```

#### Handler Layer

- Responsible for HTTP request/response handling
- Performs input validation and serialization
- Delegates business logic to services

#### Service Layer

- Contains core business logic
- Implements domain rules and computations
- Orchestrates multiple repositories if needed

#### Repository Layer

- Handles data persistence
- Abstracts database operations
- Keeps storage concerns isolated

#### Domain Layer

- Defines core entities and types
- Contains minimal logic tied to business concepts

---

## Project Structure

```
cmd/
  server/
    main.go

db/
  migrations/   # SQL schema migrations (golang-migrate)
  query/        # hand-written SQL queries (sqlc input)

internal/
  app/          # composition root, lifecycle, DB + migrations
  handler/      # HTTP transport layer
  service/      # business logic
  repository/   # data access (sqlc-generated)
  dto/          # JSON:API request/response types
  middleware/   # HTTP middleware — placeholder (chi built-in Logger+Recoverer used globally)
  config/       # environment-based configuration
  logger/       # structured slog setup
```

### Configuration

Configuration is loaded from environment variables only. All fields have sensible defaults suitable for local development.

| Env Var          | Config Field        | Default             |
| ---------------- | ------------------- | ------------------- |
| `ADDR`             | `Addr`                | `:8000`               |
| `READ_TIMEOUT`     | `ReadTimeout`         | `10s`                 |
| `WRITE_TIMEOUT`    | `WriteTimeout`        | `10s`                 |
| `IDLE_TIMEOUT`     | `IdleTimeout`         | `10s`                 |
| `LOG_LEVEL`        | `LogLevel`            | `info`                |
| `TOKEN_SECRET`     | `TokenSecret`         | `super-secret-token`  |
| `DATABASE_URL`     | `DatabaseURL`         | `db/arthika.db`       |
| `MIGRATIONS_DIR`   | `MigrationDirectory`  | `db/migrations`       |

Invalid duration values fall back to defaults without failing startup.

---

## Core Domains

### 1. Portfolio Management

Handles tracking of user investments.

#### Entities

- Instrument (Equity, Commodity, Debt)
- Transaction (Buy/Sell)
- Holding (derived, not stored initially)

#### Responsibilities

- Aggregate transactions
- Compute portfolio metrics:
  - Total investment
  - Average cost
  - Current value
  - Profit & Loss (PnL)

---

### 2. Watchlist & Research

Tracks companies and investment ideas.

#### Features

- Add/remove companies
- Store fundamental metrics (PE, ROE, etc.)
- Maintain notes and investment thesis

---

### 3. Privacy & Sharing

Enables controlled public visibility.

#### Snapshot System

```
Portfolio → Snapshot → Public URL
```

#### Characteristics

- Generated using random identifiers
- No linkage to user identity
- Read-only access

---

### 4. Authentication

#### Requirements

- Minimal data collection
- Local authentication only

#### Features

- Username + password authentication
- Password hashing (bcrypt)
- Token-based authentication (JWT or PASETO)

---

## API Design

### Principles

- RESTful endpoints
- JSON request/response format
- Stateless authentication

### Implemented Endpoints

```
GET    /users/              # placeholder user retrieval
POST   /users/register      # user registration (persists to SQLite)
POST   /login               # JWT login
```

### Planned Endpoints

```
GET    /portfolio
POST   /transactions
GET    /transactions
GET    /watchlist
POST   /watchlist
POST   /snapshots
GET    /snapshots/{id}
```

---

## Frontend Architecture (Nuxt.js) — *planned*

### Structure

```
pages/
components/
composables/
stores/
```

### Responsibilities

- UI rendering
- State management
- API communication

### Key Features

- Portfolio dashboard
- Transaction management UI
- Watchlist and research pages
- Public portfolio view (read-only)

---

## Data Storage

### Database

#### SQLite

- Ideal for single-user/self-hosted setups
- Lightweight and easy to deploy

#### Schema Migrations

- Managed via [golang-migrate](https://github.com/golang-migrate/migrate)
- Applied automatically on application startup
- Location: `db/migrations/`

#### Code Generation

- SQL queries are hand-written in `db/query/`
- Go code is generated into `internal/repository/` by [sqlc](https://sqlc.dev)
- The generated `Querier` interface enables mock-based testing

#### Current Schema

- `users`: `id TEXT PK`, `username TEXT UNIQUE`, `email TEXT UNIQUE`, `password_hash TEXT`

---

## Privacy Architecture

### Data Isolation

- All data scoped per user
- No cross-user data exposure

### Public Sharing

- Explicit snapshot generation
- No implicit data exposure

### No External Tracking

- No analytics services
- No telemetry pipelines

---

## Deployment

### Requirements

- Docker support
- Minimal configuration

### Components

- Go backend (single binary)
- Database (PostgreSQL or SQLite)
- Reverse proxy (optional, e.g., Caddy)

---

## Observability

### Logging

- Structured logging (JSON or key-value)
- No sensitive data in logs

### Metrics (Optional)

- Local-only metrics collection

---

## Testing Strategy

### Backend

- Unit tests for services
- Integration tests for API endpoints
- Repository layer uses sqlc-generated `Querier` interface, enabling mock injection in handler and service tests

### Frontend

- Component-level testing
- Optional end-to-end testing

---

## Future Enhancements

- Multi-currency support
- Broker statement import (CSV)
- Tax calculation engine
- Offline-first capabilities
- Plugin system for external data providers

---

## Non-Goals (Initial Version)

- Real-time market data integration
- High-frequency trading features
- Social or community features

---

## Summary

This architecture prioritizes privacy, simplicity, and maintainability. The
system is designed to evolve incrementally while maintaining strong separation
of concerns and a clean developer experience.
