# High Risk Integrity and Security Batch 2 Design

## Goal

Fix QA-005 through QA-013 while preserving existing endpoint names, successful response shapes, and JSON field names.

## Architecture

The implementation follows the existing controller-service-repository direction. Authorization and business invariants remain in services, relation and locking predicates live in repositories, and startup dependencies fail before the HTTP server begins accepting traffic.

## Relational Authorization

### Event Levels

Event-level updates and deletes use both the level ID and parent event ID. The event ID supplied in the request body is not trusted for persistence. A missing level under the requested event returns `domain.ErrNotFound`; an event owned by another organizer remains forbidden through the existing event ownership check.

### Assessment CRUD

The assessment repository exposes focused relation checks:

- Event level belongs to the event.
- Score category belongs to the event and its event level.
- Registration, judge, and score subcategory belong to the same event and event level.
- Every event level and score category referenced by an award belongs to the event.

Malformed cross-event relation graphs return `domain.ErrBadRequest`. Existing object update/delete lookups that target an object outside the URL event continue returning `domain.ErrNotFound`.

## Banned Accounts and Sessions

Login rejects banned users with a stable `domain.ErrAccountBanned` mapped to HTTP 403. Authentication validates the JWT and revocation status, then loads the current user record and rejects banned users. This immediately invalidates previously issued tokens without changing JWT claims or frontend token storage.

## Private Files

The broad `/public` static route is removed. Only event logos/posters and team logos remain publicly static under explicit routes.

New sensitive uploads are stored under `storage/private`:

- Team recommendation letters.
- Member ID cards and photos.
- Registration payment proofs.
- Wallet top-up proofs.

Authenticated file endpoints resolve files by resource ID and file type, never by an arbitrary filesystem path. Services enforce:

- Participants may read sensitive files belonging to their institution or registration.
- Organizers may read files for registrations in events they own.
- Admins may read wallet top-up proofs.

Existing JSON field names remain unchanged, but sensitive path values are mapped to authenticated API URLs. Legacy database paths under `public/uploads` remain readable only through the authenticated resolver after the broad static route is removed.

## Wallet Concurrency

Team approval is executed by a transaction-aware repository method. It locks the registration and wallet rows with PostgreSQL `FOR UPDATE`, verifies the registration is still `WAITING`, checks the balance, debits once, creates the withdrawal transaction, and updates payment status atomically.

Top-up approval and rejection replace `gorm:query_option` with `clause.Locking{Strength: "UPDATE"}`. PostgreSQL integration tests run only when `TEST_DATABASE_URL` is set and use isolated rows in the configured test database.

## Schema and Startup Reliability

Redis is mandatory. Its constructor returns `(RedisInterface, error)`, and route mounting returns an error so startup fails before listening when Redis cannot be reached.

Database migration returns errors for table drops, enum setup, `AutoMigrate`, default setting count, and default setting creation. The registration enum migration safely executes:

```sql
ALTER TYPE registration_status ADD VALUE IF NOT EXISTS 'KICKED';
```

The application exits before route mounting when migration fails.

## Error Handling

- Authorization failure: HTTP 403.
- Resource absent under the requested parent: HTTP 404.
- Cross-event or invalid relation graph: HTTP 400.
- Mandatory dependency or migration failure: startup failure.

## Testing

Unit tests cover event-level predicates, assessment relation rejection, banned login/token rejection, Redis constructor failure, private-file authorization, and migration error propagation.

PostgreSQL integration tests cover top-up locking, concurrent team approval, enum `KICKED`, and default migration smoke behavior. Integration tests skip with an explicit message when `TEST_DATABASE_URL` is absent.

