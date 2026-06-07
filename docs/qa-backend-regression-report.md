# QA Backend Regression Report - Paskihub

Tanggal audit: 7 Juni 2026  
Repository: `paskihub-be`  
Status: **NEEDS FIX**

## 1. Executive Summary

Backend berhasil melewati `go test`, build, `go vet`, race detector, coverage run, dan
Swagger generation ke direktori sementara. Sebagian besar temuan Critical/High pada
`QA_REPORT.md` telah ditutup secara statis dan memiliki regression test.

Backend **layak dilanjutkan ke frontend QA terkontrol**, tetapi belum direkomendasikan
untuk staging penuh. Blocker tersisa adalah kontrak private-file yang belum kompatibel
dengan cara frontend membuka file, effective owner staff organizer yang belum konsisten,
dan validasi upload backend yang masih bergantung pada extension/global body limit.

Coverage meningkat dari 3.0% menjadi **6.4%**, tetapi controller, repository, wallet
service, system settings, dan beberapa flow organizer masih dominan 0%. Tiga integration
test PostgreSQL tersedia namun tidak dijalankan karena `TEST_DATABASE_URL` tidak tersedia.

## 2. Scope

Audit mencakup bootstrap, route, middleware, auth/JWT/Redis, ownership, event/team,
registration, wallet, assessment, ranking/rekap, settings, upload/private file,
migration/enum, cron, Swagger, test inventory, dan read-only cross-check `../paskihub-fe`.

Tidak dilakukan:

- Menjalankan server terhadap PostgreSQL/Redis nyata.
- Migration/seeder/fresh migration pada database lokal.
- Authenticated API E2E.
- Dynamic upload/magic-byte/path traversal test.
- Perubahan kode produksi atau frontend.

Perubahan existing sebelum audit adalah empat file log tanggal 7 Juni 2026. File tersebut
tidak diubah atau di-revert oleh audit.

## 3. Commands Executed

| Command | Result | Notes |
|---|---|---|
| `git status --short` | PASS | Empat file log sudah berubah sebelum audit |
| `git diff --stat` | PASS | Hanya perubahan log existing pada baseline |
| `go test ./...` | PASS | 22 test file, 33 fungsi test; test PostgreSQL bersyarat skip |
| `go build -o ./tmp/main ./cmd/app/main.go` | PASS | Binary berhasil dibuat |
| `go vet ./...` | PASS | Tidak ada output |
| `go test -race ./...` | PASS | Tidak ada race pada path yang dieksekusi |
| `go test ./... -coverprofile=./tmp/coverage.out` | PASS | Total statement coverage 6.4% |
| `go tool cover -func=./tmp/coverage.out` | PASS | Banyak controller/repository masih 0% |
| `swag init -g cmd/app/main.go --output /tmp/paskihub-swagger-check` | PASS | `docs.go`, JSON, YAML berhasil dibuat |
| Focused PostgreSQL tests | SKIP | `TEST_DATABASE_URL` tidak tersedia |

Toolchain lokal: Go `1.26.1`. Module menetapkan Go `1.24.3`; compatibility pada toolchain
Go 1.24.3 belum diuji.

## 4. Regression Verification Against Previous QA

| QA ID / Area | Status | Evidence | Notes |
|---|---|---|---|
| QA-001 participant team ownership | FIXED | `participant_team_impl.go`, ownership tests | Detail/delete membandingkan institution owner |
| QA-002 pelunasan ownership | FIXED | `GetRegistrationOwnership`, service/controller test | Auth user ID diteruskan |
| QA-003 assessment bulk/finalize ownership | FIXED | `ValidateAssessmentOwnership`, transaction re-check | Registration/judge/category/violation terikat event |
| QA-004 recap cross-organizer | FIXED | `RegistrationBelongsToOrganizer`, `EventLevelBelongsToOrganizer` | Read/publish diotorisasi |
| QA-005 event level cross-event | FIXED | Predicate `id AND event_id`, tests | Rows affected diperiksa |
| QA-006 assessment relation graph | FIXED | relation validation methods, tests | Event level/category/judge/registration dibatasi |
| QA-007 banned user | FIXED | login check, auth middleware DB lookup, tests | Login dan token lama ditolak |
| QA-008 sensitive static files | FIXED | `storage/private`, authenticated file endpoint | Logo/poster tetap public sesuai intent |
| QA-009 team approval concurrency | PARTIALLY FIXED | transaction + `clause.Locking`, integration test | Dynamic PostgreSQL test skip |
| QA-010 top-up concurrency | PARTIALLY FIXED | transaction + row locks, integration test | Dynamic PostgreSQL test skip |
| QA-011 enum `KICKED` | PARTIALLY FIXED | SQL sync + static test | PostgreSQL migration test skip |
| QA-012 Redis nil panic | FIXED | constructor returns error; route mount fails | Startup fail-fast |
| QA-013 migration errors | FIXED | errors returned and fatal at bootstrap | Unit failure-path tests pass |
| QA-014 event+wallet atomicity | FIXED | `CreateEventWithWallet` transaction | Static review |
| QA-015 upload filename/content | PARTIALLY FIXED | UUID filenames, private storage | No MIME/magic-byte validation; weak per-file limits |
| QA-016 DTO validation | PARTIALLY FIXED | validator added to major DTOs/controllers | Event DTOs and some multipart/business fields remain weak |
| QA-017 grade `Undefined` | FIXED | unmatched score returns error | No longer persisted as `Undefined` |
| QA-018 recap DB errors | FIXED | query `.Error` checked | Static review |
| QA-019 unverified-user cron | FIXED | cutoff `created_at < now-24h` | Boundary test still absent |
| QA-020 rows affected | PARTIALLY FIXED | user/event critical updates fixed | Not consistently enforced in all repositories |
| QA-021 internal error exposure | FIXED | 5xx details replaced with generic text | 4xx still exposes domain/validation detail by design |
| QA-022 staff effective owner | PARTIALLY FIXED | event/EO/private-file use parent | wallet, EO-team, assessment, recap, dashboard still use staff ID |
| QA-023 public settings route | FIXED | public GET only rate-limited | Frontend can call without Bearer |
| QA-024 Bearer parsing | FIXED | exact prefix check | Malformed token returns 400, missing token 401 |
| QA-025 Swagger command | FIXED | explicit `-g` command passes | `AGENTS.md` still lists plain `swag init` |
| QA-026 test log artifacts | NOT FIXED | existing log files changed during tests | Test logger still writes package-local logs |

