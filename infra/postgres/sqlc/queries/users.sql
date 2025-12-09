-- name: GetUser :one
SELECT id, name, created_at, updated_at
FROM users
WHERE id = $1;

-- name: UpsertUser :exec
INSERT INTO users (id, name, created_at, updated_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    updated_at = EXCLUDED.updated_at;
