# Paskihub Backend QA and Testing Readiness Report

Tanggal audit: 7 Juni 2026  
Repository: `paskihub-be`  
Module: `github.com/BangNopall/paskihub-be`

## 1. Executive Summary

Status keseluruhan: **belum siap direkomendasikan untuk production release**.

Project dapat dikompilasi dan seluruh automated test yang tersedia berhasil dijalankan. `go vet`, race detector, dan regenerasi Swagger ke direktori sementara juga berhasil. Namun, tingkat keyakinan QA masih rendah karena statement coverage keseluruhan hanya **3.0%** dan sebagian besar controller, middleware, service, repository, migrasi, serta transaksi finansial belum memiliki test.

Audit statis menemukan beberapa risiko utama:

- Broken object-level authorization pada team peserta, pelunasan, assessment, rekap, dan event level.
- Endpoint assessment organizer dapat membaca atau mengubah data event organizer lain.
- User yang telah dibanned tetap dapat login dan menggunakan token.
- File sensitif seperti kartu identitas, bukti pembayaran, dan surat rekomendasi disimpan di direktori statis yang dipasang sebelum middleware API key.
- Approval team dan top-up memiliki risiko race condition atau saldo tidak konsisten.
- Ketidaksinkronan enum `KICKED` antara kode Go dan PostgreSQL.
- Kegagalan Redis menghasilkan dependency `nil` yang dapat menyebabkan panic pada endpoint terautentikasi.
- Error migrasi penting diabaikan sehingga aplikasi dapat tetap start dengan schema parsial.

## 2. Scope dan Metode

Audit dilakukan secara read-only terhadap kode produksi. Satu-satunya file repository yang ditambahkan adalah laporan ini.

Area yang diperiksa:

- Bootstrap aplikasi dan dependency wiring.
- Routing dan middleware.
- Authentication, role authorization, dan object ownership.
- Controller, service, repository, DTO, entity, dan error mapping.
- PostgreSQL migration dan enum.
- Redis/JWT/logout flow.
- Upload dan static file exposure.
- Transaction consistency pada wallet dan approval.
- Swagger readiness.
- Automated test inventory dan coverage.

Tidak dilakukan:

- Menjalankan aplikasi terhadap database lokal.
- Menjalankan migration atau seeder.
- Menjalankan `-fresh`.
- Mengubah generated Swagger.
- Pengujian end-to-end terhadap PostgreSQL, Redis, SMTP, atau deployment aktif.

Pengujian integrasi tersebut perlu environment terisolasi agar tidak berisiko mengubah data lokal.

## 3. Repository Inventory

Ringkasan struktur:

| Area | Lokasi |
|---|---|
| Entry point | `cmd/app/main.go` |
| HTTP wiring | `internal/infra/server/http_server.go` |
| Database/migration | `internal/infra/database/pgsql_conn.go` |
| Environment | `internal/infra/env/env.go` |
| Controllers | 15 file di `internal/app/*/controller/` |
| Services | 16 file di `internal/app/*/service/` |
| Repositories | 15 file di `internal/app/*/repository/` |
| DTO | 18 file di `domain/dto/` |
| Entity | 8 file di `domain/entity/` |
| Middleware | 6 file di `internal/middlewares/` |
| Test | 8 file, 11 fungsi test |
| Go package | 62 package |
| Go source non-generated docs | 118 file |

Arsitektur umumnya mengikuti alur `controller -> service -> repository -> database`. Tidak ditemukan controller yang melakukan query GORM secara langsung.

## 4. Verification Results

| Command | Hasil | Catatan |
|---|---|---|
| `git status --short` sebelum audit | PASS | Worktree bersih |
| `go test ./...` | PASS | Banyak package melaporkan `[no test files]` |
| `go build -o ./tmp/main ./cmd/app/main.go` | PASS | Binary berhasil dibuat |
| `go vet ./...` | PASS | Tidak ada output temuan |
| `go test -race ./...` | PASS | Tidak menemukan race pada test yang tersedia |
| `go test ./... -coverprofile=...` | PASS | Total statement coverage **3.0%** |
| `swag init` | FAIL | Mencari `main.go` di root |
| `swag init -g cmd/app/main.go --output /tmp/...` | PASS | `swagger.json` dan `swagger.yaml` identik; `docs.go` hanya berbeda pada nama package karena output directory |

