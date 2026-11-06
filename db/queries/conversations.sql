-- name: CreateConversation :one
INSERT INTO conversations (
    type,
    name,
    created_by
)
VALUES ($1, $2, $3)
RETURNING
    id,
    type,
    name,
    created_by,
    created_at,
    updated_at;

-- name: AddConversationMember :one
INSERT INTO conversation_members (
    conversation_id,
    user_id,
    role
)
VALUES ($1, $2, $3)
RETURNING
    conversation_id,
    user_id,
    role,
    joined_at,
    last_read_message_id;

-- name: GetConversationByID :one
SELECT
    id,
    type,
    name,
    created_by,
    created_at,
    updated_at
FROM conversations
WHERE id = $1;

-- name: ListConversationMembers :many
SELECT
    conversation_id,
    user_id,
    role,
    joined_at,
    last_read_message_id
FROM conversation_members
WHERE conversation_id = $1
ORDER BY joined_at ASC;

-- name: IsConversationMember :one
SELECT EXISTS (
    SELECT 1
    FROM conversation_members
    WHERE conversation_id = $1
      AND user_id = $2
);

-- name: ListConversationMemberIDs :many
SELECT user_id
FROM conversation_members
WHERE conversation_id = $1;