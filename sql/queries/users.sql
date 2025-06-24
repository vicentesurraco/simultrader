-- name: CreateUser :one
INSERT INTO users (name, password_hash, email, created_at, updated_at)
VALUES ($1, $2, $3, DEFAULT, DEFAULT)
RETURNING id, name, email, created_at, updated_at;

-- name: GetUser :one
select id, name, password_hash, email, created_at, updated_at
from users
where name = $1
;
