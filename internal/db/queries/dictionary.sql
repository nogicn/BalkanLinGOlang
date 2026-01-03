-- name: CreateDictionary :exec
INSERT INTO dictionary (name, language_id, image_link) VALUES (sqlc.arg(name), sqlc.arg(language_id), sqlc.arg(image_link));

-- name: GetDictionariesForUser :many
SELECT dictionary.*, language.flag_icon
FROM dictionary
LEFT JOIN language ON dictionary.language_id = language.id
LEFT JOIN dictionary_user ON dictionary.id = dictionary_user.dictionary_id
WHERE dictionary_user.dictionary_id = dictionary.id
AND dictionary_user.user_id = sqlc.arg(user_id);

-- name: GetAllDictionaries :many
SELECT * FROM dictionary;

-- name: GetAllDictionariesWithIcons :many
SELECT dictionary.*, language.flag_icon
FROM dictionary 
LEFT JOIN language ON dictionary.language_id = language.id;

-- name: DeleteDictionary :exec
DELETE FROM dictionary WHERE id = sqlc.arg(id);

-- name: GetDictionariesNotAssignedToUser :many
SELECT dictionary.*
FROM dictionary
WHERE dictionary.id NOT IN (
    SELECT dictionary_id
    FROM dictionary_user
    WHERE user_id = sqlc.arg(user_id)
);

-- name: GetDictionaryByID :one
SELECT * FROM dictionary WHERE id = sqlc.arg(id);

-- name: UpdateDictionary :exec
UPDATE dictionary SET name = sqlc.arg(name), language_id = sqlc.arg(language_id), image_link = sqlc.arg(image_link) WHERE id = sqlc.arg(id);
