# Backend to Frontend Implementation Alignment - Paskihub

Tanggal audit: 7 Juni 2026  
Backend repository: `paskihub-be`  
Frontend repository: `paskihub-fe` (`../paskihub-fe`)  
Status: **PARTIAL**

## 1. Executive Summary

Mayoritas endpoint backend sudah memiliki implementasi frontend dan header API key/Bearer
yang benar. Perbaikan frontend sebelumnya juga sudah menutup fallback rekening hardcoded,
environment example, body limit 500 MB, serta sebagian validasi upload.

Alignment belum lengkap. Blocker utama adalah file private yang dibuka langsung oleh
browser tanpa header autentikasi, enum `KICKED` yang belum diterima EO-team schema,
error/schema failure yang masih diubah menjadi empty state, dan support staff organizer
yang belum konsisten di backend.

## 2. Scope

Cross-check read-only dilakukan pada backend route/controller/service/repository/DTO/entity,
Swagger, serta frontend `src/services`, `src/actions`, `src/schemas`, auth, middleware,
environment example, upload actions, dan komponen yang membuka file.

Status:

- `VERIFIED`: contract ditemukan dan cocok secara statis.
- `NEEDS ADJUSTMENT`: implementasi ada tetapi tidak cocok/aman.
- `NOT FOUND`: endpoint frontend tidak ditemukan.
- `NOT VERIFIED`: butuh runtime/E2E.
- `NOT APPLICABLE`: tidak dibutuhkan frontend saat ini.

## 3. Backend Verification Summary

| Check | Result | Notes |
|---|---|---|
| `go test ./...` | PASS | PostgreSQL tests skip tanpa `TEST_DATABASE_URL` |
| `go build` | PASS | `./tmp/main` |
| `go vet ./...` | PASS | Tidak ada output |
| `go test -race ./...` | PASS | Hanya path test yang tersedia |
| Coverage | PASS, 6.4% | Controller/repository gap masih besar |
| Swagger check | PASS | Explicit `-g cmd/app/main.go` |

## 4. Global Contract Requirements for Frontend

- [x] Server-side integration memakai `API_BASE_URL`.
- [x] Authenticated request memakai `Authorization: Bearer <token>`.
- [x] Request backend memakai `x-api-key`.
- [ ] Production environment fail-fast saat `API_BASE_URL`/`API_KEY` kosong.
- [~] 400 ditampilkan sebagai validation/global error; belum konsisten.
- [~] 401 ditangani sebagai session failure; service belum punya typed global handler.
- [~] 403 dibedakan dari empty state pada beberapa service, belum menyeluruh.
- [~] 404 dibedakan dari empty state, belum menyeluruh.
- [~] 409 ditampilkan sebagai conflict, belum terstandardisasi.
- [~] 500 memiliki error state/retry, belum ada route error boundary menyeluruh.
- [ ] Backend/schema error tidak diubah menjadi empty state palsu.
- [x] Tidak ada fallback rekening hardcoded pada wallet service saat ini.
- [~] Upload divalidasi sebelum request; team flow masih lemah.
- [x] Backend tetap source of truth untuk ownership dan transaksi.

## 5. Auth and Session Alignment

| Backend Behavior | Frontend Required Behavior | Frontend Status | Notes |
|---|---|---|---|
| Banned user ditolak saat login | Tampilkan forbidden/banned message | NEEDS ADJUSTMENT | Auth.js mengubah semua login failure menjadi invalid credentials |
| Token lama banned user ditolak 403 | Logout/redirect dan tampilkan access denied | NEEDS CHECK | Belum ada global 403 handler |
| Bearer wajib | Semua protected service mengirim Bearer | VERIFIED | Ditemukan pada service/action protected |
| API key format `Key <key>` | Header harus menyertakan prefix `Key ` | NEEDS ADJUSTMENT | Frontend mengirim raw `API_KEY`; env harus sudah berisi prefixed value |
| Role hanya 3 enum | Schema dan middleware fail closed | NEEDS ADJUSTMENT | Login role masih `z.string()`; unknown role diarahkan ke peserta |
| Staff memiliki `parent_id` | Simpan effective organizer ID | VERIFIED | `organizerId` dibuat dari parent/id |

