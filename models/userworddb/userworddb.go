package userworddb

import (
	"BalkanLinGO/models/worddb"
	"database/sql"
)

const (
	createUserWordTable = `
        CREATE TABLE IF NOT EXISTS user_word (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            last_answered TEXT,
            delay INTEGER,
            active INTEGER NOT NULL,
            word_id INTEGER NOT NULL,
            user_id INTEGER NOT NULL,
            FOREIGN KEY (word_id) REFERENCES word(id),
            FOREIGN KEY (user_id) REFERENCES user(id)
        );
    `

	createUserIndex = `
        CREATE INDEX user_word_id_index ON user_word (user_id);
    `

	createUserWord = `
        INSERT INTO user_word (last_answered, delay, active, word_id, user_id) VALUES (@lastAnswered, @delay, @active, @wordId, @userId);
    `

	getWordsForUserForDictionary = `
        SELECT user_word.*, word.*
        FROM user_word, word
        WHERE user_word.word_id = word.id
        AND user_word.user_id = @userId
        AND word.dictionary_id = @dictionaryId;
    `

	getViableWordsForUserForDictionary = `
        SELECT user_word.*, word.*
        FROM user_word, word
        WHERE user_word.word_id = word.id
        AND user_word.user_id = @userId
        AND word.dictionary_id = @dictionaryId
        AND strftime('%s', 'now') - strftime('%s ', SUBSTR(last_answered, 1, 19)) > delay * 24 * 60 * 60;
    `

	getViableWordsForUserForDictionaryWhereItIsntActiveQuestion = `
        SELECT user_word.*, word.*
        FROM user_word, word
        WHERE user_word.word_id = word.id
        AND user_word.user_id = @userId
        AND word.dictionary_id = @dictionaryId
        AND word.id NOT IN (
            SELECT word_id FROM active_question WHERE user_id = @userId
        )
        AND strftime('%s', 'now') - strftime('%s ', SUBSTR(last_answered, 1, 19)) > delay * 24 * 60 * 60;
    `
	// if sent delay is not 0, then increase delay by 1
	setNewDelayForUser = `
        UPDATE user_word
        SET delay = CASE
			WHEN @delay = 0 THEN delay = 0
			ELSE delay + 1
			END
		WHERE user_id = @userId
		AND word_id = @wordId;
		
    `

	deactivateWordForUser = `
        UPDATE user_word
        SET active = 0
        WHERE user_id = @userId
        AND word_id = @wordId;
    `

	getUserWordByUserID = `
        SELECT * FROM user_word WHERE user_id = @userId;
    `

	getDelayForWordForUser = `
        SELECT delay FROM user_word WHERE user_id = @userId AND word_id = @wordId;
    `

	updateLastAnswered = `
        UPDATE user_word
        SET last_answered = @lastAnswered
        WHERE user_id = @userId
        AND word_id = @wordId;
    `

	deleteUserWordbyID = `
        DELETE FROM user_word WHERE word_id = @wordId;
    `
)

type UserWord struct {
	ID           int    `json:"id"`
	LastAnswered string `json:"lastAnswered"`
	Delay        int    `json:"delay"`
	Active       int    `json:"active"`
	WordID       int    `json:"wordId"`
	UserID       int    `json:"userId"`
}

type Word worddb.Word

func CreateUserWordTable(dbase *sql.DB) error {
	_, err := dbase.Exec(createUserWordTable)
	return err
}

func CreateUserWord(dbase *sql.DB, userWord *UserWord) error {
	_, err := dbase.Exec(createUserWord, userWord.LastAnswered, userWord.Delay, userWord.Active, userWord.WordID, userWord.UserID)
	return err
}

func CreateIndexForUserWord(dbase *sql.DB) error {
	_, err := dbase.Exec(createUserIndex)
	return err
}

