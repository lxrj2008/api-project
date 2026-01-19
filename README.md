# Go API Framework Skeleton

A production-ready reference skeleton for building Go (gin) APIs on SQL Server with fully hand-written SQL, JWT auth, structured logging, Swagger, and opinionated middleware.

## Features
- Gin router with CORS, rate limiting, request ID, access logging, recovery, and body size guardrails
- Configuration via Viper (YAML files + env overrides) with validation
- Structured logging (zap) and unified API responses with trace IDs
- SQL Server access through `database/sql` + `go-mssqldb`, repository + service layers with transaction boundaries
- Auth module with bcrypt password hashing and switchable HS256/RS256 JWT signing
- User CRUD sample (service + controller + DTOs) and Auth login endpoint
- Swagger UI at `/swagger` (doc template provided) and `/healthz` health probe
- SQLMock-based repository unit tests

## Project Layout
```
cmd/server/main.go     # bootstrap, graceful shutdown
internal/config/       # configuration, logger, DB initialization
server/                # gin engine + routes
controller/            # HTTP handlers + response helpers
service/               # business logic, transactions
repository/            # raw SQL data access
model/entity/          # DB entities
middleware/            # shared middleware
auth/                  # JWT + password hashing
utils/                 # helper functions (IDs, pagination, errors)
doc/                   # static Swagger template
dto/                   # request/response contracts
```

## Getting Started
1. **Install dependencies**
   ```bash
   go mod tidy
   ```
2. **Configure environment**
   - Copy `config.dev.yaml` as needed and/or set `APP_ENV` (defaults to `dev`).
   - Key settings:
     - `db.dsn`: e.g. `sqlserver://sa:123456@localhost:1433?database=go_api&encrypt=disable`
     - `auth.algorithm`: `HS256` (uses `jwtSecret`) or `RS256` (requires `privateKeyPath` & `publicKeyPath`).
3. **Run the server**
   ```bash
   go run ./cmd/server
   ```
4. **Available endpoints**
   - `GET /healthz`
   - `GET /swagger/*any`
   - `POST /api/v1/auth/login`
   - Authenticated (Bearer) user endpoints under `/api/v1/users`

## Testing
```bash
go test ./...
```
Repository tests rely on `github.com/DATA-DOG/go-sqlmock`. Add additional service/middleware tests as needed.

## Swagger
`docs/swagger.go` ships with a starter OpenAPI template. Update handler annotations and regenerate with `swag init -g cmd/server/main.go -o ./docs` when the schema evolves.

## Response Format
```
{
  "code": 0,
  "message": "OK",
  "data": {},
  "traceId": "..."
}
```
Errors reuse the same envelope with non-zero `code` and optional `details`.

## Notes
- Database schema expected: `go_api.dbo.users` (see specification)
- No migrations/ORM, caching, or containerization included per requirements
- Ensure secrets/DSNs are supplied securely via environment variables in production