Catatan penting: backend `ApiKey` memerlukan nilai header `Key <api_key>`. Frontend
mengirim `process.env.API_KEY` apa adanya, sehingga nilai environment frontend harus
berisi prefix tersebut atau helper header perlu menambahkannya secara eksplisit.

## 6. Role and Ownership Alignment

| Area | Backend Rule | Frontend Expected Handling | Frontend Status | Notes |
|---|---|---|---|---|
| Participant Team | Hanya institution owner | Tampilkan 403, bukan empty | PARTIAL | HTTP 403 dilempar; schema error masih jadi empty/null |
| Participant Registration | Pelunasan/detail hanya owner | Tampilkan forbidden | PARTIAL | Service melempar error; runtime UX belum diuji |
| Organizer Event | Hanya event owner/effective parent | Disable action dan handle 403 | PARTIAL | Event mendukung parent; area lain belum |
| Organizer Assessment | Seluruh relasi satu event owner | Validasi payload dan handle 400/403 | PARTIAL | Payload banyak memakai `any` |
| Score Recap | Read/publish hanya owner | Handle 403 | PARTIAL | Service error belum typed |
| Staff Organizer | Gunakan effective parent | Dashboard/action staff konsisten | NEEDS ADJUSTMENT | Backend wallet/assessment/team/recap/dashboard masih staff ID |
| Admin | Global admin scope | Role guard di action/UI | PARTIAL | Backend kuat; action banyak hanya cek session |

## 7. API Endpoint Alignment Matrix

### Auth/User

| Method | Backend Path | Frontend Service/Action | Auth | Frontend Status | Notes |
|---|---|---|---|---|---|
| POST | `/api/v1/users/register/:role` | `auth.service`, `auth.actions` | API key | VERIFIED | Role EO/peserta |
| POST | `/api/v1/users/login` | `lib/auth.ts` | API key | PARTIAL | Banned detail disamarkan |
| POST | `/api/v1/users/logout` | `auth.service` | Bearer | VERIFIED | Client logout tetap lanjut saat backend gagal |
| GET | `/api/v1/users/verify-email/:email/:token` | `auth.service` | API key | VERIFIED | |
| POST | `/api/v1/users/forgot-password` | `auth.service` | API key | VERIFIED | |
| PUT | `/api/v1/users/reset-password/:token` | `auth.service` | API key | VERIFIED | |
| GET | `/api/v1/users/show/:userId` | Not clearly used | Bearer | NOT FOUND | |

### Admin

| Method | Backend Path | Frontend Service/Action | Auth | Frontend Status | Notes |
|---|---|---|---|---|---|
| GET/POST | `/api/v1/admin/admins` | `admin.service` | Admin | VERIFIED | |
| DELETE/POST | `/api/v1/admin/admins/:id[/reset-password]` | `admin.service` | Admin | VERIFIED | |
| GET | `/api/v1/admin/users` | `admin.service` | Admin | PARTIAL | Zod failure becomes `[]` |
| GET | `/api/v1/admin/users/:id` | `admin.service` | Admin | PARTIAL | Response not schema-validated |
| PUT | `/api/v1/admin/users/:id/{ban,unban,verify}` | `admin.service` | Admin | VERIFIED | |
| PUT | `/api/v1/admin/users/:id/{archive,unarchive}` | `admin.service` | Admin | VERIFIED | |
| PUT | `/api/v1/admin/events/:eventId/status` | `admin.service` | Admin | VERIFIED | |
| GET | `/api/v1/admin/dashboard` | `admin-dashboard.service` | Admin | VERIFIED | Zod validated |

### Organizer/Event/Team

