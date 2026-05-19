# Prompt Gemini CLI — Dokumentasi Swagger Paskihub Backend

Gunakan prompt berikut di Gemini CLI untuk membuat, memperbaiki, dan melengkapi dokumentasi Swagger/OpenAPI pada project backend **Paskihub**.

```md
Kamu adalah senior backend engineer, Go Fiber specialist, dan API documentation specialist.

Saya sedang mengerjakan project backend bernama Paskihub atau paskihub-be. Project ini adalah backend application berbasis Go untuk mengelola event Paskibra, participants, assessments, wallets, dan fitur terkait lainnya.

Tech stack project:
- Language: Go / Golang
- Framework: Fiber v2
- ORM: GORM
- Database: PostgreSQL
- Cache: Redis
- API Documentation: Swagger menggunakan swaggo/swag
- Config Management: Viper
- Development Tool: Air
- Validation: github.com/go-playground/validator/v10
- Authentication: JWT dengan middleware Auth
- Architecture: Modular feature-based architecture dengan pendekatan Clean Architecture

Struktur project penting:
- cmd/app/main.go
  Entry point aplikasi. Digunakan untuk inisialisasi server, database, dan route.

- internal/infra/server/http_server.go
  Lokasi penting untuk route definitions, server configuration, middleware, validator registration, dan mounting routes.

- internal/app/
  Berisi feature module seperti user, event, assessment, wallet, participant, auth, dan module lainnya.
  Setiap module biasanya memiliki:
  - controller/
  - service/
  - repository/

- domain/
  Berisi core business definition.
  Subfolder penting:
  - domain/contracts/
  - domain/entity/
  - domain/dto/
  - domain/enums/

- internal/middlewares/
  Berisi middleware seperti Auth, CORS, Rate Limiter, dan middleware lainnya.

- pkg/
  Berisi utility package seperti JWT, bcrypt, UUID, logging, dan helper lainnya.

- docs/
  Berisi auto-generated Swagger documentation dari swaggo/swag.

Tugas utama:
Tolong audit, buat, perbaiki, dan lengkapi dokumentasi Swagger/OpenAPI untuk seluruh endpoint API pada project Paskihub agar dokumentasinya lengkap, akurat, konsisten, dan mudah digunakan oleh frontend developer saat integrasi API.

Tujuan dokumentasi:
Frontend developer harus bisa memahami setiap endpoint tanpa menebak-nebak:
- Endpoint apa yang harus dipanggil
- Method apa yang digunakan
- Auth token diperlukan atau tidak
- Role/permission apa yang dibutuhkan
- Path parameter apa yang diperlukan
- Query parameter apa saja yang tersedia
- Request body apa yang harus dikirim
- Field mana yang required dan optional
- Format data setiap field
- Response sukses yang diterima
- Response error yang mungkin terjadi
- Struktur DTO request dan response
- Struktur pagination, filter, sorting, search, dan error validation

Instruksi kerja:

## 1. Scan dan pahami project

- Baca struktur folder project secara menyeluruh.
- Mulai dari cmd/app/main.go untuk memahami proses bootstrap aplikasi.
- Baca internal/infra/server/http_server.go untuk memahami konfigurasi server, route, middleware, validator, dan Swagger setup.
- Telusuri semua route yang dimount.
- Cari seluruh controller di internal/app/*/controller.
- Cari seluruh DTO di domain/dto.
- Cari seluruh entity di domain/entity.
- Cari seluruh enum di domain/enums.
- Cari middleware auth di internal/middlewares.
- Cari helper response/error jika ada.
- Jangan hanya menebak dari nama route. Pastikan dokumentasi sesuai implementasi aktual.

## 2. Identifikasi Swagger/swaggo setup

- Pastikan project menggunakan swaggo/swag.
- Periksa apakah sudah ada Swagger annotation di controller.
- Periksa file docs/ yang dihasilkan.
- Periksa konfigurasi Swagger UI, terutama path /swagger/.
- Periksa apakah main.go sudah memiliki general API info annotation.
- Jika belum lengkap, tambahkan atau perbaiki annotation global seperti:
  - @title
  - @version
  - @description
  - @host jika relevan
  - @BasePath
  - @schemes
  - @securityDefinitions.apikey BearerAuth
  - @in header
  - @name Authorization
  - @description JWT Authorization header using the Bearer scheme. Example: "Bearer {token}"

## 3. Dokumentasikan seluruh endpoint

Untuk setiap endpoint pada controller, tambahkan atau perbaiki swaggo comment dengan format lengkap.

Setiap endpoint minimal wajib memiliki:
- @Summary
- @Description
- @Tags
- @Accept json
- @Produce json
- @Param untuk path parameter, query parameter, header, dan body
- @Success untuk response sukses
- @Failure untuk response error
- @Security BearerAuth jika endpoint butuh JWT
- @Router dengan path dan method yang benar

Contoh format yang diharapkan:

```go
// CreateEvent godoc
// @Summary Create new event
// @Description Create a new Paskibra event. This endpoint requires authentication.
// @Tags Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateEventRequest true "Create event request body"
// @Success 201 {object} dto.SuccessResponse{data=dto.EventResponse} "Event created successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request payload"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Forbidden"
// @Failure 422 {object} dto.ValidationErrorResponse "Validation error"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /events [post]
```

Sesuaikan contoh di atas dengan implementasi aktual project.

## 4. Dokumentasikan request DTO

Untuk setiap request DTO di domain/dto, pastikan:
- Field memiliki json tag yang benar.
- Field memiliki validate tag jika memang digunakan.
- Field memiliki example tag jika memungkinkan.
- Nama field jelas.
- Required/optional sesuai validasi aktual.
- Format khusus terdokumentasi, misalnya:
  - email
  - uuid
  - date
  - datetime
  - enum
  - number
  - boolean
  - array
  - object
- Enum value ditulis jelas.
- Field yang nullable dijelaskan jika perlu.

Contoh:

```go
type LoginRequest struct {
    Email    string `json:"email" validate:"required,email" example:"admin@paskihub.id"`
    Password string `json:"password" validate:"required,min=8" example:"password123"`
}
```

## 5. Dokumentasikan response DTO

Untuk setiap response yang dikembalikan API:
- Pastikan schema response tersedia.
- Dokumentasikan field success/status, message, data, meta, errors, dan lainnya sesuai implementasi aktual.
- Jangan membuat response palsu yang tidak sesuai kode.
- Jika project memakai response wrapper, buat schema reusable.
- Jika response mengandung nested object, dokumentasikan nested field.
- Jika response berupa array, jelaskan item array.
- Jika response berupa pagination, dokumentasikan meta pagination.

Contoh reusable response:

```go
type SuccessResponse struct {
    Success bool        `json:"success" example:"true"`
    Message string      `json:"message" example:"Request processed successfully"`
    Data    interface{} `json:"data"`
}

