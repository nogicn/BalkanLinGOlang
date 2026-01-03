-- name: AddDictionaryToUser :exec
INSERT INTO dictionary_user (user_id, dictionary_id) VALUES (sqlc.arg(user_id), sqlc.arg(dictionary_id));

-- name: DeleteDictionaryFromUser :exec
DELETE FROM dictionary_user WHERE user_id = sqlc.arg(user_id) AND dictionary_id = sqlc.arg(dictionary_id);

-- name: GetDictionaryUserByID :one
SELECT * FROM dictionary_user WHERE id = sqlc.arg(id);

-- name: GetUserDictionaries :many
SELECT * FROM dictionary_user WHERE user_id = sqlc.arg(user_id) AND dictionary_id = sqlc.arg(dictionary_id);
