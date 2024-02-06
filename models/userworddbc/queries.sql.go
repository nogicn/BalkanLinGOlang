// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: queries.sql

package userworddbc

import (
	"context"
	"database/sql"
)

const createUserIndex = `-- name: CreateUserIndex :exec
CREATE INDEX user_word_id_index ON user_word (user_id);

INSERT INTO user_word (last_answered, delay, active, word_id, user_id) 
VALUES (?1, ?2, ?3, ?4, ?5)
`

type CreateUserIndexParams struct {
	LastAnswered sql.NullString
	Delay        sql.NullInt64
	Active       int64
	WordId       int64
	UserId       int64
}

func (q *Queries) CreateUserIndex(ctx context.Context, arg CreateUserIndexParams) error {
	_, err := q.db.ExecContext(ctx, createUserIndex,
		arg.LastAnswered,
		arg.Delay,
		arg.Active,
		arg.WordId,
		arg.UserId,
	)
	return err
}

const deactivateWordForUser = `-- name: DeactivateWordForUser :exec
UPDATE user_word
SET active = 0
WHERE user_id = ?1
AND word_id = ?2
`

type DeactivateWordForUserParams struct {
	UserId int64
	WordId int64
}

func (q *Queries) DeactivateWordForUser(ctx context.Context, arg DeactivateWordForUserParams) error {
	_, err := q.db.ExecContext(ctx, deactivateWordForUser, arg.UserId, arg.WordId)
	return err
}

const deleteUserWordbyID = `-- name: DeleteUserWordbyID :exec
DELETE FROM user_word WHERE word_id = ?1
`

func (q *Queries) DeleteUserWordbyID(ctx context.Context, wordid int64) error {
	_, err := q.db.ExecContext(ctx, deleteUserWordbyID, wordid)
	return err
}

const getDelayForWordForUser = `-- name: GetDelayForWordForUser :one
SELECT delay FROM user_word WHERE user_id = ?1 AND word_id = ?2
`

type GetDelayForWordForUserParams struct {
	UserId int64
	WordId int64
}

func (q *Queries) GetDelayForWordForUser(ctx context.Context, arg GetDelayForWordForUserParams) (sql.NullInt64, error) {
	row := q.db.QueryRowContext(ctx, getDelayForWordForUser, arg.UserId, arg.WordId)
	var delay sql.NullInt64
	err := row.Scan(&delay)
	return delay, err
}

const getUserWordByUserID = `-- name: GetUserWordByUserID :many
SELECT id, last_answered, delay, active, word_id, user_id FROM user_word WHERE user_id = ?1
`