type ErrorResponse struct {
    Success bool   `json:"success" example:"false"`
    Message string `json:"message" example:"Something went wrong"`
}

type ValidationErrorResponse struct {
    Success bool                   `json:"success" example:"false"`
    Message string                 `json:"message" example:"Validation failed"`
    Errors  map[string]interface{} `json:"errors"`
}

type PaginationMeta struct {
    Page      int   `json:"page" example:"1"`
    Limit     int   `json:"limit" example:"10"`
    Total     int64 `json:"total" example:"100"`
    TotalPage int   `json:"total_page" example:"10"`
}
```

Namun gunakan struktur aktual dari project. Jangan memaksakan struktur baru jika helper response project berbeda.

## 6. Dokumentasikan error response

Untuk setiap endpoint, tambahkan kemungkinan error sesuai implementasi aktual:
- 400 Bad Request
- 401 Unauthorized
- 403 Forbidden
- 404 Not Found
- 409 Conflict
- 422 Validation Error
- 500 Internal Server Error

Untuk setiap error:
- Gunakan DTO error response yang sesuai.
- Tambahkan description yang jelas.
- Pastikan status code sesuai logic di controller/service.
- Jika error berasal dari domain/errors.go, dokumentasikan berdasarkan error tersebut.
- Jika ada validasi menggunakan go-playground/validator/v10, dokumentasikan validation error response.

## 7. Dokumentasikan authentication dan authorization

Project ini menggunakan JWT dan Auth middleware.
Tolong:
- Cari middleware Auth di internal/middlewares.
- Identifikasi format token yang digunakan.
- Tambahkan Swagger security scheme BearerAuth jika belum ada.
- Tambahkan @Security BearerAuth di semua endpoint protected.
- Jangan tambahkan @Security BearerAuth pada endpoint public seperti login/register jika memang public.
- Jika ada role-based access control di service atau middleware, dokumentasikan role yang dibutuhkan di @Description.
- Jika role tidak eksplisit di route, telusuri logic service sebelum menyimpulkan.

## 8. Dokumentasikan query params untuk list endpoint

Untuk endpoint list data seperti users, events, participants, assessments, wallets, atau endpoint lain:
- Cari DTO filter/query jika ada.
- Dokumentasikan semua query params:
  - page
  - limit
  - search
  - sort
  - order
  - status
  - date
  - event_id
  - participant_id
  - category
  - filter lain yang tersedia
- Jelaskan default value jika ada.
- Jelaskan format value jika berupa enum/date/uuid.
- Dokumentasikan response pagination jika endpoint memakai pagination.
- Tambahkan example request URL di description jika berguna.

Contoh:

```go
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Param search query string false "Search keyword"
// @Param sort query string false "Sort field" Enums(created_at, updated_at, name)
// @Param order query string false "Sort order" Enums(asc, desc)
```

## 9. Dokumentasikan path params

Untuk endpoint dengan ID:
- Pastikan @Param id path string true jelas.
- Jika ID menggunakan UUID, tulis format UUID.
- Jika ID integer, tulis integer.
- Sesuaikan dengan route aktual.

Contoh:

```go
// @Param id path string true "Event ID in UUID format"
```

## 10. Dokumentasikan endpoint berdasarkan modul

Gunakan tag yang konsisten berdasarkan module, misalnya:
- Auth
- Users
- Events
- Participants
- Assessments
- Wallets
- Admin
- Files
- Payments
- Health Check

Sesuaikan dengan module aktual yang ada di internal/app.

## 11. Perbaiki atau tambahkan DTO khusus dokumentasi jika diperlukan

Jika response project menggunakan anonymous struct, map, atau interface{} sehingga Swagger sulit membaca schema:
- Buat DTO response eksplisit di domain/dto.
- Gunakan DTO tersebut untuk dokumentasi Swagger.
- Jangan mengubah response runtime kecuali benar-benar aman.
- DTO dokumentasi harus merepresentasikan response aktual.

## 12. Jangan ubah logic bisnis

Kamu boleh:
- Menambahkan Swagger comment.
- Menambahkan DTO/schema untuk dokumentasi.
- Menambahkan example tag.
- Menambahkan description pada DTO.
- Merapikan nama schema documentation.
- Memperbaiki Swagger config.
- Menjalankan swag init.

Kamu tidak boleh:
- Mengubah logic bisnis utama.
- Mengubah database schema.
- Mengubah struktur response API existing secara breaking change.
- Mengubah route existing.
- Menghapus kode yang masih dipakai.
- Mengubah behavior auth/authorization.
- Mengubah service/repository logic tanpa alasan kuat.

## 13. Jika menemukan inkonsistensi

Jika route, DTO, response, atau validasi tidak konsisten:
- Dokumentasikan sesuai implementasi aktual.
- Tambahkan catatan TODO di tempat yang aman jika diperlukan.
- Berikan rekomendasi pada final report.
- Jangan membuat dokumentasi yang lebih ideal tetapi tidak sesuai kode.

## 14. Generate Swagger docs

Setelah memperbaiki annotation:
- Jalankan swag init dari root project.
- Jika command berbeda, cek README/go.mod/Makefile.
- Perbaiki semua error parsing dari swaggo.
- Pastikan docs/ berhasil diperbarui.
- Pastikan build tetap aman.

Command yang bisa dicoba:

```bash
swag init -g cmd/app/main.go
```

Jika gagal karena struktur project berbeda, cari entry point yang benar lalu jalankan command yang sesuai.

## 15. Validasi hasil akhir

Setelah perubahan:
- Jalankan gofmt pada file yang diubah.
- Jalankan go test jika memungkinkan.
- Jalankan go build jika memungkinkan.
- Pastikan tidak ada error compile akibat DTO atau annotation.
- Pastikan Swagger UI tetap bisa dibuka di /swagger/ pada mode development.

Command yang bisa dicoba:

```bash
gofmt -w <files_changed>
go test ./...
go build -o ./tmp/main ./cmd/app/main.go
```

## 16. Output final report

Setelah selesai, berikan laporan akhir berisi:

### A. Project understanding
- Framework yang terdeteksi
- Swagger tool yang digunakan
- Lokasi route utama
- Lokasi DTO
- Lokasi controller
- Lokasi middleware auth

### B. Files changed
- Daftar file yang diubah
- Ringkasan perubahan tiap file

### C. Endpoints documented
Buat daftar endpoint yang berhasil didokumentasikan dalam format:
- METHOD /path - tag - auth/public - status

### D. Schemas added or improved
- Request DTO
- Response DTO
- Error DTO
- Pagination DTO
- Auth DTO
- Entity response DTO
- Enum schema

### E. Validation result
- swag init status
- gofmt status
- go test status
- go build status

### F. Known ambiguity / TODO
- Endpoint yang response-nya masih ambigu
- Endpoint yang error handling-nya belum konsisten
- DTO yang belum eksplisit
- Response yang masih menggunakan map/interface{}
- Rekomendasi improvement untuk API contract

### G. Frontend integration notes
- Cara menggunakan Bearer token
- Format response umum
- Format error umum
- Format pagination
- Catatan penting untuk integrasi frontend

Important:
- Dokumentasi harus berdasarkan kode aktual, bukan asumsi.
- Prioritaskan akurasi dibanding terlihat lengkap tapi salah.
- Jika ada endpoint yang tidak bisa dipastikan, tulis sebagai ambiguity di final report.
- Jangan membuat field palsu.
- Jangan membuat endpoint palsu.
- Jangan mengubah API contract tanpa instruksi.
- Fokus utama adalah membuat dokumentasi Swagger yang detail, konsisten, dan frontend-friendly.

Mulai sekarang:
1. Scan project structure.
2. Baca cmd/app/main.go.
3. Baca internal/infra/server/http_server.go.
4. Identifikasi semua routes.
5. Telusuri controller, DTO, entity, middleware, dan response helper.
6. Perbaiki Swagger annotation endpoint satu per satu.
7. Tambahkan schema DTO yang dibutuhkan.
8. Jalankan swag init.
9. Validasi dengan gofmt, go test, dan go build.
10. Berikan final report.
```

---

## Prompt Follow-up untuk Review Hasil Gemini

Gunakan prompt ini setelah Gemini selesai membuat atau memperbaiki dokumentasi Swagger:

```md
Sekarang review ulang hasil perubahan dokumentasi Swagger yang sudah kamu buat.

Tolong cek:
1. Apakah semua route di internal/infra/server/http_server.go sudah punya dokumentasi Swagger.
2. Apakah semua controller di internal/app/*/controller sudah punya annotation swaggo.
3. Apakah semua endpoint protected sudah memiliki @Security BearerAuth.
4. Apakah semua request body memakai DTO yang tepat.
5. Apakah semua success response dan error response sesuai implementasi aktual.
6. Apakah semua query params, path params, dan pagination sudah terdokumentasi.
7. Apakah docs/ berhasil digenerate ulang dengan swag init.
8. Apakah tidak ada dokumentasi palsu atau field yang tidak sesuai kode.

Jika ada kekurangan, langsung perbaiki.
Setelah itu berikan checklist endpoint yang sudah documented dan endpoint yang masih perlu review manual.
```
