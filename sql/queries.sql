-- name: GetStory :one
SELECT * FROM story
WHERE id = $1 LIMIT 1;

-- name: ListStories :many
SELECT * FROM story
ORDER BY title;

-- name: CreateStory :one
INSERT INTO story (title)
VALUES ($1)
RETURNING *;

-- name: DeleteStory :exec
DELETE FROM story
WHERE id = $1;
