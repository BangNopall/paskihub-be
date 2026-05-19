# AGENTS.md

Standard context for Codex and other coding agents working on `paskihub-be`.

## Project Overview

Paskihub Backend is a Go service for managing Paskibra-related users, organizers, events, participants, assessments, scoring recaps, teams, wallets, and system settings.

Core stack:

- Go `1.24.3`
- Fiber v2 for HTTP routing
- GORM with PostgreSQL
- Redis for cache/session-related flows
- Viper for `.env` configuration
- Swagger via `swaggo/swag`
- Air for local live reload

The module path is:

```text
github.com/BangNopall/paskihub-be
```

## Important Commands

Run from the repository root.

```bash
# Start local development with live reload
air

# Run directly
go run cmd/app/main.go

# Build the app
go build -o ./tmp/main ./cmd/app/main.go

# Fresh migration: drops and recreates tables
go run cmd/app/main.go -fresh

# Run seeders
go run cmd/app/main.go -seeder

# Run a specific seeder
go run cmd/app/main.go -seeder -model <ModelName>

# Run tests when present
go test ./...

# Regenerate Swagger docs after changing API comments
swag init
```

Be careful with `-fresh`. It drops tables. In production, the app asks for confirmation, but agents should still treat this as destructive.

## Architecture

The project uses a modular, feature-based layout with Clean Architecture influence.

```text
cmd/app/main.go                  Application entry point
domain/                          Domain contracts, DTOs, entities, enums, errors
internal/app/<feature>/           Feature modules
internal/app/<feature>/controller Fiber handlers and route registration
internal/app/<feature>/service    Business logic
internal/app/<feature>/repository GORM data access
internal/infra/                   Database, environment, HTTP server wiring
internal/middlewares/             Fiber middlewares
pkg/                              Reusable helpers and integrations
docs/                             Generated Swagger files
data/seeders/                     Seeder data and logic
public/                           Static files served under /public
```

Key files:

- `cmd/app/main.go`: bootstraps env, PostgreSQL, migration/seeding, server, validation, and routes.
- `internal/infra/server/http_server.go`: central dependency wiring, route mounting, middleware mounting, cron job setup, and custom validation registration.
- `internal/infra/database/pgsql_conn.go`: PostgreSQL connection, enum creation, `AutoMigrate`, default settings, and seeders.
- `internal/infra/env/env.go`: `.env` loading and required environment shape.
- `domain/errors.go`: shared domain errors and HTTP status mapping.
- `GEMINI.md`: existing project context for Gemini CLI; keep this file and `AGENTS.md` conceptually aligned.

## Feature Conventions

When adding or changing a feature:

1. Put feature implementation under `internal/app/<feature>/`.
2. Keep the standard subpackages:
   - `controller`: Fiber handlers, request parsing, response writing, Swagger comments.
   - `service`: business rules, authorization-sensitive decisions, transaction orchestration, timeouts.
   - `repository`: GORM queries and persistence.
3. Define service and repository interfaces in `domain/contracts/`.
4. Put request/response structs in `domain/dto/`.
5. Put database models in `domain/entity/`.
6. Add shared constants/enums in `domain/constants/` or `domain/enums/`.
7. Wire repositories, services, controllers, and routes in `internal/infra/server/http_server.go`.
8. If a new table/entity is introduced, add it to `getInterfaces()` in `internal/infra/database/pgsql_conn.go`.

Follow the existing dependency direction:

```text
controller -> service -> repository -> database
```

Controllers should not perform database work directly. Repositories should not contain HTTP concerns.

## HTTP and Response Patterns

- Use Fiber handlers in controllers.
- Use `ctx.BodyParser`, `ctx.Params`, `ctx.Query`, and `ctx.Locals` consistently with existing controllers.
- Send responses through `pkg/helpers/http/response`.
- Map domain errors with `domain.GetCode(err)` where possible.
- Prefer returning existing domain errors from `domain/errors.go`; add new domain errors there when the API needs a stable error category.
- Preserve the current API prefix style, mostly `/api/v1/...`.
- Swagger is enabled at `/swagger/` in development mode.

## Validation

The project uses `github.com/go-playground/validator/v10`.

Custom validations are registered in `internal/infra/server/http_server.go`:

- `alphnumsympace`
- `plusnumeric`
- `validdate`
- `validpassword`

When adding DTO tags, reuse these validators when appropriate instead of inventing duplicate validation logic in controllers.

## Authentication, Authorization, and Middleware

- Global middleware includes CORS, static `/public`, and API key checking.
- Protected routes typically use `middleware.Authentication`.
- Role-specific guards include middleware such as `AuthAdmin`, `AuthPeserta`, and related methods on `internal/middlewares.Middleware`.
- Do not trust request-supplied user IDs for authorization. Compare against authenticated values from `ctx.Locals` or enforce role access in service/middleware.
- JWT logic lives in `pkg/jwt`.
- Password hashing lives in `pkg/bcrypt`.
- Redis-backed auth/session behavior lives through `pkg/redis` and middleware wiring.

## Database and Migrations

- GORM models live in `domain/entity/`.
- The app performs migrations on startup via `database.Migrate`.
- PostgreSQL enum types are created manually in `internal/infra/database/pgsql_conn.go`.
- Keep new enum values synchronized between:
  - PostgreSQL enum creation SQL
  - Go enum/constants definitions
  - entity field types
  - DTO validation/documentation
- Use GORM translated errors where available, such as `gorm.ErrDuplicatedKey`.
- Repository methods should log unexpected persistence errors and return domain errors.

## Environment

The app reads `.env` through Viper. `internal/infra/env/env.go` defines expected variables, including:

- `APP_ENV`, `APP_PORT`, `API_KEY`
- PostgreSQL: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`, `DB_NAME`
- Redis: `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASS`
- JWT: `JWT_SECRET_KEY`, `JWT_EXP_TIME`, role values
- Mail, Supabase, Firebase, and AWS-related variables

Do not commit real secrets. If documenting env requirements, create or update an example file rather than exposing local `.env` values.

## Swagger

When changing API behavior:

1. Update Swagger comments in the relevant controller.
2. Run `swag init`.
3. Review generated changes in:
   - `docs/docs.go`
   - `docs/swagger.json`
   - `docs/swagger.yaml`

Do not hand-edit generated Swagger files unless there is no viable generator path.

## Code Style

- Use `gofmt` on Go files before finishing.
- Keep imports organized by `gofmt`/`goimports` conventions.
- Prefer small, explicit functions that match existing repository/service/controller patterns.
- Keep comments useful. Swagger comments are expected for handlers; avoid noisy comments elsewhere.
- Preserve existing public API shapes unless the user explicitly asks for a breaking change.
- Avoid broad refactors while fixing a narrow issue.
- Do not rewrite unrelated files or generated artifacts.

## Testing and Verification

Current repository structure may not include many tests. For any code change, run the strongest feasible verification:

```bash
go test ./...
go build -o ./tmp/main ./cmd/app/main.go
```

If verification needs PostgreSQL, Redis, `.env`, or external services and they are unavailable, report that clearly and include what was run.

For Swagger-affecting changes, also run:

```bash
swag init
```

## Agent Workflow Rules

- Start by reading `GEMINI.md`, `AGENTS.md`, and the files directly related to the task.
- Use `rg`/`rg --files` for search.
- Check `git status --short` before editing. Treat pre-existing changes as user work and do not revert them.
- Make focused edits only.
- Use `apply_patch` for manual file edits.
- Never run destructive commands such as `git reset --hard`, `git checkout -- <file>`, table-dropping migrations, or production-impacting commands unless explicitly requested.
- When adding a feature, update contracts, DTOs, entity/migration wiring, route wiring, Swagger comments, and tests/verification together.
- When fixing a bug, first identify the data flow from route to service to repository, then patch at the layer that owns the behavior.
- Prefer existing helpers in `pkg/` over introducing new utilities.

## Known Project Notes

- The server uses a cron job to delete unverified users weekly.
- The app sets `BodyLimit` to `100 * 1024 * 1024`; the nearby comment currently says `10MB`, so verify intent before changing upload behavior.
- Some local utility scripts exist at the repository root for Swagger/controller fixes. Inspect them before relying on them; prefer native Go tooling for normal development.
- `.DS_Store`, binaries, and temporary build artifacts may exist locally. Do not include them in task changes.