Toolchain lokal adalah Go `1.26.1`, sedangkan `go.mod` menetapkan Go `1.24.3`. Build compatibility dengan toolchain Go 1.24.3 belum diverifikasi dalam audit ini.

## 5. Findings

### QA-001 - Peserta dapat melihat atau menghapus team milik peserta lain

**Risk:** Critical  
**Category:** Broken object-level authorization / data loss

Lokasi:

- `internal/app/participant_team/service/participant_team_impl.go:219`
- `internal/app/participant_team/service/participant_team_impl.go:255`
- `internal/app/participant_team/repository/participant_team_impl.go:47`
- `internal/app/participant_team/repository/participant_team_impl.go:79`

`GetTeamDetail` menerima `userID`, tetapi tidak memakainya untuk memverifikasi ownership. `DeleteTeam` juga tidak memverifikasi bahwa team dimiliki institution user yang sedang login.

Dampak:

- Peserta dapat membaca anggota, foto, kartu identitas, dan surat rekomendasi team lain jika mengetahui UUID.
- Peserta dapat menghapus team peserta lain selama kondisi payment check mengizinkan.

Rekomendasi:

- Query team harus selalu dibatasi oleh `team_id` dan owner user/institution.
- Tambahkan service-level ownership check untuk read, update, dan delete.
- Tambahkan negative authorization test dengan dua user berbeda.

### QA-002 - Pelunasan dapat mengubah registration milik peserta lain

**Risk:** Critical  
**Category:** Broken object-level authorization / payment integrity

Lokasi:

- `internal/app/participant_event/controller/participant_event_controller.go:126`
- `internal/app/participant_event/service/participant_event_impl.go:171`
- `internal/app/participant_event/repository/participant_event_impl.go:32`

Controller tidak meneruskan authenticated user ID ke service. Service mengambil registration hanya berdasarkan `regisID`, lalu mengganti bukti pembayaran dan status.

Dampak:

- Peserta dapat mengganti bukti pembayaran registration peserta lain.
- Status registration korban dapat dikembalikan menjadi `WAITING`.

Rekomendasi:

- Ubah contract agar `PelunasanEvent` menerima user ID.
- Verifikasi ownership registration sebelum menyimpan file atau mengubah status.
- Uji registration milik user lain harus menghasilkan `403`.

### QA-003 - Endpoint bulk assessment dan finalize tidak memiliki object authorization

**Risk:** Critical  
**Category:** Cross-tenant write / assessment integrity

Lokasi:

- `internal/app/assessment/controller/form_penilaian_controller.go:24`
- `internal/app/assessment/service/form_penilaian_service.go:29`
- `internal/app/assessment/service/form_penilaian_service.go:64`
- `internal/app/assessment/service/form_penilaian_service.go:86`

Role organizer diperiksa, tetapi authenticated user ID tidak diteruskan ke service. Service tidak memastikan registration, judge, violation type, subcategory, dan event berada dalam event milik organizer tersebut.

Dampak:

- Organizer dapat menambah atau menimpa score pada registration event lain.
- Organizer dapat menambah violation atau menandai assessment registration lain sebagai `COMPLETED`.
- ID dari beberapa event dapat dicampur sehingga integritas rekap rusak.

Rekomendasi:

- Resolve event dari registration dan verifikasi ownership sebelum transaksi.
- Verifikasi seluruh referenced ID berada pada event dan event level yang sama.
- Letakkan seluruh validasi dan write dalam transaksi yang sama.

### QA-004 - Rekap organizer dapat membaca dan mem-publish event level milik organizer lain

**Risk:** Critical  
**Category:** Cross-tenant read/write

Lokasi:

- `internal/app/assessment/controller/rekap_penilaian_controller.go:27`
- `internal/app/assessment/service/rekap_penilaian_service.go:20`
- `internal/app/assessment/service/rekap_penilaian_service.go:50`
- `internal/app/assessment/repository/rekap_penilaian_repository.go:181`

Endpoint detail, scoreboard, custom leaderboard, dan publish hanya memeriksa role organizer. `PublishScoreboard` menerima `userID`, tetapi service mengabaikannya dan langsung mengubah event level berdasarkan ID.

Dampak:

- Organizer dapat membaca nilai detail dan leaderboard event lain.
- Organizer dapat publish atau unpublish scoreboard event lain.

Rekomendasi:

- Derive event owner dari registration/event level pada repository.
- Wajibkan ownership check pada semua endpoint rekap organizer.
- Tambahkan test lintas dua organizer.

### QA-005 - Event level dapat dipindahkan atau dihapus lintas event

**Risk:** High  
**Category:** Broken object-level authorization

Lokasi:

- `internal/app/event/service/event_service.go:338`
- `internal/app/event/service/event_service.go:379`
- `internal/app/event/repository/event_repository.go:210`
- `internal/app/event/repository/event_repository.go:221`

Service hanya memastikan organizer memiliki `eventId` pada URL. Tidak ada verifikasi bahwa `levelId` memang milik event tersebut.

Pada update, entity disimpan dengan `levelId` target dan `EventId` dari event milik attacker. Pada delete, repository menghapus hanya berdasarkan `levelId`.

Dampak:

- Organizer dapat memindahkan event level event lain ke event miliknya.
- Organizer dapat menghapus event level event lain.

Rekomendasi:

- Update/delete dengan predicate gabungan `id = ? AND event_id = ?`.
- Jangan menggunakan `Save` untuk object yang ownership-nya belum dibuktikan.

### QA-006 - Assessment CRUD menerima relasi lintas event

**Risk:** High  
**Category:** Tenant boundary / relational integrity

Lokasi:

- `internal/app/assessment/service/assessment_service.go:113`
- `internal/app/assessment/service/assessment_service.go:294`
- `internal/app/assessment/service/assessment_service.go:441`
- `internal/app/assessment/service/assessment_service.go:490`

Ownership event pada URL diperiksa, tetapi referenced IDs pada payload tidak seluruhnya dibuktikan berada di event tersebut. Contohnya:

- `EventLevelID` saat membuat violation/category.
- `ScoreCategoriesID` saat membuat subcategory.
- `RegisID`, `JudgesID`, dan `SubCategoryID` saat input score.
- Event level dan score category pada award.

Dampak:

- Data assessment antar-event dapat tercampur.
- Nilai dapat dimasukkan menggunakan judge atau rubric dari event lain.
- Award dapat mereferensikan category/level tenant lain.

Rekomendasi:

- Validasi graph relasi lengkap sebelum write.
- Gunakan repository query yang mengikat setiap ID ke `event_id`.
- Tambahkan database constraint yang relevan bila model memungkinkan.

### QA-007 - User yang dibanned tetap dapat login dan menggunakan token

**Risk:** High  
**Category:** Authentication/account enforcement

Lokasi:

- `internal/app/user/service/user_service.go:186`
- `internal/middlewares/auth.go:15`
- `internal/app/user/service/user_service.go:507`

Flag `IsBanned` dapat diubah oleh admin, tetapi tidak diperiksa pada login atau middleware authentication. Middleware juga tidak mengambil status user dari database.

Dampak:

- Ban tidak menghentikan login baru.
- Token lama user yang dibanned tetap dapat mengakses API.

Rekomendasi:

- Tolak login saat `IsBanned`.
- Periksa status akun pada authentication atau gunakan revocation/versioning session.
- Invalidasi session aktif ketika admin melakukan ban.

### QA-008 - File sensitif tersedia melalui static route sebelum API key

**Risk:** High  
**Category:** Sensitive data exposure

Lokasi:

- `internal/infra/server/http_server.go:117`
- `internal/infra/server/http_server.go:118`
- `internal/infra/server/http_server.go:119`
- `internal/app/participant_team/service/participant_team_impl.go:75`
- `internal/app/participant_event/service/participant_event_impl.go:152`
- `internal/app/wallet/service/wallet_service.go:185`

