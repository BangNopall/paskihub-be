# High Risk Integrity and Security Batch 2 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix QA-005 through QA-013 with scoped authorization, private-file access, transaction locking, and fail-fast startup behavior.

**Architecture:** Services enforce ownership and relation invariants, repositories provide parent-bound predicates and PostgreSQL row locks, and startup functions return dependency errors. Sensitive files are served by authenticated resource endpoints while public branding assets retain static URLs.

**Tech Stack:** Go 1.24.3, Fiber v2, GORM, PostgreSQL, Redis, Swagger, `testing`

---

### Task 1: Event-Level Parent Binding

**Files:**
- Modify: `domain/contracts/event_contract.go`
- Modify: `internal/app/event/service/event_service.go`
- Modify: `internal/app/event/repository/event_repository.go`
- Test: `internal/app/event/service/event_service_test.go`

- [ ] Add failing tests for update/delete using a level from another event.
- [ ] Change repository update/delete signatures to accept `eventID`.
- [ ] Use `WHERE id = ? AND event_id = ?`, check `RowsAffected`, and return `domain.ErrNotFound` when zero.
- [ ] Persist the URL event ID instead of the body `event_id`.
- [ ] Run event service tests.

### Task 2: Assessment Relation Graph Validation

**Files:**
- Modify: `domain/contracts/assessment_contract.go`
- Modify: `internal/app/assessment/repository/assessment_repository.go`
- Modify: `internal/app/assessment/service/assessment_service.go`
- Test: `internal/app/assessment/service/assessment_service_test.go`

- [ ] Add failing tests for foreign event levels, mixed score inputs, and foreign award relations.
- [ ] Add repository methods that validate event-level, score-input, and award relation graphs.
- [ ] Reject invalid create payload relation graphs with `domain.ErrBadRequest`.
- [ ] Preserve `domain.ErrNotFound` for update/delete objects outside the URL event.
- [ ] Run assessment tests.

### Task 3: Banned Account Enforcement

**Files:**
- Modify: `domain/errors.go`
- Modify: `internal/app/user/service/user_service.go`
- Modify: `internal/middlewares/auth.go`
- Test: `internal/app/user/service/user_service_test.go`
- Test: `internal/middlewares/auth_test.go`

- [ ] Add failing tests for banned login and an existing token belonging to a newly banned user.
- [ ] Add `domain.ErrAccountBanned` mapped to 403.
- [ ] Reject banned users before token generation.
- [ ] Load the current user in authentication and reject banned accounts.
- [ ] Run user and middleware tests.

### Task 4: Authenticated Private File Gateway

**Files:**
- Create: `domain/contracts/private_file_contract.go`
- Create: `internal/app/private_file/controller/private_file_controller.go`
- Create: `internal/app/private_file/service/private_file_service.go`
- Create: `internal/app/private_file/repository/private_file_repository.go`
- Modify: `internal/infra/server/http_server.go`
- Modify: sensitive upload services and response mapping files.
- Test: `internal/app/private_file/service/private_file_service_test.go`

- [ ] Add failing authorization tests for participant, organizer, and admin file access.
- [ ] Store new sensitive uploads under `storage/private`.
- [ ] Add resource-ID based authenticated download routes.
- [ ] Convert sensitive response fields to authenticated API URLs without renaming JSON fields.
- [ ] Remove the broad `/public` static route and expose only event assets and team logos.
- [ ] Verify legacy paths resolve through database-authorized endpoints.

### Task 5: Atomic Team Approval and Top-Up Locking

**Files:**
- Modify: `domain/contracts/eo_team_contract.go`
- Modify: `internal/app/eo_team/service/eo_team_service.go`
- Modify: `internal/app/eo_team/repository/eo_team_repository.go`
- Modify: `internal/app/wallet/repository/wallet_repository.go`
- Test: `internal/app/eo_team/service/eo_team_service_test.go`
- Test: `internal/app/wallet/repository/wallet_repository_postgres_test.go`

- [ ] Add a failing service test proving approval delegates one atomic operation.
- [ ] Move registration lock, wallet lock, balance check, debit, transaction creation, and status update into one repository transaction.
- [ ] Reject registration states other than `WAITING`.
- [ ] Replace top-up query options with `clause.Locking`.
- [ ] Add opt-in PostgreSQL concurrency tests using `TEST_DATABASE_URL`.

### Task 6: Enum and Migration Reliability

**Files:**
- Modify: `internal/infra/database/pgsql_conn.go`
- Modify: `cmd/app/main.go`
- Test: `internal/infra/database/pgsql_conn_test.go`
- Test: `internal/infra/database/pgsql_conn_postgres_test.go`

- [ ] Add failing tests for migration and default-setting error propagation.
- [ ] Make `Migrate` and `EnsureDefaultSettings` return errors.
- [ ] Check every schema operation and add `KICKED` with an idempotent enum alteration.
- [ ] Fail startup before routes when migration fails.
- [ ] Add an opt-in PostgreSQL schema smoke test.

### Task 7: Mandatory Redis Startup

**Files:**
- Modify: `pkg/redis/redis.go`
- Modify: `internal/infra/server/http_server.go`
- Modify: `cmd/app/main.go`
- Test: `pkg/redis/redis_test.go`

- [ ] Add a failing constructor test for an unreachable Redis endpoint.
- [ ] Return `(RedisInterface, error)` from the constructor.
- [ ] Return route-mounting errors and fail startup before listening.
- [ ] Run Redis package and server compilation tests.

### Task 8: Documentation and Verification

**Files:**
- Modify generated Swagger only through `swag init -g cmd/app/main.go`.

- [ ] Run `gofmt` on all modified Go files.
- [ ] Run targeted unit tests after each task.
- [ ] Run opt-in PostgreSQL tests against the isolated test database.
- [ ] Run `go test ./...`.
- [ ] Run `go vet ./...`.
- [ ] Run `go build -o ./tmp/main ./cmd/app/main.go`.
- [ ] Regenerate and audit Swagger.
- [ ] Run `git diff --check` and verify QA-014 onward remain untouched.

