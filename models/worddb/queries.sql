-- name: CreateWord :exec
INSERT INTO word (foreignWord, foreignDescription, nativeWord, nativeDescription, pronunciation, dictionary_id) 
VALUES (@foreignWord, @foreignDescription, @nativeWord, @nativeDescription, @pronunciation, @dictionaryId);

-- name: DeleteWordByID :exec
DELETE FROM word WHERE id = @wordId;

-- name: DeleteWordByMeaning :exec
DELETE FROM word 
WHERE foreignWord = @foreignWord AND foreignDescription = @foreignDescription AND nativeWord = @nativeWord AND nativeDescription = @nativeDescription;

-- name: GetWordByDictionaryID :many
SELECT * FROM word WHERE dictionary_id = @dictionaryId;

-- name: DeleteWordByDictionaryID :exec
DELETE FROM word WHERE dictionary_id = @dictionaryId;

-- name: GetAllWords :many
SELECT * FROM word;

-- name: GetWordByID :one
SELECT * FROM word WHERE id = @wordId;

-- name: UpdateWord :exec
UPDATE word 
SET foreignWord = @foreignWord, foreignDescription = @foreignDescription, nativeWord = @nativeWord, nativeDescription = @nativeDescription, pronunciation = @pronunciation 
WHERE id = @wordId;

-- name: SearchWordByDictionaryID :many
SELECT * FROM word 
WHERE dictionary_id = @dictionaryId AND (foreignWord LIKE '%' || @word || '%' OR nativeWord LIKE '%' || @word || '%');

-- name: GetAllWordsNotInUserWord :many
SELECT * FROM word 
WHERE id NOT IN (SELECT word_id as id FROM user_word WHERE user_id = @userId);
