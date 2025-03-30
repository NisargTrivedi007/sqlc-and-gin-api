-- name: GetTodo :one
SELECT * FROM todos left join users on users.id = todos.user_id WHERE todos.id = $1 LIMIT 1;

-- name: GetTodos :many
SELECT * FROM todos order by created_date desc;

-- name: CreateTodo :exec
INSERT INTO todos (task, created_by) VALUES ($1, $2) RETURNING *;

-- name: UpdateTodo :exec
UPDATE todos SET task = $1, updated_date = $2 WHERE id = $3 RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todos WHERE id = $1 RETURNING *;

-- name: UpdateTodoStatus :exec
UPDATE todos SET done = $1 WHERE id = $2 RETURNING *;

-- name: GetAllUsers :many
SELECT * FROM users order by created_date desc;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: Register :exec
INSERT INTO users (username, email_id, phone_no,password) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: Login :one
SELECT * FROM users WHERE username = $1 AND password = $2 LIMIT 1;


