-- name: GetStory :one
SELECT s.id, s.title, s.description, s.first_page_id, u.username as author_name FROM story s
JOIN "user" u ON s.author = u.id
WHERE s.id = $1 LIMIT 1;

-- name: SearchStories :many
SELECT s.id, s.title, s.description, u.username as author_name FROM story s
JOIN "user" u ON s.author = u.id
WHERE title LIKE '%' || @title || '%'
ORDER BY title;

-- name: CreateStory :one
INSERT INTO story (title, description, author)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteStory :exec
DELETE FROM story
WHERE id = $1;