| Method | Backend Path | Frontend Service/Action | Auth | Frontend Status | Notes |
|---|---|---|---|---|---|
| POST | `/api/v1/events/create` | `profile.service`, action | Organizer | VERIFIED | |
| GET | `/api/v1/events/user/:userId` | `profile.service` | Organizer | PARTIAL | Error can become `[]` |
| PUT/DELETE | `/api/v1/events/{update,delete}/:id` | `profile.service` | Organizer | VERIFIED | |
| POST | `/api/v1/events/upload/:id/{logo,poster}` | event actions | Organizer | PARTIAL | Frontend 5 MB; backend global 10 MB |
| POST/PUT/DELETE | `/api/v1/events/:id/levels[/:levelId]` | event/profile service | Organizer | VERIFIED | |
| GET | `/api/v1/eo/profile` | `profile.service` | Organizer | VERIFIED | |
| GET/POST/PUT/DELETE | `/api/v1/eo/staff/*` | profile/EO UI | Organizer owner | PARTIAL | Runtime matrix not tested |
| GET | `/api/v1/eo/events/:eventId/teams` | `eo-team.service` | Organizer | PARTIAL | Staff owner backend gap |
| GET/PUT | `/api/v1/eo/events/:eventId/teams/:regId/*` | `eo-team.actions` | Organizer | PARTIAL | `KICKED` schema missing |

### Participant/Team/Event

| Method | Backend Path | Frontend Service/Action | Auth | Frontend Status | Notes |
|---|---|---|---|---|---|
| GET/PUT | `/api/v1/peserta/profile` | `profile.service` | Peserta | PARTIAL | GET error returns null |
| PUT | `/api/v1/peserta/profile/security` | `profile.service` | Peserta | VERIFIED | |
| GET/POST | `/api/v1/peserta/teams` | `team.service/actions` | Peserta | PARTIAL | Team uploads lack action validation |
| GET/PUT/DELETE | `/api/v1/peserta/teams/:id` | `team.service/actions` | Peserta owner | PARTIAL | Schema error becomes null |
| GET | `/api/v1/peserta/events/open` | participant event service | Peserta | VERIFIED | |
| POST | `/api/v1/peserta/events/register` | participant event action | Peserta owner | VERIFIED | 5 MB client validation |
| PUT | `/api/v1/peserta/events/register/:id/pelunasan` | participant event action | Peserta owner | VERIFIED | |
| GET | `/api/v1/peserta/events/active` | participant event service | Peserta | VERIFIED | |
| GET | `/api/v1/peserta/events/registrations/:id` | participant event service | Peserta owner | VERIFIED | |
| GET | `/api/v1/peserta/assessment/recap/:regis_id` | participant event service | Peserta owner | VERIFIED | |
| GET | `/api/v1/peserta/rekap/scoreboard/:levelId` | participant event service | Peserta | VERIFIED | Published-only contract |

### Wallet/Transaction

| Method | Backend Path | Frontend Service/Action | Auth | Frontend Status | Notes |
|---|---|---|---|---|---|
| GET | `/api/v1/wallets/:eventId` | `wallet.service` | Organizer owner | PARTIAL | Staff owner backend gap |
| GET | `/api/v1/wallets/:eventId/logs` | `wallet.service` | Organizer owner | PARTIAL | Private proof URLs need proxy |
| POST | `/api/v1/wallets/:eventId/topup` | `wallet.actions` | Organizer owner | NEEDS ADJUSTMENT | Frontend permits PDF; backend rejects PDF |
| GET | `/api/v1/wallets/admin/transactions` | `admin.service` | Admin | VERIFIED | |
| PUT | `/api/v1/wallets/admin/transactions/:id/approve` | `admin.service/action` | Admin | VERIFIED | Backend idempotent |
| PUT | `/api/v1/wallets/admin/transactions/:id/reject` | `admin.service/action` | Admin | VERIFIED | |

### Assessment/Ranking/Score Recap