`/public` dipasang sebelum middleware API key. Direktori tersebut menyimpan kartu identitas, foto, surat rekomendasi, dan bukti pembayaran.

Dampak:

- File berpotensi dapat diakses tanpa API key atau JWT jika URL diketahui atau bocor.
- Data pribadi dan dokumen pembayaran dapat terekspos.

Rekomendasi:

- Simpan dokumen sensitif di storage private.
- Sajikan melalui endpoint terautentikasi dengan ownership check atau signed URL berumur pendek.
- Pisahkan asset publik seperti logo/poster dari dokumen privat.

### QA-009 - Approval team memiliki risiko lost update atau saldo negatif saat concurrent

**Risk:** High  
**Category:** Financial consistency / concurrency

Lokasi:

- `internal/app/eo_team/service/eo_team_service.go:150`
- `internal/app/eo_team/service/eo_team_service.go:151`
- `internal/app/eo_team/service/eo_team_service.go:156`
- `internal/app/eo_team/service/eo_team_service.go:161`

Transaction dibuka melalui `s.db`, tetapi wallet dibaca melalui `s.walletRepo` yang masih menggunakan base DB, bukan transaction handle. Tidak ada row lock pada pembacaan saldo.

Dampak:

- Dua approval paralel dapat sama-sama lolos balance check.
- Update saldo dapat hilang atau menghasilkan debit yang tidak konsisten.
- Registration yang sama juga belum memiliki idempotency/state transition guard yang kuat.

Rekomendasi:

- Gunakan transaction-bound repository.
- Lock wallet dan registration dengan PostgreSQL `FOR UPDATE`.
- Gunakan atomic conditional update dan uji dua approval paralel.

### QA-010 - Locking top-up perlu dikonfirmasi dengan integration test PostgreSQL

**Risk:** High  
**Category:** Financial consistency / concurrency  
**Status:** Perlu dikonfirmasi secara dinamis

Lokasi:

- `internal/app/wallet/repository/wallet_repository.go:46`
- `internal/app/wallet/repository/wallet_repository.go:50`
- `internal/app/wallet/repository/wallet_repository.go:68`

Kode menggunakan `Set("gorm:query_option", "FOR UPDATE")`. Efektivitas pola ini pada GORM v2/PostgreSQL perlu dibuktikan; pola yang eksplisit adalah `clause.Locking`.

Dampak jika lock tidak diterapkan:

- Approval transaksi paralel dapat membaca saldo lama.
- Jaminan idempotency dan serialisasi approval melemah.

Rekomendasi:

- Verifikasi SQL aktual dengan logger/test database.
- Gunakan `Clauses(clause.Locking{Strength: "UPDATE"})`.
- Tambahkan concurrency integration test.

### QA-011 - Enum `KICKED` tidak dibuat pada PostgreSQL

**Risk:** High  
**Category:** Schema mismatch

Lokasi:

- `domain/enums/enums.go:36`
- `internal/app/eo_team/service/eo_team_service.go:216`
- `internal/infra/database/pgsql_conn.go:135`

Go mendefinisikan dan menulis status `KICKED`, tetapi enum PostgreSQL `registration_status` hanya berisi `WAITING`, `DP_PAID`, `FULL_PAID`, dan `REJECTED`.

Dampak:

- Endpoint kick team dapat gagal dengan PostgreSQL enum error.
- Error kemungkinan dikembalikan sebagai status yang tidak konsisten.

Rekomendasi:

- Tambahkan migration aman untuk enum `KICKED`.
- Tambahkan schema integration test yang mengeksekusi kick flow.

### QA-012 - Redis unavailable menghasilkan dependency nil dan potensi panic

**Risk:** High  
**Category:** Availability / nil dependency

Lokasi:

- `pkg/redis/redis.go:23`
- `pkg/redis/redis.go:30`
- `pkg/redis/redis.go:34`
- `internal/infra/server/http_server.go:135`
- `internal/middlewares/auth.go:62`
- `internal/app/user/service/user_service.go:477`