func GetWordsForUserForDictionary(dbase *sql.DB, userID, dictionaryID int) ([]UserWord, []Word, error) {
	rows, err := dbase.Query(getWordsForUserForDictionary, userID, dictionaryID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var userWords []UserWord
	var words []Word

	for rows.Next() {
		var userWord UserWord
		var word Word

		if err := rows.Scan(
			&userWord.ID, &userWord.LastAnswered, &userWord.Delay, &userWord.Active, &userWord.WordID, &userWord.UserID,
			&word.ID, &word.ForeignWord, &word.ForeignDescription, &word.NativeWord, &word.NativeDescription, &word.Pronunciation, &word.DictionaryID,
		); err != nil {
			return nil, nil, err
		}

		userWords = append(userWords, userWord)
		words = append(words, word)
	}

	return userWords, words, nil
}

func GetViableWordsForUserForDictionary(dbase *sql.DB, userID, dictionaryID int) ([]UserWord, []Word, error) {
	rows, err := dbase.Query(getViableWordsForUserForDictionary, userID, dictionaryID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var userWords []UserWord
	var words []Word

	for rows.Next() {
		var userWord UserWord
		var word Word

		if err := rows.Scan(
			&userWord.ID, &userWord.LastAnswered, &userWord.Delay, &userWord.Active, &userWord.WordID, &userWord.UserID,
			&word.ID, &word.ForeignWord, &word.ForeignDescription, &word.NativeWord, &word.NativeDescription, &word.Pronunciation, &word.DictionaryID,
		); err != nil {
			return nil, nil, err
		}

		userWords = append(userWords, userWord)
		words = append(words, word)
	}

	return userWords, words, nil
}

func GetViableWordsForUserForDictionaryWhereItIsntActiveQuestion(dbase *sql.DB, userID, dictionaryID int) ([]UserWord, []Word, error) {
	rows, err := dbase.Query(getViableWordsForUserForDictionaryWhereItIsntActiveQuestion, userID, dictionaryID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var userWords []UserWord
	var words []Word

	for rows.Next() {
		var userWord UserWord
		var word Word

		if err := rows.Scan(
			&userWord.ID, &userWord.LastAnswered, &userWord.Delay, &userWord.Active, &userWord.WordID, &userWord.UserID,
			&word.ID, &word.ForeignWord, &word.ForeignDescription, &word.NativeWord, &word.NativeDescription, &word.Pronunciation, &word.DictionaryID,
		); err != nil {
			return nil, nil, err
		}

		userWords = append(userWords, userWord)
		words = append(words, word)
	}

	return userWords, words, nil
}

func SetNewDelayForUser(dbase *sql.DB, userID, wordID, delay int) error {
	_, err := dbase.Exec(setNewDelayForUser, delay, userID, wordID)
	return err
}

func DeactivateWordForUser(dbase *sql.DB, userID, wordID int) error {
	_, err := dbase.Exec(deactivateWordForUser, userID, wordID)
	return err
}

func GetUserWordsByUserID(dbase *sql.DB, userID int) ([]UserWord, error) {
	rows, err := dbase.Query(getUserWordByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userWords []UserWord
	for rows.Next() {
		var userWord UserWord
		if err := rows.Scan(
			&userWord.ID, &userWord.LastAnswered, &userWord.Delay, &userWord.Active, &userWord.WordID, &userWord.UserID,
		); err != nil {
			return nil, err
		}

		userWords = append(userWords, userWord)
	}

	return userWords, nil
}

func GetDelayForWordForUser(dbase *sql.DB, userID, wordID int) (int, error) {
	var delay int
	err := dbase.QueryRow(getDelayForWordForUser, userID, wordID).Scan(&delay)
	return delay, err
}

func UpdateLastAnswered(dbase *sql.DB, userID, wordID int, lastAnswered string) error {
	_, err := dbase.Exec(updateLastAnswered, lastAnswered, userID, wordID)
	return err
}

func DeleteUserWordByID(dbase *sql.DB, wordID int) error {
	_, err := dbase.Exec(deleteUserWordbyID, wordID)
	return err
}