| Method | Backend Path | Frontend Service/Action | Auth | Frontend Status | Notes |
|---|---|---|---|---|---|
| CRUD | `/api/v1/eo/events/:eventId/assessment/judges` | `judge.service/actions` | Organizer owner | PARTIAL | Staff owner backend gap |
| CRUD | `.../violation-types` | `assessment.service/actions` | Organizer owner | PARTIAL | Request uses `any` |
| CRUD | `.../score-categories` | `assessment.service/actions` | Organizer owner | PARTIAL | Request uses `any` |
| CRUD | `.../score-sub-categories` | `assessment.service/actions` | Organizer owner | PARTIAL | Grade payload weakly typed |
| GET | `.../unified?level_id=` | `assessment.service` | Organizer owner | VERIFIED | |
| POST | `/api/v1/assessment/{scores/bulk,violations/bulk,finalize}` | assessment service/action | Organizer owner | PARTIAL | Finalize found; bulk endpoints not clearly used |
| CRUD | `.../awards` | `ranking.service/actions` | Organizer owner | VERIFIED | Zod response |
| GET/POST/PUT | `/api/v1/rekap/*` | `rekap.service/actions` | Organizer owner | PARTIAL | Staff owner backend gap |

### System Setting/Dashboard/File

| Method | Backend Path | Frontend Service/Action | Auth | Frontend Status | Notes |
|---|---|---|---|---|---|
| GET | `/api/v1/settings/public` | settings/wallet service | API key | VERIFIED | Wallet fails closed to null |
| GET/PATCH | `/api/v1/settings` | setting service/action | Admin | VERIFIED | Zod schema |
| GET | `/api/v1/organizer/dashboard` | organizer dashboard service | Organizer | PARTIAL | Staff owner backend gap |
| GET | `/api/v1/peserta/dashboard` | participant dashboard service | Peserta | VERIFIED | |
| GET | `/api/v1/public/home-stats` | home stats service | API key | VERIFIED | |
| GET | `/api/v1/files/:type/:id` | Direct browser URLs | Bearer + API key | NEEDS ADJUSTMENT | Direct `<img>/<a>` cannot send headers |

## 8. Status and Enum Alignment

| Backend Enum/Status | Meaning | Frontend Must Support? | Frontend Status | Notes |
|---|---|---|---|---|
| `WAITING` | Awaiting approval | Yes | VERIFIED | |
| `DP_PAID` | Down payment approved | Yes | VERIFIED | |
| `FULL_PAID` | Fully paid | Yes | VERIFIED | |
| `REJECTED` | Registration/payment rejected | Yes | VERIFIED | |
| `KICKED` | Team removed | Yes | NEEDS ADJUSTMENT | EO-team Zod enum excludes it |
| `PENDING` | Transaction/assessment pending | Yes | VERIFIED | Some schemas use free string |
| `APPROVE` | Transaction approved | Yes | VERIFIED | Note spelling is not `APPROVED` |
| `IN_PROGRESS` | Assessment started | Yes | PARTIAL | EO schema uses free string |
| `COMPLETED` | Assessment finalized | Yes | PARTIAL | EO schema uses free string |
| `DRAFT/OPEN/CLOSED/ARCHIVED` | Event lifecycle | Yes | PARTIAL | Several frontend schemas use free string |
| `ADMIN/ORGANIZER/PESERTA` | User role | Yes | NEEDS ADJUSTMENT | Login schema role is free string |
| `SD/SMP/SMA/PURNA/UMUM` | Institution type | Yes | PARTIAL | Often free string |
| `PASUKAN/DANPAS/OFFICIAL/PELATIH` | Team member role | Yes | PARTIAL | Team form role is free string |
| `TOPUP/WITHDRAW` | Wallet transaction type | Yes | PARTIAL | Not consistently schema-enforced |

## 9. Upload and File Access Alignment

