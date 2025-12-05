## Quick context

- Clean architecture Go HTTP service using PostgreSQL, sqlc for SQL generation, and Docker for orchestration.
- Project structure: `cmd/app/` (entry point), `internal/domain/` (entities/interfaces), `internal/usecase/` (business logic), `internal/infrastructure/` (DB/repos), `internal/api/handler/` (HTTP handlers), `pkg/` (shared utilities).
- Server listens on configurable port (default `:8080`), exposes RESTful endpoints for CRUD operations on authors.

## Architecture layers

- **Domain** (`internal/domain/author/`): Pure business logic. `entity.go` defines `Author` model; `repository.go` defines `Repository` interface (no implementation).
- **Use Cases** (`internal/usecase/author/`): Application logic orchestrating domain + repositories. Each file = one use case (list, get, create, update, delete); pattern: `New*UseCase(repo) → Execute(ctx, params)`.
- **Infrastructure** (`internal/infrastructure/`): DB drivers, persistence adapters. `repository/author.go` implements `Repository` using sqlc-generated `tutorial` queries. `database/postgres.go` wraps pgx connection.
- **API/Handlers** (`internal/api/handler/`): HTTP transport layer. `author.go` handles HTTP requests; calls use cases; returns JSON responses. Route registration via `RegisterRoutes(mux)`.
- **Config** (`config/`): Environment-based configuration; loaded in `main`.
- **Entry point** (`cmd/app/main.go`): Wires dependencies, starts server with graceful shutdown.

## What matters for code edits

- **Keep the interface contracts**: Domain interfaces in `internal/domain/author/repository.go` define the boundaries; implementations must satisfy them.
- **Use case pattern**: Each use case in `internal/usecase/author/` should be a small struct with a `repo` field and `Execute(ctx, ...)` method. No side effects outside Execute.
- **Handler pattern**: HTTP handlers in `internal/api/handler/author.go` decode request → call use case → encode response. Always set `Content-Type: application/json` first.
- **Errors**: Use `pkg/errors/` for domain errors or wrap `github.com/jackc/pgx` errors. Return HTTP status codes: 200 (OK), 201 (Created), 204 (No Content), 400 (Bad Request), 404 (Not Found), 500 (Internal Server Error).
- **Route registration**: Use Go 1.22+ http.ServeMux with `GET /authors`, `POST /authors`, `GET /authors/{id}`, `PUT /authors/{id}`, `DELETE /authors/{id}` patterns.
- **sqlc integration**: Use `tutorial.New(conn)` to get queries; repository adapts them to domain models.

## Data shapes / contract (important examples)

- **Author entity** (`internal/domain/author/entity.go`):
  ```go
  type Author struct {
    ID   int64
    Name string
    Bio  pgtype.Text
  }
  ```
- **Repository interface** (`internal/domain/author/repository.go`):
  ```go
  type Repository interface {
    GetAuthor(ctx context.Context, id int64) (*Author, error)
    ListAuthors(ctx context.Context) ([]*Author, error)
    CreateAuthor(ctx context.Context, params CreateAuthorParams) (*Author, error)
    UpdateAuthor(ctx context.Context, params UpdateAuthorParams) error
    DeleteAuthor(ctx context.Context, id int64) error
  }
  ```
- **HTTP request/response** (POST /authors):
  - Request: `{"name":"Alice","bio":"author"}`
  - Response (201): `{"id":1,"name":"Alice","bio":"author"}`

## Run / build / debug (concrete commands)

- **Development run** (from `cmd/app/`):
  - `go run ./cmd/app` (requires DATABASE_URL env var or default connection)
  - Default: `localhost:8080`, connects to `localhost:5432` (Postgres)

- **Build binary**:
  - `go build -o ./bin/app ./cmd/app`

- **Environment variables**:
  - `DATABASE_URL`: PostgreSQL connection string (default: `user=sqlc dbname=sqlc_db sslmode=disable host=localhost`)
  - `SERVER_PORT`: HTTP server port (default: `8080`)
  - `ENVIRONMENT`: deployment environment, e.g. `production` (default: `development`)

