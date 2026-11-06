-- name: CreateMessage :one
INSERT INTO messages (
    conversation_id,
    sender_id,
    client_message_id,
    type,
    content
)
VALUES ($1, $2, $3, $4, $5)
RETURNING
    id,
    conversation_id,
    sender_id,
    client_message_id,
    type,
    content,
    created_at,
    updated_at,
    deleted_at;

-- name: GetMessageByID :one
SELECT
    id,
    conversation_id,
    sender_id,
    client_message_id,
    type,
    content,
    created_at,
    updated_at,
    deleted_at
FROM messages
WHERE id = $1;

-- name: GetMessageByClientID :one
SELECT
    id,
    conversation_id,
    sender_id,
    client_message_id,
    type,
    content,
    created_at,
    updated_at,
    deleted_at
FROM messages
WHERE sender_id = $1
  AND client_message_id = $2;

-- name: ListMessages :many
SELECT
    id,
    conversation_id,
    sender_id,
    client_message_id,
    type,
    content,
    created_at,
    updated_at,
    deleted_at
FROM messages
WHERE conversation_id = $1
ORDER BY created_at DESC, id DESC
LIMIT $2;

-- name: ListMessagesBefore :many
SELECT
    id,
    conversation_id,
    sender_id,
    client_message_id,
    type,
    content,
    created_at,
    updated_at,
    deleted_at
FROM messages
WHERE conversation_id = $1
  AND (
      created_at < $2
      OR (created_at = $2 AND id < $3)
  )
ORDER BY created_at DESC, id DESC
LIMIT $4;