## 5. Backend Route and Contract Review

All API routes pass global `x-api-key: Key <API_KEY>` except development Swagger and
the two public static asset directories. Protected routes additionally require
`Authorization: Bearer <token>`.

| Domain | Main paths | Auth / role | Ownership |
|---|---|---|---|
| Auth/User | `/api/v1/users/*` | Public auth flows; Bearer for logout/show | Show endpoint restricts non-admin self access |
| Admin | `/api/v1/admin/*` | Admin | Global admin scope |
| Event | `/api/v1/events/*` | Organizer; admin show route | Event owner; event level bound to parent |
| EO staff | `/api/v1/eo/profile`, `/staff` | Organizer | Effective parent owner |
| EO team | `/api/v1/eo/events/:eventId/teams/*` | Organizer | Event owner; staff parent not consistently used |
| Participant | `/api/v1/peserta/profile`, `/teams`, `/events` | Peserta | Institution/team/registration owner |
| Wallet | `/api/v1/wallets/*` | Organizer or admin | Event owner; staff parent not consistently used |
| Assessment | `/api/v1/eo/events/:eventId/assessment/*`, `/api/v1/assessment/*` | Organizer | Full relation ownership; staff parent gap |
| Recap/ranking | `/api/v1/rekap/*`, participant scoreboard | Organizer or peserta | Organizer ownership / published participant data |
| Settings | `/api/v1/settings/public`, `/api/v1/settings` | Public or admin | Admin mutation |
| Private files | `/api/v1/files/:resourceType/:resourceId` | Any authenticated role | Admin or participant/organizer resource owner |
| Dashboard | `/api/v1/{admin,organizer,peserta}/dashboard` | Matching role | Organizer staff parent gap |

Success responses use:

```json
{"code": 200, "status": "success", "message": "...", "data": {}}
```

Error responses use:

```json
{"error": {"code": "fail", "message": "...", "details": "..."}}
```

## 6. Authentication and Authorization Review

- API key format is `Key <api-key>`.
- JWT format is `Bearer <token>`.
- Middleware resolves the current user from PostgreSQL, so deleted/banned users are
  rejected even with a previously issued token.
- Redis deny-list errors fail closed with 500 instead of bypassing revocation.
- Role guards are fail closed for `ADMIN`, `ORGANIZER`, and `PESERTA`.
- Object ownership is now present on the critical participant and organizer flows.
- Staff organizer remains inconsistent because several controllers use `ctx.Locals("id")`
  instead of the authenticated `parent_id`.

## 7. Financial and Transaction Review

- Event and wallet creation are in one GORM transaction.
- Top-up approve/reject lock the transaction row and reject non-`PENDING` replay.
- Top-up approval locks the wallet row before balance update.
- Team approval locks registration and wallet, rejects replay, checks balance, debits,
  writes a withdrawal transaction, and updates registration atomically.
- System settings failure propagates; no backend fallback bank account is used.
- Dynamic concurrency assurance remains blocked until PostgreSQL tests run with
  `TEST_DATABASE_URL`.

## 8. Upload and File Access Review

Positive changes:

- Sensitive documents moved to `storage/private`.
- Server-generated UUID filenames replace client filenames.
- Sensitive paths returned to clients are logical `/api/v1/files/...` URLs.
- Private-file repository enforces participant/event organizer ownership.
- Public static exposure is limited to event assets and team logos.

Remaining gaps:

- Team, registration, and event uploads do not validate magic bytes or trusted MIME.
- Wallet validates extension only and currently rejects PDF although frontend permits PDF.
- No explicit per-file size limits below the global 10 MB request limit.
- Update flows can leave orphaned files after partial failure.

