# Ryujin Backend

Go REST API for the Ryujin personal finance platform. Built with **DDD + Ports/Adapters (Hexagonal)** architecture, structured around bounded contexts.

---

## Stack

| Concern        | Technology                          |
| -------------- | ----------------------------------- |
| Language       | Go 1.25                             |
| HTTP Framework | Gin                                 |
| ORM            | GORM + PostgreSQL (pgx driver)      |
| Auth           | JWT (golang-jwt/jwt v5)             |
| Migrations     | golang-migrate (embedded SQL files) |
| Hot Reload     | Air                                 |
| Containerized  | Docker + Docker Compose             |

---

## Project Structure

```
backend/
├── cmd/server/                   # Entry point + wiring
│   ├── main.go                   # Startup, graceful shutdown, background workers
│   ├── dependencies.go           # Dependency injection (repositories → services → controllers)
│   └── router.go                 # Route registration for all modules
│
├── internal/
│   ├── modules/                  # Bounded contexts — each self-contained
│   │   ├── user/                 # Auth, registration, profile, JWT sessions
│   │   ├── finance/              # Income sources, expenses, debts, accounts, categories
│   │   ├── investment/           # Holdings, portfolio summary, stock quotes, price history
│   │   ├── goal/                 # Purchase goals (planned)
│   │   └── dashboard/            # Cross-cutting aggregations (planned)
│   │
│   ├── shared/                   # Shared kernel — no module-specific logic
│   │   ├── domain/               # Value objects (Money, Currency, UserID, DateRange), domain errors
│   │   ├── infrastructure/       # CORS, Logger, RateLimit middlewares, base repository
│   │   └── utils/                # JWT helpers, password hashing, pagination, validator
│   │
│   └── config/                   # Config loading (env vars → struct) and database setup
│
├── migrations/                   # SQL migration files (golang-migrate format)
│   ├── 000001_create_users_table.{up,down}.sql
│   ├── 000002_create_finance_tables.{up,down}.sql
│   └── 000003_create_investment_tables.{up,down}.sql
│
└── pkg/                          # Pure functions — zero DB or HTTP dependencies
    ├── finance/indices.go         # Savings ratio, debt ratio, cash flow, emergency coverage
    └── investment/projections.go  # Future value, compound interest, CAGR
```

---

## Module Structure

Every module follows the same three-layer layout:

```
module/
├── domain/
│   ├── entities/          # Rich domain models with behavior (not anemic)
│   ├── services/          # Domain services for multi-entity business logic
│   ├── repositories/      # Repository interfaces (ports)
│   ├── errors.go          # Domain-specific error types
│   └── value_objects.go   # Immutable domain primitives
│
├── application/
│   ├── services/          # Use case orchestration, transaction boundaries
│   └── dto/               # Request/response data transfer objects
│
└── infrastructure/
    ├── http/
    │   ├── controllers/   # Thin: validate input → call app service → return response
    │   ├── middlewares/   # Module-scoped middlewares (e.g. auth)
    │   └── router.go      # Route registration for this module
    ├── persistence/
    │   ├── models/        # GORM models (separate from domain entities)
    │   ├── mappers/       # Domain entity ↔ GORM model conversion
    │   └── repositories/  # Repository implementations (adapters)
    ├── external/          # Third-party API clients (Yahoo Finance, Alpha Vantage)
    └── worker/            # Background workers (e.g. price refresh)
```

**Dependency rule:** `Infrastructure → Application → Domain`  
Domain has zero external dependencies. Application depends only on domain interfaces.

---

## API Routes

All routes are prefixed with `/api/v1`. Protected routes require `Authorization: Bearer <token>`.

### Health

| Method | Path      | Auth | Description       |
| ------ | --------- | ---- | ----------------- |
| GET    | `/health` | No   | Liveness check    |

### User Module

| Method | Path                    | Auth | Description                |
| ------ | ----------------------- | ---- | -------------------------- |
| POST   | `/auth/register`        | No   | Register a new user        |
| POST   | `/auth/login`           | No   | Login, returns JWT pair    |
| POST   | `/auth/refresh`         | No   | Refresh access token       |
| POST   | `/auth/logout`          | Yes  | Invalidate refresh token   |
| GET    | `/users/me`             | Yes  | Get current user profile   |
| PUT    | `/users/me`             | Yes  | Update profile             |
| PUT    | `/users/me/password`    | Yes  | Change password            |

### Finance Module

| Method | Path                       | Auth | Description                    |
| ------ | -------------------------- | ---- | ------------------------------ |
| GET    | `/categories`              | Yes  | List categories                |
| POST   | `/categories`              | Yes  | Create category                |
| PUT    | `/categories/:id`          | Yes  | Update category                |
| DELETE | `/categories/:id`          | Yes  | Delete category                |
| GET    | `/income-sources`          | Yes  | List income sources            |
| POST   | `/income-sources`          | Yes  | Create income source           |
| PUT    | `/income-sources/:id`      | Yes  | Update income source           |
| DELETE | `/income-sources/:id`      | Yes  | Delete income source           |
| GET    | `/expenses`                | Yes  | List expenses (paginated)      |
| POST   | `/expenses`                | Yes  | Create expense                 |
| PUT    | `/expenses/:id`            | Yes  | Update expense                 |
| DELETE | `/expenses/:id`            | Yes  | Delete expense                 |
| GET    | `/debts`                   | Yes  | List debts                     |
| POST   | `/debts`                   | Yes  | Create debt                    |
| PUT    | `/debts/:id`               | Yes  | Update debt                    |
| DELETE | `/debts/:id`               | Yes  | Delete debt                    |
| GET    | `/accounts`                | Yes  | List accounts                  |
| POST   | `/accounts`                | Yes  | Create account                 |
| PUT    | `/accounts/:id`            | Yes  | Update account                 |
| DELETE | `/accounts/:id`            | Yes  | Delete account                 |

