-- name: CreateUser :exec
INSERT INTO user (name, surname, email, password) VALUES (sqlc.arg(name), sqlc.arg(surname), sqlc.arg(email), sqlc.arg(password));

-- name: CreateAdmin :exec
INSERT INTO user (name, surname, email, password, is_admin) VALUES (sqlc.arg(name), sqlc.arg(surname), sqlc.arg(email), sqlc.arg(password), 1);

-- name: LoginEmailPassword :one
SELECT * FROM user WHERE email = sqlc.arg(email) AND password = sqlc.arg(password);

-- name: GetAllUsers :many
SELECT * FROM user;

-- name: GetUserByID :one
SELECT * FROM user WHERE id = sqlc.arg(id);

-- name: GetUserByEmail :one
SELECT * FROM user WHERE email = sqlc.arg(email);

-- name: UpdateTokenByEmail :one
UPDATE user SET token = sqlc.arg(token) WHERE email = sqlc.arg(email) RETURNING *;

-- name: UpdateTokenByID :one
UPDATE user SET token = sqlc.arg(token) WHERE id = sqlc.arg(id) RETURNING *;

-- name: UpdatePasswordByEmail :one
UPDATE user SET password = sqlc.arg(password) WHERE email = sqlc.arg(email) RETURNING *;

-- name: UpdateUserByToken :one
UPDATE user SET name = sqlc.arg(name), surname = sqlc.arg(surname) WHERE token = sqlc.arg(token) RETURNING *;

-- name: DeleteUserByID :exec
DELETE FROM user WHERE id = sqlc.arg(id);

-- name: SetAdminByEmail :one
UPDATE user SET is_admin = NOT is_admin WHERE email = sqlc.arg(email) RETURNING *;

-- name: SetAdminByID :one
UPDATE user SET is_admin = NOT is_admin WHERE id = sqlc.arg(id) RETURNING *;

-- name: GetAllUsersLikeEmail :many
SELECT * FROM user WHERE email LIKE '%' || sqlc.arg(email) || '%';
