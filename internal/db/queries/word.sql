-- name: CreateWord :exec
INSERT INTO word (foreignWord, foreignDescription, nativeWord, nativeDescription, pronunciation, dictionary_id) 
VALUES (sqlc.arg(foreign_word), sqlc.arg(foreign_description), sqlc.arg(native_word), sqlc.arg(native_description), sqlc.arg(pronunciation), sqlc.arg(dictionary_id));

-- name: DeleteWordByID :exec
DELETE FROM word WHERE id = sqlc.arg(id);

-- name: DeleteWordByMeaning :exec
DELETE FROM word WHERE foreignWord = sqlc.arg(foreign_word) AND foreignDescription = sqlc.arg(foreign_description) AND nativeWord = sqlc.arg(native_word) AND nativeDescription = sqlc.arg(native_description);

-- name: GetWordsByDictionaryID :many
SELECT * FROM word WHERE dictionary_id = sqlc.arg(dictionary_id);

-- name: DeleteWordsByDictionaryID :exec
DELETE FROM word WHERE dictionary_id = sqlc.arg(dictionary_id);

-- name: GetAllWords :many
SELECT * FROM word;

-- name: GetWordByID :one
SELECT * FROM word WHERE id = sqlc.arg(id);

-- name: UpdateWord :exec
UPDATE word SET foreignWord = sqlc.arg(foreign_word), foreignDescription = sqlc.arg(foreign_description), nativeWord = sqlc.arg(native_word), nativeDescription = sqlc.arg(native_description), pronunciation = sqlc.arg(pronunciation) WHERE id = sqlc.arg(id);

-- name: SearchWordByDictionaryID :many
SELECT * FROM word WHERE dictionary_id = sqlc.arg(dictionary_id) AND (foreignWord LIKE '%' || sqlc.arg(search_term) || '%' OR nativeWord LIKE '%' || sqlc.arg(search_term) || '%');

-- name: GetAllWordsNotInUserWord :many
SELECT * FROM word WHERE id NOT IN (SELECT word_id as id FROM user_word WHERE user_id = sqlc.arg(user_id));