| File Type | Backend Rule | Frontend Validation | Frontend Status | Notes |
|---|---|---|---|---|
| Event logo/poster | Public, UUID filename, 10 MB request cap | 5 MB, JPEG/PNG | PARTIAL | Backend lacks MIME/magic-byte check |
| Team logo | Public, UUID filename | Mainly UI `accept` | NEEDS ADJUSTMENT | Add action size/type validation |
| ID card | Private authenticated endpoint | Mainly UI `accept` | NEEDS ADJUSTMENT | Direct URL cannot send auth headers |
| Recommendation letter | Private authenticated endpoint | Mainly UI `accept` | NEEDS ADJUSTMENT | Direct URL cannot send auth headers |
| Registration proof | Private authenticated endpoint | 5 MB JPG/PNG/PDF | PARTIAL | Backend does not enforce matching type |
| Wallet proof | Private authenticated endpoint | 5 MB JPG/PNG/PDF | NEEDS ADJUSTMENT | Backend only accepts jpg/jpeg/png |

Required frontend design for private files:

1. Add a server-side route/handler that authenticates the NextAuth session.
2. Fetch backend `/api/v1/files/:type/:id` with API key and Bearer headers.
3. Stream the response with safe content type/disposition.
4. Point `<img>`/download links to that same-origin frontend route.
5. Handle 401/403/404 explicitly.

## 10. Wallet and Payment Alignment

- [x] Tidak ada fallback rekening hardcoded.
- [~] Settings failure returns `null`; UI must disable/hide top-up form.
- [x] Top-up double approval is rejected by backend transaction state.
- [x] Transaction statuses `PENDING/APPROVE/REJECTED` are represented.
- [~] Rejection reason is available and rendered.
- [ ] Private proof is fetched through authenticated proxy.
- [ ] PDF allowance is aligned between backend and frontend.
- [~] Double-submit UI prevention needs runtime verification.

## 11. Assessment, Ranking, and Score Recap Alignment

- [x] Backend rejects cross-event registration/judge/category/violation relations.
- [~] Frontend derives IDs from selected event/level, but many actions accept `any`.
- [x] 400/403 responses cause service/action failure.
- [ ] Typed UX distinguishes invalid relation (400) from ownership denial (403).
- [x] Grade mismatch is rejected by backend.
- [x] Recap publish/unpublish calls the correct backend endpoint.
- [~] Services generally throw on HTTP error, but error parsing often misses
  backend `error.details`.
- [ ] Staff organizer behavior is aligned after backend effective-owner fix.

## 12. Error Handling Alignment

| HTTP Status | Backend Meaning | Frontend Expected UX | Frontend Status | Notes |
|---|---|---|---|---|
| 400 | Invalid payload/API key format/business request | Field/global validation | PARTIAL | API key also returns 400 |
| 401 | Missing/expired/deny-listed token | Redirect login/session expired | PARTIAL | No common service handler |
| 403 | Role/ownership/banned | Access denied, not empty | PARTIAL | Some service paths preserve, some lose detail |
| 404 | Resource absent | Not found, not generic empty | PARTIAL | Several services return null/[] |
| 409 | Duplicate/conflict | Conflict message | NEEDS ADJUSTMENT | Not standardized |
| 413 | Request too large | File-size message | PARTIAL | Team flow parses non-JSON response |
| 429 | Rate limit | Retry-later UI | NOT VERIFIED | |
| 500 | Internal/dependency failure | Error boundary/retry | PARTIAL | No consistent route boundaries |

## 13. Environment Alignment

- [x] `.example-env` uses `API_BASE_URL`.
- [x] `API_KEY`, `NEXTAUTH_URL`, and `NEXTAUTH_SECRET` are documented.
- [x] Production fallback no longer points to localhost; it becomes empty string.
- [ ] Missing production env fails fast with a clear startup/build error.
- [ ] API key prefix contract is documented (`Key <api-key>`).
- [~] `NEXT_PUBLIC_API_BASE_URL` is used for assets; private assets must stop using it directly.
- [ ] Staging/prod values and secret ownership are documented outside source control.

## 14. Frontend QA Checklist Generated from Backend

