// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package todos

import (
	"context"
	"database/sql"
)

const createTodo = `-- name: CreateTodo :exec
INSERT INTO todos (task, created_by, created_date, updated_date) VALUES ($1, $2, $3, $4) RETURNING id, task, created_by, created_date, updated_date
`

type CreateTodoParams struct {
	Task        sql.NullString
	CreatedBy   sql.NullInt64
	CreatedDate sql.NullTime
	UpdatedDate sql.NullTime
}

func (q *Queries) CreateTodo(ctx context.Context, arg CreateTodoParams) error {
	_, err := q.db.ExecContext(ctx, createTodo,
		arg.Task,
		arg.CreatedBy,
		arg.CreatedDate,
		arg.UpdatedDate,
	)
	return err
}

const deleteTodo = `-- name: DeleteTodo :exec
DELETE FROM todos WHERE id = $1 RETURNING id, task, created_by, created_date, updated_date
`

func (q *Queries) DeleteTodo(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteTodo, id)
	return err
}

const getTodo = `-- name: GetTodo :one
SELECT id, task, created_by, created_date, updated_date FROM todos WHERE id = $1 LIMIT 1
`

func (q *Queries) GetTodo(ctx context.Context, id int32) (Todo, error) {
	row := q.db.QueryRowContext(ctx, getTodo, id)
	var i Todo
	err := row.Scan(
		&i.ID,
		&i.Task,
		&i.CreatedBy,
		&i.CreatedDate,
		&i.UpdatedDate,
	)
	return i, err
}

const getTodos = `-- name: GetTodos :many
SELECT id, task, created_by, created_date, updated_date FROM todos order by created_date desc
`

func (q *Queries) GetTodos(ctx context.Context) ([]Todo, error) {
	rows, err := q.db.QueryContext(ctx, getTodos)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Todo
	for rows.Next() {
		var i Todo
		if err := rows.Scan(
			&i.ID,
			&i.Task,
			&i.CreatedBy,
			&i.CreatedDate,
			&i.UpdatedDate,
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

const updateTodo = `-- name: UpdateTodo :exec
UPDATE todos SET task = $1, updated_date = $2 WHERE id = $3 RETURNING id, task, created_by, created_date, updated_date
`

type UpdateTodoParams struct {
	Task        sql.NullString
	UpdatedDate sql.NullTime
	ID          int32
}

func (q *Queries) UpdateTodo(ctx context.Context, arg UpdateTodoParams) error {
	_, err := q.db.ExecContext(ctx, updateTodo, arg.Task, arg.UpdatedDate, arg.ID)
	return err
}