Jika ping Redis gagal, constructor mengembalikan `nil`, tetapi server tetap dirakit dan start. Authentication dan logout kemudian memanggil method pada interface tersebut.

Dampak:

- Request terautentikasi atau logout dapat panic saat Redis down.
- Health/readiness aplikasi tidak mencerminkan dependency yang gagal.

Rekomendasi:

- Fail fast saat Redis mandatory tidak tersedia, atau gunakan explicit degraded-mode implementation.
- Constructor sebaiknya mengembalikan `(RedisInterface, error)`.
- Tambahkan startup dan middleware test untuk Redis failure.

### QA-013 - Migration dan default setting mengabaikan error

**Risk:** High  
**Category:** Deployment reliability / schema integrity

Lokasi:

- `internal/infra/database/pgsql_conn.go:81`
- `internal/infra/database/pgsql_conn.go:87`
- `internal/infra/database/pgsql_conn.go:165`
- `internal/infra/database/pgsql_conn.go:172`
- `internal/infra/database/pgsql_conn.go:184`

Error dari `DropTable`, enum creation, `AutoMigrate`, count, dan create default setting tidak diperiksa.

Dampak:

- Aplikasi dapat start dengan schema parsial atau enum belum tersedia.
- Default setting dapat tidak ada, lalu wallet flow gagal.
- Deployment terlihat sehat meskipun migration gagal.

Rekomendasi:

- Periksa setiap error dan hentikan startup bila schema tidak valid.
- Gunakan versioned migration, bukan hanya startup `AutoMigrate`.
- Tambahkan migration smoke test pada database kosong dan database versi sebelumnya.

### QA-014 - Create event dan wallet tidak atomic

**Risk:** Medium-High  
**Category:** Transaction consistency

Lokasi:

- `internal/app/event/service/event_service.go:87`
- `internal/app/event/service/event_service.go:103`

Event dibuat lebih dahulu, kemudian wallet dibuat melalui operasi terpisah tanpa transaction.

Dampak:

- Jika wallet creation gagal, event tetap tersimpan tanpa wallet.
- Flow berikutnya dapat gagal karena asumsi satu event memiliki satu wallet.

Rekomendasi:

- Buat event dan wallet dalam satu database transaction.
- Tambahkan unique constraint pada `wallets.event_id`.

### QA-015 - Upload menggunakan filename dari client dan validasi content lemah

**Risk:** Medium-High  
**Category:** File upload security  
**Status:** Path traversal perlu dikonfirmasi dengan dynamic test

Lokasi:

- `internal/app/event/service/event_service.go:178`
- `internal/app/participant_team/service/participant_team_impl.go:40`
- `internal/app/participant_event/service/participant_event_impl.go:40`
- `internal/app/wallet/service/wallet_service.go:173`
- `internal/infra/server/http_server.go:85`

Nama file client dimasukkan ke path penyimpanan. Sebagian besar upload tidak memvalidasi extension, MIME, signature, atau ukuran per jenis file. Wallet hanya memeriksa extension. Global body limit adalah 100 MB meskipun komentar menyebut 10 MB.

Dampak:

- Upload file non-image atau file berbahaya ke static directory.
- Potensi path manipulation bergantung pada normalisasi filename oleh multipart stack.
- Resource exhaustion melalui upload besar.

Rekomendasi:

- Abaikan filename asli; generate nama server-side dan gunakan basename hanya untuk metadata.
- Validasi MIME dan magic bytes.
- Terapkan size limit per endpoint dan storage quota.
- Simpan file sensitif di luar static web root.

### QA-016 - Validasi request tidak konsisten dan sebagian tag tidak pernah dijalankan

**Risk:** Medium  
**Category:** Input validation

Lokasi:

- `domain/dto/user_dto.go:10`
- `domain/dto/event_dto.go:10`
- `internal/app/event/controller/event_controller.go:76`
- `internal/app/system_setting/controller/system_setting_controller.go:131`
- `internal/app/user/controller/user_controller.go:77`

