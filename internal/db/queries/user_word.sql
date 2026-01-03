-- name: CreateUserWord :exec
INSERT INTO user_word (last_answered, delay, active, word_id, user_id) VALUES (sqlc.arg(last_answered), sqlc.arg(delay), sqlc.arg(active), sqlc.arg(word_id), sqlc.arg(user_id));

-- name: GetWordsForUserForDictionary :many
SELECT sqlc.embed(user_word), sqlc.embed(word)
FROM user_word, word
WHERE user_word.word_id = word.id
AND user_word.user_id = sqlc.arg(user_id)
AND word.dictionary_id = sqlc.arg(dictionary_id);

-- name: GetViableWordsForUserForDictionary :many
SELECT sqlc.embed(user_word), sqlc.embed(word)
FROM user_word, word
WHERE user_word.word_id = word.id
AND user_word.user_id = sqlc.arg(user_id)
AND word.dictionary_id = sqlc.arg(dictionary_id)
AND strftime('%s', 'now') - strftime('%s ', SUBSTR(last_answered, 1, 19)) > delay * 24 * 60 * 60;

-- name: GetViableWordsForUserForDictionaryWhereItIsntActiveQuestion :many
SELECT sqlc.embed(user_word), sqlc.embed(word)
FROM user_word, word
WHERE user_word.word_id = word.id
AND user_word.user_id = sqlc.arg(user_id)
AND word.dictionary_id = sqlc.arg(dictionary_id)
AND word.id NOT IN (
    SELECT word_id FROM active_question WHERE active_question.user_id = sqlc.arg(user_id)
)
AND strftime('%s', 'now') - strftime('%s ', SUBSTR(last_answered, 1, 19)) > delay * 24 * 60 * 60;

-- name: SetNewDelayForUser :exec
UPDATE user_word
SET delay = CASE
    WHEN sqlc.arg(is_correct) = 0 THEN 0
    ELSE delay + 1
END
WHERE user_id = sqlc.arg(user_id) AND word_id = sqlc.arg(word_id);

-- name: DeactivateWordForUser :exec
UPDATE user_word SET active = 0 WHERE user_id = sqlc.arg(user_id) AND word_id = sqlc.arg(word_id);

-- name: GetUserWordByUserID :many
SELECT * FROM user_word WHERE user_id = sqlc.arg(user_id);

-- name: GetDelayForWordForUser :one
SELECT delay FROM user_word WHERE user_id = sqlc.arg(user_id) AND word_id = sqlc.arg(word_id);

-- name: UpdateLastAnswered :exec
UPDATE user_word SET last_answered = sqlc.arg(last_answered) WHERE user_id = sqlc.arg(user_id) AND word_id = sqlc.arg(word_id);

-- name: DeleteUserWordByID :exec
DELETE FROM user_word WHERE word_id = sqlc.arg(word_id);
