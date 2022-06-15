-- name: CreateUser :one
INSERT INTO users (
  email,
  first_name,
  last_name,
  password,
  user_active
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetAll :many
SELECT * FROM users
ORDER BY last_name;

-- name: GetByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET email = $1,
first_name = $2,
last_name = $3,
user_active = $4,
updated_at = $5
WHERE id = $6
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: ResetPassword :one
UPDATE users
SET password = $1
WHERE id = $2
RETURNING *;