// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: user.sql

package db

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :one
INSERT INTO "user" (username, password_hash, email)
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

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM "user"
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUserDetails = `-- name: GetUserDetails :many
SELECT u.id, u.username, u.password_hash, u.is_admin, u.email, s.id as story_id, s.title, s.description FROM "user" u
LEFT JOIN story s ON u.id = s.author
WHERE username = $1
`

type GetUserDetailsRow struct {
	ID           int64
	Username     string
	PasswordHash string
	IsAdmin      bool
	Email        string
	StoryID      sql.NullInt64
	Title        sql.NullString
	Description  sql.NullString
}

func (q *Queries) GetUserDetails(ctx context.Context, username string) ([]GetUserDetailsRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserDetails, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserDetailsRow
	for rows.Next() {
		var i GetUserDetailsRow
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.PasswordHash,
			&i.IsAdmin,
			&i.Email,
			&i.StoryID,
			&i.Title,
			&i.Description,
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

const getUserWithEmail = `-- name: GetUserWithEmail :one
SELECT id, username, password_hash, is_admin, email FROM "user"
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

const getUserWithId = `-- name: GetUserWithId :one
SELECT id, username, password_hash, is_admin, email FROM "user"
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUserWithId(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserWithId, id)
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
SELECT id, username, password_hash, is_admin, email FROM "user"
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

const searchUsers = `-- name: SearchUsers :many
SELECT id, username, password_hash, is_admin, email FROM "user"
WHERE username LIKE '%' || $1 || '%'
ORDER BY id
LIMIT 40
OFFSET 40 * ($2 - 1)
`

type SearchUsersParams struct {
	Column1 sql.NullString
	Column2 interface{}
}

func (q *Queries) SearchUsers(ctx context.Context, arg SearchUsersParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, searchUsers, arg.Column1, arg.Column2)
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

const setAdmin = `-- name: SetAdmin :exec
UPDATE "user"
SET is_admin = NOT is_admin
WHERE id = $1
`

func (q *Queries) SetAdmin(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, setAdmin, id)
	return err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE "user"
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