### Public

- [ ] Register organizer and participant with valid/invalid payload.
- [ ] Login normal user; login banned user shows forbidden-specific message.
- [ ] Verify email and reset password expiry/error flows.
- [ ] Home stats backend failure is not shown as legitimate zero data.

### Admin

- [ ] User ban invalidates an already active session.
- [ ] Approve/reject top-up once; repeated action shows state conflict.
- [ ] Transaction proof loads through authenticated frontend proxy.
- [ ] Missing user/event returns not-found UX.

### Organizer

- [ ] Owner can create/update event and levels.
- [ ] Different organizer receives 403 for event/team/assessment/recap.
- [ ] Staff can access parent event features after backend fix.
- [ ] Wallet settings unavailable disables transfer/top-up UI.
- [ ] Wallet upload type matches backend.
- [ ] Assessment rejects mixed event IDs with 400.
- [ ] Recap publish/unpublish handles 403 distinctly.
- [ ] Team private documents load only for owner/admin.

### Peserta

- [ ] Team owner can view/edit/delete own team.
- [ ] Foreign team/detail/delete receives access denied.
- [ ] Foreign registration pelunasan/detail receives access denied.
- [ ] Registration and pelunasan upload size/type errors are clear.
- [ ] `KICKED` status renders intentionally.
- [ ] Private team/payment files load only through authenticated proxy.

### Negative/Error

- [ ] Missing API key and malformed `Key` prefix.
- [ ] Missing/malformed/expired Bearer token.
- [ ] 400 validation, 403 ownership, 404 missing, 409 conflict, 413 upload, 429 rate limit.
- [ ] Backend/Redis/settings unavailable.
- [ ] Unsupported extension, MIME mismatch, and magic-byte mismatch.
- [ ] Unknown role fails closed.

## 15. Gaps and Action Items

| Priority | Gap | Backend Evidence | Frontend Action Needed | Owner | Notes |
|---|---|---|---|---|---|
| P0 | Private files opened without auth headers | `/api/v1/files/*` requires API key + Bearer | Build authenticated same-origin proxy/blob flow | Frontend | Staging blocker |
| P0 | Staff organizer ownership inconsistent | Controllers use staff `id` in several domains | Hold staff QA until backend centralizes owner | Backend | Functional blocker |
| P1 | EO schema rejects `KICKED` | Backend enum includes `KICKED` | Add enum and UI state/test | Frontend | |
| P1 | Wallet PDF mismatch | Backend extension whitelist excludes PDF | Agree contract and align both sides | Both | |
| P1 | Team upload validation weak | Backend lacks content validation | Validate in action; backend validates independently | Both | |
| P1 | Error becomes empty/null | Backend has stable HTTP statuses | Introduce typed service error; remove false empty fallback | Frontend | |
| P1 | API key prefix implicit | Middleware requires `Key ` | Centralize header builder and env validation | Frontend | |
| P1 | PostgreSQL tests skipped | Tests require `TEST_DATABASE_URL` | Run in isolated CI service | Backend/DevOps | |
| P2 | Free-string enum schemas | Backend enums are closed sets | Convert role/status/type schemas to enums | Frontend | |
| P2 | Controller error mapping inconsistent | Some business errors return 500 | Normalize domain errors/status mapping | Backend | |

## 16. Release Recommendation

**Backend ready for controlled frontend QA, but not safe for staging integration yet.**

Staging gate:

1. Private-file frontend proxy implemented and tested.
2. Backend staff effective-owner behavior made consistent or staff features disabled.
3. `KICKED` and wallet upload contracts aligned.
4. PostgreSQL transaction/migration tests executed successfully.

## 17. Notes

Audit frontend bersifat statis/read-only. Tidak tersedia account fixture untuk tiga role,
server backend aktif, PostgreSQL test database, atau browser E2E authenticated. Status
`VERIFIED` berarti kontrak kode ditemukan dan cocok secara statis, bukan bukti runtime
production.
