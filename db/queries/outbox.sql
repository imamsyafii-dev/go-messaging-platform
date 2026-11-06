-- name: CreateOutboxEvent :one
INSERT INTO outbox_events (
    aggregate_type,
    aggregate_id,
    event_type,
    payload
)
VALUES ($1, $2, $3, $4)
RETURNING
    id,
    aggregate_type,
    aggregate_id,
    event_type,
    payload,
    created_at,
    published_at,
    attempts;

-- name: ListPendingOutboxEvents :many
SELECT
    id,
    aggregate_type,
    aggregate_id,
    event_type,
    payload,
    created_at,
    published_at,
    attempts
FROM outbox_events
WHERE published_at IS NULL
ORDER BY created_at ASC, id ASC
LIMIT $1;

-- name: MarkOutboxEventPublished :exec
UPDATE outbox_events
SET published_at = NOW()
WHERE id = $1;

-- name: IncrementOutboxEventAttempts :exec
UPDATE outbox_events
SET attempts = attempts + 1
WHERE id = $1;