### Investment Module

| Method | Path                           | Auth | Description                          |
| ------ | ------------------------------ | ---- | ------------------------------------ |
| GET    | `/holdings`                    | Yes  | List holdings (paginated)            |
| POST   | `/holdings`                    | Yes  | Create holding (upserts stock quote) |
| GET    | `/holdings/:id`                | Yes  | Get holding by ID                    |
| PUT    | `/holdings/:id`                | Yes  | Update holding                       |
| DELETE | `/holdings/:id`                | Yes  | Delete holding                       |
| POST   | `/holdings/:id/refresh-price`  | Yes  | Force price refresh from market API  |
| GET    | `/portfolio/summary`           | Yes  | Portfolio totals grouped by currency |
| GET    | `/portfolio/performance`       | Yes  | Per-holding unrealized P&L           |
| GET    | `/stocks`                      | Yes  | List cached stock quotes             |
| GET    | `/stocks/:symbol/quote`        | Yes  | Get quote for a symbol               |
| GET    | `/stocks/:symbol/history`      | Yes  | Price history for a symbol           |

---

## API Response Format

All endpoints return a consistent JSON envelope:

```json
{
  "success": true,
  "data": {},
  "message": "optional human-readable message",
  "errors": []
}
```

Paginated list responses include:

```json
{
  "success": true,
  "data": {
    "items": [],
    "total": 100,
    "page": 1,
    "limit": 20
  }
}
```

---

## Environment Variables

Copy `.env.example` to `.env` and fill in your values:

```bash
cp .env.example .env
```

| Variable                 | Default       | Description                              |
| ------------------------ | ------------- | ---------------------------------------- |
| `PORT`                   | `8080`        | HTTP server port                         |
| `GIN_MODE`               | `debug`       | `debug` or `release`                     |
| `LOG_LEVEL`              | `info`        | Log verbosity                            |
| `DB_HOST`                | `localhost`   | PostgreSQL host                          |
| `DB_PORT`                | `5432`        | PostgreSQL port                          |
| `DB_USER`                | —             | Database user                            |
| `DB_PASSWORD`            | —             | Database password                        |
| `DB_NAME`                | `ryujin_db`   | Database name                            |
| `DB_SSLMODE`             | `disable`     | SSL mode for Postgres                    |
| `JWT_SECRET`             | —             | Secret key for signing JWTs              |
| `JWT_ACCESS_DURATION`    | `15m`         | Access token TTL                         |
| `JWT_REFRESH_DURATION`   | `24h`         | Refresh token TTL                        |
| `CORS_ALLOWED_ORIGINS`   | `localhost:*` | Comma-separated allowed origins          |
| `RATE_LIMIT_ENABLED`     | `true`        | Enable/disable rate limiting             |
| `RATE_LIMIT_RPS`         | `10`          | Requests per second per IP               |
| `RATE_LIMIT_BURST`       | `20`          | Burst size                               |

---

## Running Locally

### With Docker (recommended)

From the project root:

```bash
docker compose up --build
```

This starts PostgreSQL, runs migrations automatically, and boots the API on `http://localhost:8080`.

### Without Docker

Requirements: Go 1.21+, PostgreSQL 15+

```bash
# Install Air for hot reload
go install github.com/air-verse/air@latest

# Set up environment
cp .env.example .env
# edit .env with your DB credentials

# Run with hot reload
air

# Or run directly
go run ./cmd/server
```

---

## Database Migrations

Migrations run automatically on startup via embedded SQL files (golang-migrate).  
Migration files live in `migrations/` and follow the `{version}_{description}.{up,down}.sql` convention.

To add a new migration:

```bash
# Naming convention: sequential number + snake_case description
# Example:
touch migrations/000004_create_goals_table.up.sql
touch migrations/000004_create_goals_table.down.sql
```

> The `migrations.go` file uses `//go:embed *.sql` so new files are picked up automatically on next build.

---

## Key Conventions

- **Money is never a float.** All monetary values are stored as `BIGINT` cents. Use the `Money` value object — never `float64`.
- **All queries are scoped by `user_id`.** This is a multi-tenant system.
- **GORM models are separate from domain entities.** Mappers convert between them — domain stays pure.
- **Controllers are thin.** Validate input → call application service → return response. No business logic.
- **`pkg/` contains only pure functions.** No DB, no HTTP, no side effects — easy to test in isolation.
- **Monetary amounts in the DB:** `BIGINT` (cents). Ratios and percentages: `NUMERIC(10,4)`.
- **UUIDs** for all primary keys. Generated at the DB level (`uuid_generate_v4()`).
- **Soft deletes** on users (`deleted_at`). Hard deletes on financial records.