- **Docker/Make targets** (requires Docker):
  - `make db-up`: Start Postgres via docker-compose
  - `make db-wait`: Wait for Postgres to be healthy
  - `make db-down`: Stop Postgres and remove volume
  - `make lint`: Run golangci-lint in Docker
  - `make sqlc`: Run sqlc code generation in Docker

- **Basic runtime checks** (PowerShell examples):
  - Health check:
    - `curl.exe http://localhost:8080/health`
    - `Invoke-RestMethod http://localhost:8080/health`
  - List authors (GET):
    - `curl.exe http://localhost:8080/authors`
    - `Invoke-RestMethod http://localhost:8080/authors`
  - Create author (POST):
    - `curl.exe -X POST http://localhost:8080/authors -H "Content-Type: application/json" -d '{"name":"Alice","bio":{"String":"Researcher","Valid":true}}'`
    - `Invoke-RestMethod -Method Post -Uri http://localhost:8080/authors -Body '{"name":"Alice","bio":{"String":"Researcher","Valid":true}}' -ContentType 'application/json'`

## Patterns & conventions to preserve

- **Graceful shutdown**: the main function uses `signal.Notify` with a 5s `context.WithTimeout` and `srv.Shutdown(ctx)`. If you change server lifecycle code, keep graceful shutdown behavior and the timeout logic unless intentionally adjusting semantics.
- **Dependency injection**: Use constructor functions (`NewAuthorHandler`, `NewListAuthorsUseCase`, etc.) to inject dependencies. Avoid global state.
- **Clean architecture boundaries**: Domain logic lives in `internal/domain/` and never imports from other layers. Use Cases import Domain but not Infrastructure. Handlers import Use Cases. Infrastructure implements Domain interfaces.
- **Error handling**: Use custom errors from `pkg/errors/` or wrap sqlc/pgx errors. Always include context (status code, error message) in HTTP responses.
- **JSON handling**: Use `json.NewDecoder`/`json.NewEncoder` for streaming large responses. Use `json.Unmarshal`/`json.Marshal` for small payloads. Always set `Content-Type: application/json` first.
- **HTTP methods & status codes**:
  - `GET` → 200 (OK), 404 (Not Found), 500 (Internal Server Error)
  - `POST` → 201 (Created), 400 (Bad Request), 500 (Internal Server Error)
  - `PUT` → 204 (No Content), 400 (Bad Request), 404 (Not Found), 500 (Internal Server Error)
  - `DELETE` → 204 (No Content), 404 (Not Found), 500 (Internal Server Error)

## When writing new code

- Add tests in `_test.go` files using the standard `testing` package. For handlers, use `httptest.NewRequest` and `httptest.NewRecorder`. For use cases, mock the repository using interfaces.
- Run `go test ./...` to run all tests.
- Run `gofmt` / `go vet` as standard pre-commit checks.
- Follow the layered structure: domain entities → interfaces → use cases → handlers.

## Integration & dependencies

- **External packages**: Uses `github.com/jackc/pgx/v5` (database driver) and `github.com/jackc/pgtype` (PostgreSQL type mappings).
- **sqlc**: Uses sqlc-generated code in the `tutorial` package. Regenerate with `make sqlc` after modifying `query.sql`.
- **Configuration**: Environment-based via `config/config.go`. No hardcoded credentials.
- **Logging**: Uses the standard `log` package for now. Consider structured logging (e.g., `slog`) for production.

## Short checklist for pull requests

1. **Architecture**: Preserve layer boundaries. Domain code never imports infrastructure.
2. **Interfaces**: If modifying repositories or use cases, update the corresponding interface in `internal/domain/`.
3. **HTTP**: Set `Content-Type: application/json` first; use appropriate status codes; validate input.
4. **Tests**: Add unit tests for new use cases and handlers; use `httptest` for HTTP tests.
5. **SQL changes**: After modifying `query.sql`, run `make sqlc` to regenerate sqlc code.
6. **Linting**: Run `make lint` (Docker) or `gofmt` / `go vet` before commit.
7. **Graceful shutdown**: If changing server startup/shutdown, preserve the 5s timeout logic.

---

If any of the above is incorrect or you'd like the instructions framed differently (more detail on tests, CI, or packaging), tell me which parts to expand or change and I will iterate.
