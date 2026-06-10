# Critical Authorization Batch 1 Design

## Goal

Fix QA-001 through QA-004 without changing public response payloads or addressing unrelated QA findings.

## Architecture

Authorization remains a service-layer decision. Controllers pass the authenticated user ID from Fiber locals, services return stable domain errors, and repositories expose narrowly scoped ownership or relation checks backed by database predicates.

## Changes

### Participant Teams

`GetTeamDetail` and `DeleteTeam` resolve the authenticated participant's institution and reject teams whose `InstiId` differs with `domain.ErrForbidden`.

### Participant Payment Completion

`PelunasanEvent` receives the authenticated participant ID. It verifies registration ownership before reading or updating the registration and before saving a new payment proof.

### Assessment Writes

Bulk scores, bulk violations, and finalize receive the authenticated organizer ID. The repository validates that:

- The registration belongs to an event owned by the organizer.
- The judge belongs to that event.
- Every score subcategory belongs to that event.
- Every violation type belongs to that event.

Authorization validation occurs before any write. Invalid ownership or mixed-event references return `domain.ErrForbidden`.

### Organizer Recaps

All organizer recap reads and publish operations receive the authenticated organizer ID. Registration and event-level ownership are resolved by repository queries before data is read or mutated.

## Error Handling

Authorization failures use `domain.ErrForbidden`, mapped to HTTP 403 by `domain.GetCode`. Existing success payloads and route paths remain unchanged.

## Testing

Regression tests cover denied access for:

- Participant team detail and deletion.
- Participant payment completion.
- Organizer assessment write operations.
- Organizer recap read and publish operations.

Tests also assert that denied operations do not invoke write methods.
