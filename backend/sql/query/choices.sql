-- name: ListChoices :many
SELECT p.action, c.path_id FROM choices c
JOIN page p ON c.path_id = p.id
WHERE page_id = $1;

-- name: CreateChoices :one
INSERT INTO choices (page_id, path_id)
VALUES ($1, $2)
RETURNING *;