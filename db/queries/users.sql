-- name: CreateUser :one
INSERT INTO users (
    email,
    password_hash,
    username,
    display_name
)
VALUES ($1, $2, $3, $4)
RETURNING
    id,
    email,
    password_hash,
    username,
    display_name,
    avatar_url,
    created_at,
    updated_at;

-- name: GetUserByID :one
SELECT
    id,
    email,
    password_hash,
    username,
    display_name,
    avatar_url,
    created_at,
    updated_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT
    id,
    email,
    password_hash,
    username,
    display_name,
    avatar_url,
    created_at,
    updated_at
FROM users
WHERE email = $1;