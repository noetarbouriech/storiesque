-- name: GetShelf :many
SELECT st.id, st.title, st.description, st.has_img, u.username as author_name FROM shelf sh
JOIN story st ON sh.story_id = st.id
JOIN "user" u ON st.author = u.id
WHERE sh.owner_id = $1
ORDER BY st.id
LIMIT 30
OFFSET 30 * ($2 - 1);

-- name: GetOnShelf :one
SELECT * FROM shelf
WHERE owner_id = $1 AND story_id = $2
LIMIT 1;

-- name: AddToShelf :exec
INSERT INTO shelf (owner_id, story_id)
VALUES ($1, $2)
RETURNING *;

-- name: RemoveFromShelf :exec
DELETE FROM shelf
WHERE owner_id = $1 AND story_id = $2;