Sebagian DTO memakai tag `binding`, sementara project menggunakan `validator` dengan tag default `validate`. Banyak controller hanya `BodyParser` tanpa menjalankan validator.

Contoh dampak:

- Admin dapat dibuat dengan field kosong bila service/database tidak menolak.
- Update setting menerima zero value; `CoinRate = 0` dapat menyebabkan pembagian dengan nol pada wallet response.
- Tanggal event invalid diabaikan dan berubah menjadi zero time.
- Event status dapat ditulis tanpa validasi enum pada beberapa flow.

Rekomendasi:

- Standardisasi tag `validate`.
- Jalankan validator pada seluruh request DTO.
- Tambahkan business validation untuk date order, team member limit, payment type, enum, dan monetary values.

### QA-017 - Nilai relasi assessment dapat disimpan dengan grade `Undefined`

**Risk:** Medium  
**Category:** Scoring correctness

Lokasi:

- `internal/app/assessment/service/form_penilaian_service.go:41`
- `internal/app/assessment/service/form_penilaian_service.go:48`
- `internal/app/assessment/service/form_penilaian_service.go:148`

Service hanya memeriksa bahwa hasil query rules tidak kosong, bukan bahwa setiap subcategory memiliki rule dan setiap score masuk ke range. `getGrade` mengembalikan `"Undefined"` dan hasil tersebut tetap disimpan.

Dampak:

- Score valid secara request dapat menghasilkan grade invalid.
- Rekap dan ranking dapat mengandung data yang tidak sesuai rubric.

Rekomendasi:

- Pastikan setiap subcategory memiliki rules.
- Tolak gap/overlap range dan score yang tidak cocok.
- Tambahkan unit test batas bawah, batas atas, gap, overlap, dan missing rules.

### QA-018 - Error database pada rekap detail diabaikan

**Risk:** Medium  
**Category:** Silent partial response

Lokasi:

- `internal/app/assessment/repository/rekap_penilaian_repository.go:46`
- `internal/app/assessment/repository/rekap_penilaian_repository.go:80`

Error dari query flat scores dan team violations tidak diperiksa.

Dampak:

- API dapat mengembalikan sukses dengan score atau violation kosong/parsial.
- Gangguan database sulit dideteksi dari response.

Rekomendasi:

- Periksa `.Error` pada seluruh query.
- Jangan mengembalikan response sukses jika komponen rekap gagal dimuat.

### QA-019 - Cron menghapus semua user unverified tanpa batas usia

**Risk:** Medium-High  
**Category:** Data retention / unintended deletion

Lokasi:

- `internal/infra/server/http_server.go:208`
- `internal/app/user/service/user_service.go:238`
- `internal/app/user/repository/user_repository.go:42`

Cron mingguan menjalankan delete untuk seluruh user dengan `email_is_verified = false`, tanpa filter umur atau expiry token.

Dampak:

- User yang baru registrasi menjelang jadwal cron dapat langsung terhapus.
- Retention behavior tidak sesuai dengan expiry token satu jam secara eksplisit.

Rekomendasi:

- Hapus hanya user unverified yang dibuat atau expired sebelum cutoff tertentu.
- Log jumlah record dan tambahkan test boundary waktu.

### QA-020 - Ban, verify, dan update lain tidak memastikan row ditemukan

**Risk:** Medium  
**Category:** API correctness

Lokasi:

- `internal/app/user/repository/user_repository.go:107`
- `internal/app/user/repository/user_repository.go:118`
- `internal/app/event/repository/event_repository.go:232`

Update hanya memeriksa `.Error`, tidak `RowsAffected`.

Dampak:

- Endpoint dapat mengembalikan sukses untuk user/event yang tidak ada.
- QA dan client tidak dapat membedakan update sukses dan no-op.

Rekomendasi:

- Periksa `RowsAffected`.
- Kembalikan `ErrNotFound` saat target tidak ada.

### QA-021 - Error response dapat mengekspos detail internal

**Risk:** Medium  
**Category:** Information disclosure

Lokasi:

