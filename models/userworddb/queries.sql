-- name: CreateUserIndex :exec
CREATE INDEX user_word_id_index ON user_word (user_id);

-- name: CreateUserWord :exec
INSERT INTO user_word (last_answered, delay, active, word_id, user_id) 
VALUES (@lastAnswered, @delay, @active, @wordId, @userId);

-- name: GetWordsForUserForDictionary :many
SELECT user_word.*, word.*
FROM user_word, word
WHERE user_word.word_id = word.id
AND user_word.user_id = @userId
AND word.dictionary_id = @dictionaryId;

-- name: GetViableWordsForUserForDictionary :many
SELECT user_word.*, word.*
FROM user_word, word
WHERE user_word.word_id = word.id
AND user_word.user_id = @userId
AND word.dictionary_id = @dictionaryId
AND strftime('%s', 'now') - strftime('%s ', SUBSTR(last_answered, 1, 19)) > delay * 24 * 60 * 60;

-- name: GetViableWordsForUserForDictionaryWhereItIsntActiveQuestion :many
SELECT user_word.*, word.*
FROM user_word, word
WHERE user_word.word_id = word.id
AND user_word.user_id = @userId
AND word.dictionary_id = @dictionaryId
AND word.id NOT IN (
    SELECT word_id FROM active_question WHERE user_id = @userId
)
AND strftime('%s', 'now') - strftime('%s ', SUBSTR(last_answered, 1, 19)) > delay * 24 * 60 * 60;

-- name: SetNewDelayForUser :exec
UPDATE user_word
SET delay = CASE
    WHEN @delay = 0 THEN delay = 0
    ELSE delay + 1
END
WHERE user_id = @userId
AND word_id = @wordId;

-- name: DeactivateWordForUser :exec
UPDATE user_word
SET active = 0
WHERE user_id = @userId
AND word_id = @wordId;

-- name: GetUserWordByUserID :many
SELECT * FROM user_word WHERE user_id = @userId;

-- name: GetDelayForWordForUser :one
SELECT delay FROM user_word WHERE user_id = @userId AND word_id = @wordId;

-- name: UpdateLastAnswered :exec
UPDATE user_word
SET last_answered = @lastAnswered
WHERE user_id = @userId
AND word_id = @wordId;

-- name: DeleteUserWordbyID :exec
DELETE FROM user_word WHERE word_id = @wordId;
