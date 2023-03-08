// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0
// source: crud.sql

package repository

import (
	"context"
)

const createReport = `-- name: CreateReport :exec
INSERT INTO "report" (
    "product_id", "description"
) VALUES (
    $1, $2
)
`

type CreateReportParams struct {
	ProductID   int64
	Description string
}

func (q *Queries) CreateReport(ctx context.Context, arg CreateReportParams) error {
	_, err := q.db.ExecContext(ctx, createReport, arg.ProductID, arg.Description)
	return err
}

const deleteReport = `-- name: DeleteReport :exec
DELETE FROM "report"
WHERE "id" = $1
`

func (q *Queries) DeleteReport(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteReport, id)
	return err
}

const getAllReport = `-- name: GetAllReport :many
SELECT id, product_id, description, status FROM "report" WHERE "status" = 'waiting'
`

func (q *Queries) GetAllReport(ctx context.Context) ([]Report, error) {
	rows, err := q.db.QueryContext(ctx, getAllReport)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Report
	for rows.Next() {
		var i Report
		if err := rows.Scan(
			&i.ID,
			&i.ProductID,
			&i.Description,
			&i.Status,
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

const handleReport = `-- name: HandleReport :one
UPDATE "report"
SET "status" = 'handled'
WHERE "id" = $1
RETURNING id, product_id, description, status
`

func (q *Queries) HandleReport(ctx context.Context, id int64) (Report, error) {
	row := q.db.QueryRowContext(ctx, handleReport, id)
	var i Report
	err := row.Scan(
		&i.ID,
		&i.ProductID,
		&i.Description,
		&i.Status,
	)
	return i, err
}