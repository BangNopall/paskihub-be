# Design Spec: Ranking System (Event Award) Backend

**Date:** 2026-05-04  
**Topic:** Ranking System Configuration for Organizer Dashboard  
**Status:** Draft (Pending User Review)

---

## 1. Overview
Halaman **Ranking System** pada Organizer Dashboard membutuhkan API untuk mengelola konfigurasi juara (`EventAward`). Konfigurasi ini menentukan bagaimana pemenang dikalkulasi berdasarkan akumulasi nilai dari beberapa `ScoreCategory`.

Satu juara (misal: "Juara Umum") dapat mengambil nilai dari satu atau lebih kategori penilaian. Satu kategori penilaian juga dapat digunakan di beberapa jenis juara. Oleh karena itu, relasi yang digunakan adalah **Many-to-Many**.

## 2. Architecture & Data Model

### 2.1 Entity (`domain/entity/assessment.go`)
Kita akan menambahkan entitas `EventAward` dan mendefinisikan relasi many-to-many ke `ScoreCategory`.

```go
type EventAward struct {
    ID              uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
    EventID         uuid.UUID       `gorm:"type:uuid;not null"`
    EventLevelID    uuid.UUID       `gorm:"type:uuid;not null"`
    Name            string          `gorm:"type:varchar(255);not null"`
    LimitRank       int             `gorm:"not null;default:1"`
    CreatedAt       time.Time
    UpdatedAt       time.Time

    // Relation
    Event           Event           `gorm:"foreignKey:EventID;references:Id"`
    EventLevel      EventLevel      `gorm:"foreignKey:EventLevelID;references:Id"`
    ScoreCategories []ScoreCategory `gorm:"many2many:event_award_score_categories;"`
}
```
*GORM akan secara otomatis membuat tabel join `event_award_score_categories`.*

### 2.2 DTO (`domain/dto/assessment_dto.go`)

```go
type CreateAwardReq struct {
    EventLevelID     uuid.UUID   `json:"event_level_id" validate:"required"`
    Name             string      `json:"name" validate:"required"`
    LimitRank        int         `json:"limit_rank" validate:"required,min=1"`
    ScoreCategoryIDs []uuid.UUID `json:"score_category_ids" validate:"required,min=1"`
}

type UpdateAwardReq struct {
    Name             string      `json:"name" validate:"required"`
    LimitRank        int         `json:"limit_rank" validate:"required,min=1"`
    ScoreCategoryIDs []uuid.UUID `json:"score_category_ids" validate:"required,min=1"`
}

type AwardRes struct {
    ID               uuid.UUID          `json:"id"`
    EventID          uuid.UUID          `json:"event_id"`
    EventLevelID     uuid.UUID          `json:"event_level_id"`
    Name             string             `json:"name"`
    LimitRank        int                `json:"limit_rank"`
    ScoreCategories  []ScoreCategoryRes `json:"score_categories,omitempty"`
}
```

## 3. API Endpoints

Semua endpoint dilindungi oleh middleware `Authentication`, `ApiKey`, dan `AuthOrganizer`.

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/api/v1/eo/events/:eventId/awards` | List awards (optional query: `level_id`) |
| `POST` | `/api/v1/eo/events/:eventId/awards` | Create new award configuration |
| `PUT` | `/api/v1/eo/events/:eventId/awards/:id` | Update award configuration |
| `DELETE` | `/api/v1/eo/events/:eventId/awards/:id` | Delete award configuration |

## 4. Component Logic

### 4.1 Repository (`internal/app/assessment/repository/assessment_repository.go`)
- `CreateAward`: Simpan entitas dan asosiasi many-to-many.
- `GetAwardsByEvent`: Fetch awards dengan `Preload("ScoreCategories")`.
- `UpdateAward`: Update basic info dan ganti asosiasi (`FullSave` atau `Association.Replace`).
- `DeleteAward`: Hapus entitas (asosiasi di tabel join akan terhapus otomatis oleh GORM jika diatur cascade).

### 4.2 Service (`internal/app/assessment/service/assessment_service.go`)
- Menangani pengecekan kepemilikan event (`ensureOwnership`).
- Validasi bahwa `EventLevelID` dan `ScoreCategoryIDs` valid dan milik event yang sama.
- Konversi antara Entity dan DTO.

### 4.3 Controller (`internal/app/assessment/controller/assessment_controller.go`)
- Routing baru di bawah group `/api/v1/eo/events/:eventId/assessment/awards`.
- Validasi input menggunakan `s.validator`.

## 5. Implementation Plan
1. Update `domain/entity/assessment.go` dengan `EventAward`.
2. Tambahkan DTO di `domain/dto/assessment_dto.go`.
3. Update Contract di `domain/contracts/assessment_contract.go`.
4. Implementasi logic di `assessment_repository.go` dan `assessment_service.go`.
5. Tambahkan handler di `assessment_controller.go` dan daftarkan route.
6. Verifikasi manual/test.

---
**Review Check:**
- [x] No TBD/Placeholders.
- [x] Internal consistency maintained.
- [x] Focused scope.
- [x] Explicit requirements.
