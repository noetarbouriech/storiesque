-- name: GetUser :one
SELECT * FROM "user"
WHERE id = $1 LIMIT 1;

-- name: GetUserWithEmail :one
SELECT * FROM "user"
WHERE email = $1 LIMIT 1;

-- name: GetUserWithUsername :one
SELECT * FROM "user"
WHERE username = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO "user" (username, password_hash, email)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListUsers :many
SELECT * FROM "user"
ORDER BY username;

-- name: UpdateUser :exec
UPDATE "user"
SET
  username = CASE WHEN @username_do_update::boolean
    THEN @username::VARCHAR(24) ELSE username END,

  password_hash = CASE WHEN @password_hash_do_update::boolean
    THEN @password_hash::VARCHAR(64) ELSE password_hash END,

  email = CASE WHEN @email_do_update::boolean
    THEN @email::VARCHAR(128) ELSE email END
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM "user"
WHERE id = $1;
