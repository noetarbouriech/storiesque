-- name: GetStory :one
SELECT * FROM story
WHERE id = $1 LIMIT 1;

-- name: ListStories :many
SELECT * FROM story
ORDER BY title;

-- name: CreateStory :one
INSERT INTO story (title, description)
VALUES ($1, $2)
RETURNING *;

-- name: DeleteStory :exec
DELETE FROM story
WHERE id = $1;



-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserWithEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserWithUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (username, password_hash, email)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY username;

-- name: UpdateUser :exec
UPDATE users
SET
  username = CASE WHEN @username_do_update::boolean
    THEN @username::VARCHAR(24) ELSE username END,

  password_hash = CASE WHEN @password_hash_do_update::boolean
    THEN @password_hash::VARCHAR(64) ELSE password_hash END,

  email = CASE WHEN @email_do_update::boolean
    THEN @email::VARCHAR(128) ELSE email END
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

