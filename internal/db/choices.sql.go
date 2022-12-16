// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: choices.sql

package db

import (
	"context"
)

const createChoices = `-- name: CreateChoices :one
INSERT INTO choices (page_id, path_id)
VALUES ($1, $2)
RETURNING page_id, path_id
`

type CreateChoicesParams struct {
	PageID int64
	PathID int64
}

func (q *Queries) CreateChoices(ctx context.Context, arg CreateChoicesParams) (Choice, error) {
	row := q.db.QueryRowContext(ctx, createChoices, arg.PageID, arg.PathID)
	var i Choice
	err := row.Scan(&i.PageID, &i.PathID)
	return i, err
}

const listChoices = `-- name: ListChoices :many
SELECT path_id FROM choices
WHERE page_id = $1
`

func (q *Queries) ListChoices(ctx context.Context, pageID int64) ([]int64, error) {
	rows, err := q.db.QueryContext(ctx, listChoices, pageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var path_id int64
		if err := rows.Scan(&path_id); err != nil {
			return nil, err
		}
		items = append(items, path_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
