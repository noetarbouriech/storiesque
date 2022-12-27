-- name: GetPage :one
SELECT * FROM page
WHERE id = $1 LIMIT 1;

-- name: CreatePage :one
INSERT INTO page (action, body)
VALUES ($1, $2)
RETURNING *;

-- name: UpdatePage :exec
UPDATE page
SET
  action = CASE WHEN @action_do_update::boolean
    THEN @action::VARCHAR(32) ELSE action END,

  body = CASE WHEN @body_do_update::boolean
    THEN @body::VARCHAR(4096) ELSE body END
WHERE id = $1;

-- name: DeletePage :exec
DELETE FROM page
WHERE id = $1;
