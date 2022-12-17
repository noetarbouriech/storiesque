-- name: GetPage :one
SELECT * FROM page
WHERE id = $1 LIMIT 1;

-- name: CreatePage :one
INSERT INTO page (title, body)
VALUES ($1, $2)
RETURNING *;

-- name: UpdatePage :exec
UPDATE page
SET
  title = CASE WHEN @title_do_update::boolean
    THEN @title::VARCHAR(32) ELSE title END,

  body = CASE WHEN @body_do_update::boolean
    THEN @body::VARCHAR(4096) ELSE body END
WHERE id = $1;

-- name: DeletePage :exec
DELETE FROM page
WHERE id = $1;
