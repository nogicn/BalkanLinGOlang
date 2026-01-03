-- name: CreateLanguage :exec
INSERT INTO language (name, shorthand, flag_icon) VALUES (sqlc.arg(name), sqlc.arg(shorthand), sqlc.arg(flag_icon));

-- name: GetAllLanguages :many
SELECT * FROM language;

-- name: GetShorthands :many
SELECT shorthand FROM language;

-- name: GetLanguageByID :one
SELECT * FROM language WHERE id = sqlc.arg(id);

-- name: DeleteLanguageByID :exec
DELETE FROM language WHERE id = sqlc.arg(id);

-- name: UpdateLanguage :exec
UPDATE language SET name = sqlc.arg(name), shorthand = sqlc.arg(shorthand), flag_icon = sqlc.arg(flag_icon) WHERE id = sqlc.arg(id);