- `pkg/helpers/http/response/response.go:45`
- `pkg/helpers/http/response/response.go:54`

`SendErrResp` selalu memasukkan `err.Error()` ke response. Banyak repository mengembalikan raw GORM/database error.

Dampak:

- Client dapat menerima detail schema, constraint, path, atau dependency internal.

Rekomendasi:

- Log detail internal di server.
- Response production hanya mengirim stable public error code/message.

### QA-022 - Staff organizer tidak konsisten menggunakan parent ownership

**Risk:** Medium  
**Category:** Authorization model / functional correctness

Lokasi:

- `internal/middlewares/auth.go:89`
- `internal/app/event/controller/event_controller.go:166`
- `internal/app/event/controller/event_controller.go:218`
- `internal/app/eo/controller/eo_controller.go:36`

Token staff membawa `parent_id`, tetapi hanya sebagian flow menggunakannya. Banyak operasi event, wallet, assessment, dashboard, dan EO staff management menggunakan staff's own `id`.

Dampak:

- Staff dapat gagal mengelola resource parent yang seharusnya diizinkan.
- Staff dapat membuat resource terpisah atas ID sendiri.
- Staff dapat mencoba mengelola staff di bawah dirinya sendiri.

Rekomendasi:

- Definisikan satu konsep `effective_owner_id`.
- Bedakan permission owner dan staff secara eksplisit.
- Tambahkan matrix test owner/staff/organizer lain/admin.

### QA-023 - Public settings route tidak sesuai nama dan Swagger

**Risk:** Low-Medium  
**Category:** API contract mismatch

Lokasi:

- `internal/app/system_setting/controller/system_setting_controller.go:29`
- `internal/app/system_setting/controller/system_setting_controller.go:45`

Route bernama `/public` dan Swagger hanya menyatakan API key, tetapi implementasi mewajibkan JWT organizer.

Dampak:

- Client mengikuti Swagger akan menerima unauthorized.
- Peserta atau public client tidak dapat mengambil setting yang disebut public.

Rekomendasi:

- Konfirmasi intended audience.
- Selaraskan route middleware, nama endpoint, dan Swagger.

### QA-024 - Authentication header parsing tidak mewajibkan skema Bearer

**Risk:** Low  
**Category:** Protocol correctness

Lokasi:

- `internal/middlewares/auth.go:32`
- `internal/middlewares/auth.go:45`

Middleware menerima token dari elemen kedua tanpa memeriksa elemen pertama adalah `Bearer`.

Rekomendasi:

- Wajibkan tepat dua elemen dan prefix `Bearer`.
- Tambahkan test malformed header.

### QA-025 - `swag init` pada dokumentasi project tidak dapat dijalankan apa adanya

**Risk:** Low  
**Category:** Developer readiness

Perintah `swag init` gagal karena entry point berada di `cmd/app/main.go`. Perintah berikut berhasil dan menghasilkan output identik:

```bash
swag init -g cmd/app/main.go
```

Rekomendasi:

- Perbarui dokumentasi command menjadi bentuk yang eksplisit.
- Tambahkan Swagger generation check pada CI.

### QA-026 - Test JWT menulis log ke direktori package

**Risk:** Low  
**Category:** Test isolation / repository hygiene

Saat test dijalankan, logger membuat `pkg/jwt/data/logs/<date>.log` karena working directory test adalah package directory.

Dampak:

- Test mengotori worktree dengan file untracked.
- CI atau developer dapat mendapat artefak berbeda berdasarkan package yang menjalankan logger.
- Log test memperlihatkan warning konfigurasi JWT kosong karena environment sengaja dilewati pada binary `.test`.

Rekomendasi:

- Arahkan logger test ke temporary directory atau no-op test logger.
- Inject konfigurasi JWT pada test secara eksplisit.
- Tambahkan ignore hanya sebagai pertahanan tambahan, bukan pengganti test isolation.

## 6. Automated Test Readiness

Total statement coverage: **3.0%**.

Package dengan coverage berarti:

- `internal/app/dashboard/service`: 50.0%
- `internal/app/participant_assessment/service`: 70.8%
- `internal/app/participant_event/controller`: 27.0%
- `internal/app/participant_event/service`: 27.3%
- `pkg/html`: 60.0%
- `pkg/jwt`: 76.9%
- `internal/app/assessment/repository`: 1.3%
- `internal/infra/database/seeders`: 1.9%

Area kritis dengan coverage 0%:

- Seluruh middleware.
- User controller/service/repository.
- Event controller/service/repository.
- Wallet controller/service/repository.
- EO team controller/service/repository.
- Assessment controller dan service utama.
- Form assessment service.
- Rekap service.
- Participant team controller/service/repository.
- System setting.
- Migration dan HTTP server wiring.
- Redis helper dan response helper.

Race detector lulus, tetapi hasil tersebut hanya mencakup path yang dijalankan oleh 11 test saat ini. Race pada transaksi database paralel belum diuji.

## 7. Recommended Test Plan

### Priority P0 - Authorization regression

- Dua peserta: read/update/delete team milik peserta lain.
- Dua peserta: pelunasan registration milik peserta lain.
- Dua organizer: read/write assessment, rekap, publish, event level update/delete lintas organizer.
- Staff organizer: permission terhadap resource parent dan larangan terhadap organizer lain.
- Banned user: login baru dan token lama.

### Priority P0 - Financial and transaction integrity

- Dua approval team berjalan paralel dengan saldo hanya cukup untuk satu approval.
- Dua approval top-up untuk transaction yang sama.
- Rollback event jika wallet creation gagal.
- Rollback finalize assessment jika salah satu insert gagal.
- Idempotency approval/reject/finalize.

### Priority P1 - PostgreSQL integration

- Fresh schema pada database disposable.
- Migration dari schema sebelumnya.
- Semua enum termasuk `KICKED`.
- Foreign key dan delete behavior.
- `FOR UPDATE` benar-benar muncul pada SQL.
- RowsAffected untuk target yang tidak ada.

### Priority P1 - Middleware and dependency failure

- Missing/invalid API key.
- Missing/malformed/expired/wrong-signature JWT.
- Non-Bearer authorization header.
- Redis unavailable saat startup, authentication, dan logout.
- CORS preflight.
- Rate limiting per IP dan perilaku di belakang reverse proxy.

### Priority P1 - Upload

- Filename dengan `../`, slash, Unicode, dan nama sangat panjang.
- MIME mismatch dan executable content.
- Oversized file.
- Unauthorized direct access ke ID card dan payment proof.
- Cleanup file saat database transaction gagal.

### Priority P2 - Contract and validation

- Invalid/empty DTO untuk seluruh endpoint.
- Invalid event date dan urutan tanggal.
- Invalid enum/status/payment type.
- `CoinRate = 0`.
- Score pada gap/overlap grade range.
- Swagger route and schema contract test.
- Konsistensi HTTP status untuk `400`, `401`, `403`, `404`, `409`, dan `500`.

## 8. Release Gate Recommendation

Minimum gate sebelum production:

1. Selesaikan QA-001 sampai QA-013.
2. Tambahkan integration test PostgreSQL untuk authorization dan transaksi finansial.
3. Tambahkan middleware test untuk banned user dan Redis failure.
4. Pisahkan private upload dari `/public`.
5. Pastikan migration fail fast dan enum schema sinkron.
6. Jalankan CI dengan Go 1.24.3, `go test -race ./...`, `go vet ./...`, build, dan Swagger drift check.
7. Naikkan coverage pada business-critical service/repository; angka target perlu disepakati, tetapi 3.0% belum memadai sebagai release signal.

## 9. Final Assessment

Project memiliki struktur modular yang cukup jelas, build sehat, dan beberapa test terarah sudah ada. Kelemahan utama bukan pada kemampuan compile, tetapi pada enforcement authorization per object, konsistensi transaksi, keamanan file, startup dependency handling, dan rendahnya regression coverage.

Rekomendasi QA saat ini adalah **hold production release** sampai temuan Critical dan High divalidasi serta diperbaiki di branch terpisah dengan automated regression tests.