func (q *Queries) GetUserWordByUserID(ctx context.Context, userid int64) ([]UserWord, error) {
	rows, err := q.db.QueryContext(ctx, getUserWordByUserID, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UserWord
	for rows.Next() {
		var i UserWord
		if err := rows.Scan(
			&i.ID,
			&i.LastAnswered,
			&i.Delay,
			&i.Active,
			&i.WordID,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getViableWordsForUserForDictionary = `-- name: GetViableWordsForUserForDictionary :many
SELECT user_word.id, user_word.last_answered, user_word.delay, user_word.active, user_word.word_id, user_word.user_id, word.id, word.foreignword, word.foreigndescription, word.nativeword, word.nativedescription, word.pronunciation, word.dictionary_id
FROM user_word, word
WHERE user_word.word_id = word.id
AND user_word.user_id = ?1
AND word.dictionary_id = ?2
AND strftime('%s', 'now') - strftime('%s ', SUBSTR(last_answered, 1, 19)) > delay * 24 * 60 * 60
`

type GetViableWordsForUserForDictionaryParams struct {
	UserId       int64
	DictionaryId int64
}

type GetViableWordsForUserForDictionaryRow struct {
	ID                 int64
	LastAnswered       sql.NullString
	Delay              sql.NullInt64
	Active             int64
	WordID             int64
	UserID             int64
	ID_2               int64
	Foreignword        string
	Foreigndescription string
	Nativeword         string
	Nativedescription  string
	Pronunciation      string
	DictionaryID       int64
}

func (q *Queries) GetViableWordsForUserForDictionary(ctx context.Context, arg GetViableWordsForUserForDictionaryParams) ([]GetViableWordsForUserForDictionaryRow, error) {
	rows, err := q.db.QueryContext(ctx, getViableWordsForUserForDictionary, arg.UserId, arg.DictionaryId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetViableWordsForUserForDictionaryRow
	for rows.Next() {
		var i GetViableWordsForUserForDictionaryRow
		if err := rows.Scan(
			&i.ID,
			&i.LastAnswered,
			&i.Delay,
			&i.Active,
			&i.WordID,
			&i.UserID,
			&i.ID_2,
			&i.Foreignword,
			&i.Foreigndescription,
			&i.Nativeword,
			&i.Nativedescription,
			&i.Pronunciation,
			&i.DictionaryID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getViableWordsForUserForDictionaryWhereItIsntActiveQuestion = `-- name: GetViableWordsForUserForDictionaryWhereItIsntActiveQuestion :many
SELECT user_word.id, user_word.last_answered, user_word.delay, user_word.active, user_word.word_id, user_word.user_id, word.id, word.foreignword, word.foreigndescription, word.nativeword, word.nativedescription, word.pronunciation, word.dictionary_id
FROM user_word, word
WHERE user_word.word_id = word.id
AND user_word.user_id = ?1
AND word.dictionary_id = ?2
AND word.id NOT IN (
    SELECT word_id FROM active_question WHERE user_id = ?1
)
AND strftime('%s', 'now') - strftime('%s ', SUBSTR(last_answered, 1, 19)) > delay * 24 * 60 * 60
`

type GetViableWordsForUserForDictionaryWhereItIsntActiveQuestionParams struct {
	UserId       int64
	DictionaryId int64
}

type GetViableWordsForUserForDictionaryWhereItIsntActiveQuestionRow struct {
	ID                 int64
	LastAnswered       sql.NullString
	Delay              sql.NullInt64
	Active             int64
	WordID             int64
	UserID             int64
	ID_2               int64
	Foreignword        string
	Foreigndescription string
	Nativeword         string
	Nativedescription  string
	Pronunciation      string
	DictionaryID       int64
}

func (q *Queries) GetViableWordsForUserForDictionaryWhereItIsntActiveQuestion(ctx context.Context, arg GetViableWordsForUserForDictionaryWhereItIsntActiveQuestionParams) ([]GetViableWordsForUserForDictionaryWhereItIsntActiveQuestionRow, error) {
	rows, err := q.db.QueryContext(ctx, getViableWordsForUserForDictionaryWhereItIsntActiveQuestion, arg.UserId, arg.DictionaryId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetViableWordsForUserForDictionaryWhereItIsntActiveQuestionRow
	for rows.Next() {
		var i GetViableWordsForUserForDictionaryWhereItIsntActiveQuestionRow
		if err := rows.Scan(
			&i.ID,
			&i.LastAnswered,
			&i.Delay,
			&i.Active,
			&i.WordID,
			&i.UserID,
			&i.ID_2,
			&i.Foreignword,
			&i.Foreigndescription,
			&i.Nativeword,
			&i.Nativedescription,
			&i.Pronunciation,
			&i.DictionaryID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getWordsForUserForDictionary = `-- name: GetWordsForUserForDictionary :many
SELECT user_word.id, user_word.last_answered, user_word.delay, user_word.active, user_word.word_id, user_word.user_id, word.id, word.foreignword, word.foreigndescription, word.nativeword, word.nativedescription, word.pronunciation, word.dictionary_id
FROM user_word, word
WHERE user_word.word_id = word.id
AND user_word.user_id = ?1
AND word.dictionary_id = ?2
`

type GetWordsForUserForDictionaryParams struct {
	UserId       int64
	DictionaryId int64
}

type GetWordsForUserForDictionaryRow struct {
	ID                 int64
	LastAnswered       sql.NullString
	Delay              sql.NullInt64
	Active             int64
	WordID             int64
	UserID             int64
	ID_2               int64
	Foreignword        string
	Foreigndescription string
	Nativeword         string
	Nativedescription  string
	Pronunciation      string
	DictionaryID       int64
}

func (q *Queries) GetWordsForUserForDictionary(ctx context.Context, arg GetWordsForUserForDictionaryParams) ([]GetWordsForUserForDictionaryRow, error) {
	rows, err := q.db.QueryContext(ctx, getWordsForUserForDictionary, arg.UserId, arg.DictionaryId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetWordsForUserForDictionaryRow
	for rows.Next() {
		var i GetWordsForUserForDictionaryRow
		if err := rows.Scan(
			&i.ID,
			&i.LastAnswered,
			&i.Delay,
			&i.Active,
			&i.WordID,
			&i.UserID,
			&i.ID_2,
			&i.Foreignword,
			&i.Foreigndescription,
			&i.Nativeword,
			&i.Nativedescription,
			&i.Pronunciation,
			&i.DictionaryID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const setNewDelayForUser = `-- name: SetNewDelayForUser :exec
UPDATE user_word
SET delay = CASE
    WHEN ?1 = 0 THEN delay = 0
    ELSE delay + 1
END
WHERE user_id = ?2
AND word_id = ?3
`

type SetNewDelayForUserParams struct {
	Delay  interface{}
	UserId int64
	WordId int64
}

func (q *Queries) SetNewDelayForUser(ctx context.Context, arg SetNewDelayForUserParams) error {
	_, err := q.db.ExecContext(ctx, setNewDelayForUser, arg.Delay, arg.UserId, arg.WordId)
	return err
}

const updateLastAnswered = `-- name: UpdateLastAnswered :exec
UPDATE user_word
SET last_answered = ?1
WHERE user_id = ?2
AND word_id = ?3
`

type UpdateLastAnsweredParams struct {
	LastAnswered sql.NullString
	UserId       int64
	WordId       int64
}

func (q *Queries) UpdateLastAnswered(ctx context.Context, arg UpdateLastAnsweredParams) error {
	_, err := q.db.ExecContext(ctx, updateLastAnswered, arg.LastAnswered, arg.UserId, arg.WordId)
	return err
}
