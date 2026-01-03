-- name: DeleteActiveQuestionByUserID :exec
DELETE FROM active_question WHERE user_id = sqlc.arg(user_id);

-- name: DeleteActiveQuestionByWordID :exec
DELETE FROM active_question WHERE word_id = sqlc.arg(word_id);

-- name: SetActiveQuestion :exec
INSERT OR REPLACE INTO active_question (user_id, word_id, type) VALUES (sqlc.arg(user_id), sqlc.arg(word_id), sqlc.arg(type));

-- name: GetActiveQuestionByUserID :one
SELECT * FROM active_question WHERE user_id = sqlc.arg(user_id);

-- name: IncreaseActiveQuestionTypeByUserID :exec
UPDATE active_question
SET type = CASE
    WHEN type >= 1 AND type < 4 THEN type + 1
    ELSE 1 
END
WHERE user_id = sqlc.arg(user_id);
