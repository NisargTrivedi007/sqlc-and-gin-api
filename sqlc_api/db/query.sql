-- name: GetTodo :one
SELECT * FROM todos WHERE id = $1 LIMIT 1;

-- name: GetTodos :many
SELECT * FROM todos order by created_date desc;

-- name: CreateTodo :exec
INSERT INTO todos (task, created_by, created_date, updated_date) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateTodo :exec
UPDATE todos SET task = $1, updated_date = $2 WHERE id = $3 RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todos WHERE id = $1 RETURNING *;
