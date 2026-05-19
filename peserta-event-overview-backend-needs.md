# Peserta Event Overview Backend Needs

Dokumen ini merangkum gap backend untuk halaman peserta:

`/peserta/dashboard/event/[id]/overview`

Analisis berdasarkan Swagger backend:

`/Users/noxval/_PROJECT_/paskihub-be/docs/swagger.json`

## Endpoint yang Sudah Tersedia

- `GET /api/v1/peserta/events/registrations/{regis_id}`
- `GET /api/v1/peserta/assessment/recap/{regis_id}`
- `GET /api/v1/peserta/rekap/scoreboard/{event_level_id}`
- `PUT /api/v1/peserta/events/register/{id}/pelunasan`

## Gap yang Masih Dibutuhkan

### 1. Tambahkan `event_level_id` ke Detail Registrasi

Halaman overview perlu menampilkan leaderboard peserta melalui endpoint:

`GET /api/v1/peserta/rekap/scoreboard/{event_level_id}`

Namun response `dto.RegistrationDetailResponse` saat ini belum menyediakan `event_level_id`, sehingga frontend tidak punya sumber SSR yang andal untuk memanggil scoreboard.

Tambahkan field berikut di root response detail registrasi:

```json
{
  "event_level_id": "uuid"
}
```

### 2. Tambahkan Logo Event ke Detail Registrasi

UI overview menampilkan logo event di header dan modal pelunasan. Response `dto.RegistrationDetailResponse.event` saat ini hanya memuat:

- `id`
- `title`
- `description`
- `date`
- `location`
- `price`
- `target_date`

Tambahkan salah satu field berikut ke object `event`:

```json
{
  "event": {
    "logo_url": "/uploads/events/logo.png"
  }
}
```

Atau gunakan nama yang konsisten dengan endpoint event lain:

```json
{
  "event": {
    "logo_path": "/uploads/events/logo.png"
  }
}
```

Frontend sudah siap menerima keduanya.

### 3. Opsional: Tambahkan Nilai Maksimal Rekap

Response `dto.AssessmentRecapResponse` sudah menyediakan:

- `final_score`
- `total_score`
- `total_violation_points`
- `categories`
- `violations`

Progress bar di frontend butuh denominator nilai maksimal. Untuk sementara frontend memakai fallback `500`.

Jika nilai maksimal bersifat dinamis per event/level, tambahkan:

```json
{
  "max_score": 500
}
```

ke `dto.AssessmentRecapResponse`.

## Response Ideal Detail Registrasi

```json
{
  "event_level_id": "uuid",
  "event": {
    "id": "uuid",
    "title": "Nama Event",
    "description": "Deskripsi event",
    "date": "2026-05-19",
    "location": "Jakarta",
    "price": "250000",
    "target_date": "2026-06-01T00:00:00Z",
    "logo_url": "/uploads/events/logo.png"
  },
  "team": {
    "id": "uuid",
    "name": "Nama Tim",
    "logo_url": "/uploads/teams/logo.png",
    "official_count": 2,
    "pasukan_count": 20
  },
  "payment": {
    "status": "DP_PAID",
    "amount_paid": "100000",
    "total_amount": "250000",
    "remaining_amount": "150000",
    "proof_url": "/uploads/payments/proof.png"
  }
}
```