## 9. Validation and Error Response Review

- Core auth, settings, assessment bulk, participant profile, event update/level, and
  registration requests invoke validator.
- Some event DTO fields have no validation tags; invalid dates may still become zero time.
- Multipart team DTOs have no validator tags for required files/member roles.
- `PelunasanEventRequest.PaymentProof` is not validator-enforced.
- Several controllers map business validation/not-found errors to 500. In particular,
  participant registration/detail and some event/assessment paths do not consistently
  use `domain.GetCode`.
- Pagination defaults invalid values instead of rejecting them; status filters are not
  enum-validated at the controller boundary.

## 10. Database, Migration, Redis, and Cron Review

- Migration errors are propagated and stop startup.
- `KICKED` exists in Go and migration SQL, including `ADD VALUE IF NOT EXISTS`.
- Default settings creation is fail-fast.
- Redis is mandatory and startup fails when unavailable.
- Cron deletes only unverified users older than 24 hours.
- Startup migration still relies on `AutoMigrate`; versioned production migrations are
  recommended for schema history and rollback control.

## 11. Swagger and Documentation Review

Swagger generation passes with:

```bash
swag init -g cmd/app/main.go --output /tmp/paskihub-swagger-check
```

The checked controller comments cover the main routes. Some assessment and EO-team
comments omit possible `403` responses, and `AGENTS.md`/`GEMINI.md` still document
plain `swag init`, which is not the verified command.

## 12. Test Coverage and Gaps

- Test files: 22
- Test functions: 33
- Go packages: 65
- Statement coverage: 6.4%

Strongly covered examples include JWT, participant assessment, private-file service,
dashboard service, and targeted ownership services. Critical gaps remain:

- HTTP/controller contract tests for most domains.
- Wallet service and system setting tests.
- EO staff effective-owner matrix.
- Upload type/size/magic-byte tests.
- Cron cutoff boundary tests.
- PostgreSQL transaction/migration tests in CI.
- Authenticated API E2E across two participants and two organizers.

Recommended priorities:

- **P0:** Run PostgreSQL concurrency/migration tests in CI; private-file frontend E2E;
  staff ownership matrix; cross-tenant API tests.
- **P1:** Upload hardening tests; status/error mapping tests; wallet/settings integration.
- **P2:** Controller response schema, pagination/filter, cron boundary, and Swagger drift.

## 13. New Findings

| Severity | Issue | Location | Impact | Recommendation |
|---|---|---|---|---|
| High | Private file contract is not browser-direct compatible | Backend `/api/v1/files/*`; frontend direct `<img>/<a>` | Browser omits API key/Bearer, sensitive files fail to load | Frontend authenticated proxy/blob route; add E2E |
| Medium | Staff organizer effective owner remains partial | EO-team, wallet, assessment, recap, dashboard controllers | Staff receives forbidden/empty data for parent resources | Centralize effective owner middleware/helper |
| Medium | Upload content validation incomplete | Event/team/participant-event/wallet services | Extension spoofing and storage abuse remain possible | Validate size, MIME, magic bytes, extension |
| Medium | Status/error mapping inconsistent | Participant event and assorted controllers | 400/403/404 can surface as 500 | Use stable domain errors and `domain.GetCode` |
| Medium | Frontend EO team enum excludes `KICKED` | `paskihub-fe/src/schemas/eo-team.schema.ts` | Kicked registration can fail Zod parsing | Add `KICKED` and regression test |
| Low | Test execution writes repository logs | logger paths under packages | Dirty worktree/CI artifacts | Inject no-op/temp logger for tests |

## 14. Frontend Impact Summary

Frontend must:

- Fetch private files server-side or as authenticated blobs; direct URLs do not work.
- Add `KICKED` to EO team payment schemas and UI states.
- Keep wallet settings fail closed; current `null` behavior is improved but needs a
  clear unavailable state.
- Stop converting Zod/network failures into empty arrays/null where that changes meaning.
- Treat 401 as session expiry, 403 as forbidden/ownership/banned, and 409 as conflict.
- Use `parent_id`/`organizerId` consistently once backend staff ownership is unified.
- Align upload PDF support with backend; wallet backend currently accepts image extensions only.

## 15. Release Recommendation

- **Safe for frontend QA:** Yes, with known blockers documented.
- **Safe for staging:** No.
- **Needs backend fixes:** Yes, staff effective-owner consistency and upload hardening.
- **Needs frontend fixes:** Yes, authenticated private-file delivery and enum/error alignment.

## 16. Next Steps

1. Implement an authenticated frontend file proxy/blob flow for `/api/v1/files/*`.
2. Centralize backend effective owner resolution for all organizer routes.
3. Add shared backend upload validation and align allowed types with frontend.
4. Run PostgreSQL integration tests with `TEST_DATABASE_URL` in isolated CI.
5. Add API E2E for cross-user/cross-organizer denial, staff access, and private files.

