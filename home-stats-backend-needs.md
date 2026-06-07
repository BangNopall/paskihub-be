# Home Stats Backend Needs

## Context

The public home page at `src/app/(home)/page.tsx` renders four aggregate stats with SSR:

- `Jumlah event`
- `Event Organizer`
- `Peserta`
- `Tim`

The current Swagger file does not expose a public endpoint that matches these needs. Existing dashboard endpoints contain partial stats, but they require `BearerAuth` and are scoped to admin or organizer dashboards.

## Required Endpoint

```http
GET /api/v1/public/home-stats
```

Recommended security:

- `ApiKeyAuth`
- No `BearerAuth`

This keeps the endpoint usable from a Next.js Server Component while keeping `API_KEY` server-side.

## Required Response Shape

```json
{
  "data": {
    "total_events": 120,
    "total_organizers": 45,
    "total_participants": 980,
    "total_teams": 210
  }
}
```

## Field Notes

- `total_events`: Total non-archived events in the platform.
- `total_organizers`: Total active organizer accounts.
- `total_participants`: Total active participant accounts.
- `total_teams`: Total non-archived participant teams in the platform.

All fields should be non-negative integers. If backend filtering rules differ, please document them in Swagger so the frontend copy can be aligned with the displayed values.

## Current Frontend Fallback

Until this endpoint exists, the frontend service returns:

```json
{
  "total_events": 0,
  "total_organizers": 0,
  "total_participants": 0,
  "total_teams": 0
}
```

## Swagger Status

As of the Swagger file checked from `/Users/noxval/_PROJECT_/paskihub-be/docs/swagger.json`, no endpoint currently provides these four public aggregate values.
