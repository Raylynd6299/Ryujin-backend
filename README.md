# Ryujin Backend

This folder contains the Go API following DDD + Ports/Adapters.

## Structure

- `cmd/` Entry points
- `internal/modules/` Bounded contexts (user, finance, investment, goal, dashboard)
- `internal/shared/` Shared kernel (domain, infrastructure, utils)
- `internal/config/` Configuration and setup
- `migrations/` Database migrations
- `pkg/` Pure functions (finance, investment)
