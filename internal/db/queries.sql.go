// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: queries.sql

package db

import (
	"context"
	"database/sql"
)

const createStory = `-- name: CreateStory :one
INSERT INTO story (title, description)
VALUES ($1, $2)
RETURNING id, title, description, first_page_id
`

type CreateStoryParams struct {
	Title       string
	Description sql.NullString
}

func (q *Queries) CreateStory(ctx context.Context, arg CreateStoryParams) (Story, error) {
	row := q.db.QueryRowContext(ctx, createStory, arg.Title, arg.Description)
	var i Story
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.FirstPageID,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (username, password_hash, email)
VALUES ($1, $2, $3)
RETURNING id, username, password_hash, is_admin, email
`

type CreateUserParams struct {
	Username     string
	PasswordHash string
	Email        string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Username, arg.PasswordHash, arg.Email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.PasswordHash,
		&i.IsAdmin,
		&i.Email,
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

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getStory = `-- name: GetStory :one
SELECT id, title, description, first_page_id FROM story
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetStory(ctx context.Context, id int64) (Story, error) {
	row := q.db.QueryRowContext(ctx, getStory, id)
	var i Story
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.FirstPageID,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, username, password_hash, is_admin, email FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.PasswordHash,
		&i.IsAdmin,
		&i.Email,
	)
	return i, err
}

const getUserWithEmail = `-- name: GetUserWithEmail :one
SELECT id, username, password_hash, is_admin, email FROM users
WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUserWithEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserWithEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.PasswordHash,
		&i.IsAdmin,
		&i.Email,
	)
	return i, err
}

const getUserWithUsername = `-- name: GetUserWithUsername :one
SELECT id, username, password_hash, is_admin, email FROM users
WHERE username = $1 LIMIT 1
`

func (q *Queries) GetUserWithUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserWithUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.PasswordHash,
		&i.IsAdmin,
		&i.Email,
	)
	return i, err
}

const listStories = `-- name: ListStories :many
SELECT id, title, description, first_page_id FROM story
ORDER BY title
`

func (q *Queries) ListStories(ctx context.Context) ([]Story, error) {
	rows, err := q.db.QueryContext(ctx, listStories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Story
	for rows.Next() {
		var i Story
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.FirstPageID,
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

const listUsers = `-- name: ListUsers :many
SELECT id, username, password_hash, is_admin, email FROM users
ORDER BY username
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.PasswordHash,
			&i.IsAdmin,
			&i.Email,
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

const updateUser = `-- name: UpdateUser :exec
UPDATE users
SET
  username = CASE WHEN $2::boolean
    THEN $3::VARCHAR(24) ELSE username END,

  password_hash = CASE WHEN $4::boolean
    THEN $5::VARCHAR(64) ELSE password_hash END,

  email = CASE WHEN $6::boolean
    THEN $7::VARCHAR(128) ELSE email END
WHERE id = $1
`

type UpdateUserParams struct {
	ID                   int64
	UsernameDoUpdate     bool
	Username             string
	PasswordHashDoUpdate bool
	PasswordHash         string
	EmailDoUpdate        bool
	Email                string
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.ExecContext(ctx, updateUser,
		arg.ID,
		arg.UsernameDoUpdate,
		arg.Username,
		arg.PasswordHashDoUpdate,
		arg.PasswordHash,
		arg.EmailDoUpdate,
		arg.Email,
	)
	return err
}
