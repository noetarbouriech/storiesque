-- name: ListChoices :many
SELECT path_id FROM choices
WHERE page_id = $1;

-- name: CreateChoices :one
INSERT INTO choices (page_id, path_id)
VALUES ($1, $2)
RETURNING *;

