# Critical Authorization Batch 1 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Enforce object ownership for QA-001 through QA-004 with stable 403 responses.

**Architecture:** Controllers pass authenticated IDs to services. Services own authorization decisions and use focused repository checks. Existing endpoint paths and response payloads remain intact.

**Tech Stack:** Go, Fiber v2, GORM, PostgreSQL, `testing`

---

### Task 1: Participant Team Ownership

**Files:**
- Modify: `internal/app/participant_team/service/participant_team_impl.go`
- Test: `internal/app/participant_team/service/participant_team_impl_test.go`
- Modify: `internal/app/participant_team/controller/participant_team_controller.go`

- [ ] Add failing tests proving foreign team detail and deletion return `domain.ErrForbidden`.
- [ ] Run the participant team service tests and confirm the authorization tests fail.
- [ ] Reuse authenticated institution lookup in detail/delete and compare `InstiId`.
- [ ] Map service errors through `domain.GetCode` in affected controllers.
- [ ] Run participant team tests.

### Task 2: Participant Payment Ownership

**Files:**
- Modify: `domain/contracts/participant_event_contract.go`
- Modify: `internal/app/participant_event/controller/participant_event_controller.go`
- Modify: `internal/app/participant_event/service/participant_event_impl.go`
- Test: `internal/app/participant_event/service/participant_event_impl_test.go`
- Test: `internal/app/participant_event/controller/participant_event_controller_test.go`

- [ ] Add a failing test proving a foreign registration is rejected before update.
- [ ] Run participant event tests and confirm failure.
- [ ] Pass authenticated user ID through the contract and controller.
- [ ] Check `GetRegistrationOwnership` before reading or saving a proof.
- [ ] Return `domain.ErrForbidden` and map it to 403.
- [ ] Run participant event tests.

### Task 3: Assessment Write Ownership

**Files:**
- Modify: `domain/contracts/form_penilaian_contract.go`
- Modify: `internal/app/assessment/controller/form_penilaian_controller.go`
- Modify: `internal/app/assessment/service/form_penilaian_service.go`
- Modify: `internal/app/assessment/repository/form_penilaian_repository.go`
- Test: `internal/app/assessment/service/form_penilaian_service_test.go`

- [ ] Add failing tests for denied score, violation, and finalize writes.
- [ ] Run assessment service tests and confirm failure.
- [ ] Add a repository authorization method covering registration, judge, subcategories, and violations.
- [ ] Pass authenticated organizer ID from controller to service.
- [ ] Reject denied operations with `domain.ErrForbidden` before transactions.
- [ ] Run assessment service tests.

### Task 4: Organizer Recap Ownership

**Files:**
- Modify: `domain/contracts/rekap_penilaian_contract.go`
- Modify: `internal/app/assessment/controller/rekap_penilaian_controller.go`
- Modify: `internal/app/assessment/service/rekap_penilaian_service.go`
- Modify: `internal/app/assessment/repository/rekap_penilaian_repository.go`
- Test: `internal/app/assessment/service/rekap_penilaian_service_test.go`

- [ ] Add failing tests for foreign registration/event-level read and publish operations.
- [ ] Run recap service tests and confirm failure.
- [ ] Add repository ownership checks for registration and event level.
- [ ] Pass organizer ID through all organizer recap service methods.
- [ ] Return `domain.ErrForbidden` before reads or writes.
- [ ] Run recap service tests.

### Task 5: Verification

**Files:**
- Verify all modified Go files and generated artifacts.

- [ ] Run `gofmt` on modified Go files.
- [ ] Run targeted tests.
- [ ] Run `go test ./...`.
- [ ] Run `go vet ./...`.
- [ ] Run `go build -o ./tmp/main ./cmd/app/main.go`.
- [ ] Confirm generated Swagger files and unrelated production files were not changed.
