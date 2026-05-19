# Admin Dashboard Backend Contract

## Context

The admin dashboard page at `src/app/admin/dashboard/page.tsx` should render with SSR from a dedicated backend endpoint:

```http
GET /api/v1/admin/dashboard
```

The current Swagger file does not expose this endpoint. Existing endpoints can provide some raw data, but they do not provide reliable aggregate values for the dashboard without fetching large lists and calculating totals in the frontend.

## Current Frontend Needs

- `data.stats.total_revenue.value`
- `data.stats.total_revenue.trend`
- `data.stats.total_eo.value`
- `data.stats.total_eo.trend`
- `data.stats.total_participants.value`
- `data.stats.total_participants.trend`
- `data.stats.pending_topups.value`
- `data.stats.pending_topups.trend`
- `data.recent_transactions[].id`
- `data.recent_transactions[].eo_name`
- `data.recent_transactions[].amount`
- `data.recent_transactions[].amount_koin`
- `data.recent_transactions[].time_ago`
- `data.recent_transactions[].status`
- `data.eo_registrations[].id`
- `data.eo_registrations[].name`
- `data.eo_registrations[].email`
- `data.eo_registrations[].registered_at`

## Recommended Response Shape

```json
{
  "data": {
    "stats": {
      "total_revenue": {
        "value": 125000000,
        "trend": "+12% bulan ini"
      },
      "total_eo": {
        "value": 42,
        "trend": "+5 baru"
      },
      "total_participants": {
        "value": 1240,
        "trend": "+156 baru"
      },
      "pending_topups": {
        "value": 8,
        "trend": "Perlu approval"
      }
    },
    "recent_transactions": [
      {
        "id": "transaction-id",
        "eo_name": "SMA 1 Jakarta",
        "amount": 500000,
        "amount_koin": 500,
        "time_ago": "10 menit yang lalu",
        "status": "PENDING"
      }
    ],
    "eo_registrations": [
      {
        "id": "user-id",
        "name": "Lomba Jaya Abadi",
        "email": "contact@lomajaya.com",
        "registered_at": "2026-04-24T00:00:00Z"
      }
    ]
  }
}
```

## Field Notes

- `total_revenue.value`: Total approved top-up revenue in IDR.
- `total_revenue.trend`: Backend-calculated trend label for the current dashboard period.
- `total_eo.value`: Total active organizer accounts.
- `total_eo.trend`: New organizer count or percentage for the current dashboard period.
- `total_participants.value`: Total active participant accounts.
- `total_participants.trend`: New participant count or percentage for the current dashboard period.
- `pending_topups.value`: Total top-up transactions with `PENDING` status.
- `pending_topups.trend`: Short label for pending approval state.
- `recent_transactions`: Latest top-up transactions for the admin approval preview.
- `eo_registrations`: Latest organizer registrations.
- `status`: Should use the existing transaction enum values: `PENDING`, `APPROVE`, `REJECTED`.

## Existing Swagger Gaps

- No `GET /api/v1/admin/dashboard` endpoint exists.
- `GET /api/v1/wallets/admin/transactions` provides recent transactions, but it is paginated list data, not dashboard aggregate data.
- `GET /api/v1/admin/users` can provide organizer and participant lists, but counting dashboard totals from the frontend is not ideal for large datasets.
- `total_revenue`, stat trends, and accurate total pending top-up count are not available as dedicated dashboard aggregate fields.

## Existing Swagger Mismatch To Review

The current frontend transaction schema expects top-level pagination fields:

```json
{
  "transactions": [],
  "total": 10,
  "page": 1,
  "limit": 10
}
```

Swagger currently documents:

```json
{
  "transactions": [],
  "pagination": {
    "page": 1,
    "total_pages": 1
  }
}
```

Please align either the backend response or the frontend schema before relying on transaction pagination totals for admin dashboard statistics.
