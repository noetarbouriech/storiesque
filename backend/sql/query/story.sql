-- name: GetStory :one
SELECT s.id, s.title, s.description, s.first_page_id, s.has_img, u.id as author, u.username as author_name FROM story s
JOIN "user" u ON s.author = u.id
WHERE s.id = $1 LIMIT 1;

-- name: GetStoryAuthor :one
SELECT author FROM story
WHERE id = $1 LIMIT 1;

-- name: SearchStories :many
SELECT s.id, s.title, s.description, s.has_img, s.featured, u.username as author_name FROM story s
JOIN "user" u ON s.author = u.id
WHERE title LIKE '%' || $1 || '%'
ORDER BY s.id
LIMIT 30
OFFSET 30 * ($2 - 1);

-- name: CreateStory :one
INSERT INTO story (title, description, author)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateStory :exec
UPDATE story
SET
  title = CASE WHEN @title_do_update::boolean
    THEN @title::VARCHAR(48) ELSE title END,

  description = CASE WHEN @description_do_update::boolean
    THEN @description::VARCHAR(512) ELSE description END
WHERE id = $1;

-- name: DeleteStory :exec
DELETE FROM story
WHERE id = $1;

-- name: SetImgStory :exec
UPDATE story
SET has_img = $2
WHERE id = $1;

-- name: GetFeaturedStories :many
SELECT s.id, s.title, s.description, s.has_img, u.username as author_name FROM story s
JOIN "user" u ON s.author = u.id
WHERE s.featured = true
ORDER BY s.id
LIMIT 3;

-- name: SetStoryAsFeatured :exec
UPDATE story
SET featured = NOT featured
WHERE id = $1;