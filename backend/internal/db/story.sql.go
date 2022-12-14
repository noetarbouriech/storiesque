// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: story.sql

package db

import (
	"context"
	"database/sql"
)

const createStory = `-- name: CreateStory :one
INSERT INTO story (title, description, author)
VALUES ($1, $2, $3)
RETURNING id, title, description, has_img, featured, author, first_page_id
`

type CreateStoryParams struct {
	Title       string
	Description sql.NullString
	Author      int64
}

func (q *Queries) CreateStory(ctx context.Context, arg CreateStoryParams) (Story, error) {
	row := q.db.QueryRowContext(ctx, createStory, arg.Title, arg.Description, arg.Author)
	var i Story
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.HasImg,
		&i.Featured,
		&i.Author,
		&i.FirstPageID,
	)
	return i, err
}

const deleteStory = `-- name: DeleteStory :exec
DELETE FROM story
WHERE id = $1
`

func (q *Queries) DeleteStory(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteStory, id)
	return err
}

const getFeaturedStories = `-- name: GetFeaturedStories :many
SELECT s.id, s.title, s.description, s.has_img, u.username as author_name FROM story s
JOIN "user" u ON s.author = u.id
WHERE s.featured = true
ORDER BY s.id
LIMIT 3
`

type GetFeaturedStoriesRow struct {
	ID          int64
	Title       string
	Description sql.NullString
	HasImg      bool
	AuthorName  string
}

func (q *Queries) GetFeaturedStories(ctx context.Context) ([]GetFeaturedStoriesRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeaturedStories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeaturedStoriesRow
	for rows.Next() {
		var i GetFeaturedStoriesRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.HasImg,
			&i.AuthorName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getStory = `-- name: GetStory :one
SELECT s.id, s.title, s.description, s.first_page_id, s.has_img, u.id as author, u.username as author_name FROM story s
JOIN "user" u ON s.author = u.id
WHERE s.id = $1 LIMIT 1
`

type GetStoryRow struct {
	ID          int64
	Title       string
	Description sql.NullString
	FirstPageID sql.NullInt64
	HasImg      bool
	Author      int64
	AuthorName  string
}

func (q *Queries) GetStory(ctx context.Context, id int64) (GetStoryRow, error) {
	row := q.db.QueryRowContext(ctx, getStory, id)
	var i GetStoryRow
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.FirstPageID,
		&i.HasImg,
		&i.Author,
		&i.AuthorName,
	)
	return i, err
}

const getStoryAuthor = `-- name: GetStoryAuthor :one
SELECT author FROM story
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetStoryAuthor(ctx context.Context, id int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, getStoryAuthor, id)
	var author int64
	err := row.Scan(&author)
	return author, err
}

const searchStories = `-- name: SearchStories :many
SELECT s.id, s.title, s.description, s.has_img, s.featured, u.username as author_name FROM story s
JOIN "user" u ON s.author = u.id
WHERE title LIKE '%' || $1 || '%'
ORDER BY s.id
LIMIT 30
OFFSET 30 * ($2 - 1)
`

type SearchStoriesParams struct {
	Column1 sql.NullString
	Column2 interface{}
}

type SearchStoriesRow struct {
	ID          int64
	Title       string
	Description sql.NullString
	HasImg      bool
	Featured    bool
	AuthorName  string
}

func (q *Queries) SearchStories(ctx context.Context, arg SearchStoriesParams) ([]SearchStoriesRow, error) {
	rows, err := q.db.QueryContext(ctx, searchStories, arg.Column1, arg.Column2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SearchStoriesRow
	for rows.Next() {
		var i SearchStoriesRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.HasImg,
			&i.Featured,
			&i.AuthorName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const setImgStory = `-- name: SetImgStory :exec
UPDATE story
SET has_img = $2
WHERE id = $1
`

type SetImgStoryParams struct {
	ID     int64
	HasImg bool
}

func (q *Queries) SetImgStory(ctx context.Context, arg SetImgStoryParams) error {
	_, err := q.db.ExecContext(ctx, setImgStory, arg.ID, arg.HasImg)
	return err
}

const setStoryAsFeatured = `-- name: SetStoryAsFeatured :exec
UPDATE story
SET featured = NOT featured
WHERE id = $1
`

func (q *Queries) SetStoryAsFeatured(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, setStoryAsFeatured, id)
	return err
}

const updateStory = `-- name: UpdateStory :exec
UPDATE story
SET
  title = CASE WHEN $2::boolean
    THEN $3::VARCHAR(48) ELSE title END,

  description = CASE WHEN $4::boolean
    THEN $5::VARCHAR(512) ELSE description END
WHERE id = $1
`

type UpdateStoryParams struct {
	ID                  int64
	TitleDoUpdate       bool
	Title               string
	DescriptionDoUpdate bool
	Description         string
}

func (q *Queries) UpdateStory(ctx context.Context, arg UpdateStoryParams) error {
	_, err := q.db.ExecContext(ctx, updateStory,
		arg.ID,
		arg.TitleDoUpdate,
		arg.Title,
		arg.DescriptionDoUpdate,
		arg.Description,
	)
	return err
}
