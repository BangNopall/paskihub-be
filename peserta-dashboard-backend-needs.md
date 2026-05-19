# Peserta Dashboard Backend Needs

## Context

The participant dashboard page at `src/app/peserta/dashboard/page.tsx` is integrated with:

```http
GET /api/v1/peserta/dashboard
```

The current Swagger contract supports dashboard statistics and recent activities:

- `data.stats.total_team`
- `data.stats.active_event`
- `data.stats.finished_event`
- `data.stats.pending_payment`
- `data.recent_activities[].title`
- `data.recent_activities[].description`
- `data.recent_activities[].time`

The frontend can render those fields with SSR.

## Missing Fields

The current UI has an `Event Mendatang` section, but `dto.ParticipantDashboardRes` does not expose upcoming event data yet.

Please add an `upcoming_events` array to `GET /api/v1/peserta/dashboard`.

Recommended response shape:

```json
{
  "data": {
    "stats": {
      "total_team": 2,
      "active_event": 1,
      "finished_event": 0,
      "pending_payment": 1
    },
    "recent_activities": [
      {
        "title": "Lomba Paskibra Nasional 2026",
        "description": "Status pendaftaran tim Garuda Muda: WAITING",
        "time": "5 jam yang lalu"
      }
    ],
    "upcoming_events": [
      {
        "id": "event-or-registration-id",
        "title": "Lomba Paskibra Nasional 2026",
        "date": "2026-03-15",
        "registered_teams": 18,
        "status": "OPEN",
        "detail_url_id": "registration-id-or-event-id"
      }
    ]
  }
}
```

## Field Notes

- `id`: Stable identifier for the event item.
- `title`: Event name shown in the card.
- `date`: Competition date or closest relevant event date.
- `registered_teams`: Number of registered teams for the event.
- `status`: Event status label source, for example `OPEN`, `CLOSED`, or `FINISHED`.
- `detail_url_id`: Identifier the frontend can use for detail navigation. If participant dashboard should link to `/peserta/dashboard/event/[id]/overview`, this should be the registration ID. If it should link to a public event detail page later, this can be the event ID.

## Existing Stat Accuracy

`stats.finished_event` exists in Swagger, but the backend repository currently appears to return `0` with a `Need logic for finished` comment. The frontend can display this field, but the backend should calculate it from completed participant registrations/events when the rule is finalized.
