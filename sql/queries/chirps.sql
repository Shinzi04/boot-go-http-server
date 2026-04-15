-- name: NewChirp :one
INSERT INTO
    chirps (id, created_at, updated_at, body, user_id)
VALUES
    (gen_random_uuid(), NOW(), NOW(), $1, $2)
RETURNING *;

-- name: GetChirps :many
SELECT * FROM chirps ORDER BY created_at ASC;

-- name: GetChirpById :one
SELECT * FROM chirps WHERE id = $1 LIMIT 1;

-- name: GetUserByChirpID :one
SELECT users.*
FROM users INNER JOIN chirps
ON users.id = chirps.user_id
WHERE chirps.id = $1;

-- name: DeleteChirpById :exec
DELETE FROM chirps
WHERE id = $1;

-- name: GetChirpByUserID :many
SELECT * FROM chirps
WHERE user_id = $1
ORDER BY created_at ASC;