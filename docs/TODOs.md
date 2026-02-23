# Backend TODOs - Ryujin

## 🔴 Crítico

### User Module — Application Layer
- [ ] Implementar DTOs: `RegisterRequest`, `LoginRequest`, `UpdateProfileRequest`, `ChangePasswordRequest`, `UserResponse`
- [ ] Implementar `AuthService` (register, login, refresh token, logout)
- [ ] Implementar `ProfileService` (get profile, update profile, change password, update currencies)

### User Module — Infrastructure / Persistence
- [ ] Implementar GORM model `UserModel` en `persistence/models/`
- [ ] Implementar mapper `UserMapper` (domain entity ↔ GORM model)
- [ ] Implementar `UserRepositoryGorm` (implementación del port `UserRepository`)

### User Module — Infrastructure / HTTP
- [ ] Implementar `AuthController` (POST /register, POST /login, POST /refresh, POST /logout)
- [ ] Implementar `ProfileController` (GET /me, PUT /me, PATCH /me/password, PATCH /me/currencies)
- [ ] Implementar `AuthMiddleware` (validar JWT en requests protegidos)
- [ ] Implementar `router.go` del módulo user y registrarlo en el router principal

### Wiring
- [ ] Registrar el User Module en `cmd/server/dependencies.go`
- [ ] Descomentar y conectar rutas en `cmd/server/router.go`

## 🟡 Importante

### Seguridad
- [ ] Migrar a refresh tokens con persistencia en DB (actualmente stateless — no permite revocación)
- [ ] Implementar blacklist de tokens revocados (logout real)
- [ ] Agregar rate limiting específico para endpoints de auth

### Otros Módulos (pendientes después de User)
- [ ] Finance Module: entities, repositories, application services, HTTP layer
- [ ] Investment Module: entities, repositories, application services, HTTP layer, external API clients
- [ ] Goal Module: entities, repositories, application services, HTTP layer

### Dashboard Module — Endpoint `GET /api/v1/dashboard/summary`
> ⚠️ Pendiente hasta tener Investment Module y Goal Module terminados.
> El frontend usa datos calculados del Finance Module en el interim.

- [ ] `internal/modules/dashboard/application/dto/dashboard_dto.go`
  - `KPIsDTO` con: TotalMonthlyIncome, TotalMonthlyExpenses, NetCashFlow, SavingsRatio, TotalAccountsBalance, TotalDebtRemaining, DebtRatio, EmergencyCoverage
  - `ExpenseByCategoryDTO` con: CategoryName, CategoryIcon, CategoryColor, TotalAmount, Currency, Percentage
  - `DebtSummaryDTO` con: ID, Name, DebtType, RemainingAmount, TotalAmount, ProgressPercent, Currency
  - `AccountSummaryDTO` con: ID, Name, AccountType, Balance, Currency
  - `DashboardSummaryDTO` como respuesta raíz

- [ ] `internal/modules/dashboard/application/services/dashboard_service.go`
  - Inyectar repos: IncomeSourceRepo, ExpenseRepo, DebtRepo, AccountRepo
  - Ejecutar queries en **paralelo con goroutines + errgroup**
  - Calcular KPIs en base a `MonthlyEquivalent()` de cada entidad
  - Agregar gastos por categoría con JOIN a categories

- [ ] `internal/modules/dashboard/infrastructure/http/controllers/dashboard_controller.go`
  - `GET /dashboard/summary` → protegido con AuthMiddleware
  - Thin controller: extraer userID del JWT → llamar service → responder

- [ ] `internal/modules/dashboard/infrastructure/http/router.go`
  - Registrar la ruta con AuthMiddleware

- [ ] `cmd/server/dependencies.go`
  - Instanciar DashboardService y DashboardController
  - Inyectar los repos existentes del Finance Module (sin duplicar)

- [ ] `cmd/server/router.go`
  - Registrar rutas del Dashboard Module

*Nota: Una vez implementado, reemplazar el cálculo manual del frontend por este endpoint.*

### Base de Datos
- [ ] Crear migraciones SQL para producción
- [ ] Implementar seeds de datos de prueba

## 🟢 Mejoras Futuras

### Performance
- [ ] Implementar caché Redis para sesiones y datos frecuentes
- [ ] Implementar paginación con cursor (más eficiente que offset para listas grandes)

### Testing
- [ ] Tests de integración para repositorios GORM
- [ ] Tests E2E del API con Docker Compose

### Observabilidad
- [ ] Integrar structured logging (zerolog o zap)
- [ ] Agregar métricas con Prometheus
- [ ] Implementar tracing distribuido

---

*Última actualización: 2026-02-22*
