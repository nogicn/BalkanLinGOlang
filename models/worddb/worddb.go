package worddb

import "database/sql"

const (
	createWordTable = `
        CREATE TABLE IF NOT EXISTS word (
            id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
            foreignWord TEXT NOT NULL,
            foreignDescription TEXT NOT NULL,
            nativeWord TEXT NOT NULL,
            nativeDescription TEXT NOT NULL,
            pronunciation TEXT NOT NULL,
            dictionary_id INTEGER NOT NULL,
            FOREIGN KEY (dictionary_id) REFERENCES dictionary(id),
            UNIQUE (foreignWord, foreignDescription, nativeWord, nativeDescription, dictionary_id)
        );
    `

	createWord = `
        INSERT INTO word (foreignWord, foreignDescription, nativeWord, nativeDescription, pronunciation, dictionary_id) VALUES (@foreignWord, @foreignDescription, @nativeWord, @nativeDescription, @pronunciation, @dictionaryId);
    `

	deleteWordByID = `
        DELETE FROM word WHERE id = @wordId;
    `

	deleteWordByMeaning = `
        DELETE FROM word WHERE foreignWord = @foreignWord AND foreignDescription = @foreignDescription AND nativeWord = @nativeWord AND nativeDescription = @nativeDescription;
    `

	getWordByDictionaryID = `
        SELECT * FROM word WHERE dictionary_id = @dictionaryId;
    `

	deleteWordByDictionaryID = `
        DELETE FROM word WHERE dictionary_id = @dictionaryId;
    `

	getAllWords = `
        SELECT * FROM word;
    `

	getWordByID = `
        SELECT * FROM word WHERE id = @wordId;
    `

	updateWord = `
        UPDATE word SET foreignWord = @foreignWord, foreignDescription = @foreignDescription, nativeWord = @nativeWord, nativeDescription = @nativeDescription, pronunciation = @pronunciation WHERE id = @wordId;
    `

	searchWordByDictionaryID = `
        SELECT * FROM word WHERE dictionary_id = @dictionaryId AND (foreignWord LIKE '%' || @word || '%' OR nativeWord LIKE '%' || @word || '%');
    `

	getAllWordsNotInUserWord = `
        SELECT * FROM word WHERE id NOT IN (SELECT word_id as id FROM user_word WHERE user_id = @userId);
    `
)

type Word struct {
	ID                 int    `json:"id"`
	ForeignWord        string `json:"foreignWord"`
	ForeignDescription string `json:"foreignDescription"`
	NativeWord         string `json:"nativeWord"`
	NativeDescription  string `json:"nativeDescription"`
	Pronunciation      string `json:"pronunciation"`
	DictionaryID       int    `json:"dictionaryId"`
}

func CreateWordTable(db *sql.DB) error {
	_, err := db.Exec(createWordTable)
	return err
}

func CreateWord(db *sql.DB, word *Word) error {
	_, err := db.Exec(createWord, word.ForeignWord, word.ForeignDescription, word.NativeWord, word.NativeDescription, word.Pronunciation, word.DictionaryID)
	return err
}

func DeleteWordByID(db *sql.DB, wordID int) error {
	_, err := db.Exec(deleteWordByID, wordID)
	return err
}

func DeleteWordByMeaning(db *sql.DB, foreignWord, foreignDescription, nativeWord, nativeDescription string) error {
	_, err := db.Exec(deleteWordByMeaning, foreignWord, foreignDescription, nativeWord, nativeDescription)
	return err
}

func GetWordsByDictionaryID(db *sql.DB, dictionaryID int) ([]Word, error) {
	rows, err := db.Query(getWordByDictionaryID, dictionaryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []Word
	for rows.Next() {
		var word Word
		if err := rows.Scan(&word.ID, &word.ForeignWord, &word.ForeignDescription, &word.NativeWord, &word.NativeDescription, &word.Pronunciation, &word.DictionaryID); err != nil {
			return nil, err
		}
		words = append(words, word)
	}

	return words, nil
}

func DeleteWordsByDictionaryID(db *sql.DB, dictionaryID int) error {
	_, err := db.Exec(deleteWordByDictionaryID, dictionaryID)
	return err
}

func GetAllWords(db *sql.DB) ([]Word, error) {
	rows, err := db.Query(getAllWords)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []Word
	for rows.Next() {
		var word Word
		if err := rows.Scan(&word.ID, &word.ForeignWord, &word.ForeignDescription, &word.NativeWord, &word.NativeDescription, &word.Pronunciation, &word.DictionaryID); err != nil {
			return nil, err
		}
		words = append(words, word)
	}

	return words, nil
}

func GetWordByID(db *sql.DB, wordID int) (Word, error) {
	var word Word
	err := db.QueryRow(getWordByID, wordID).Scan(&word.ID, &word.ForeignWord, &word.ForeignDescription, &word.NativeWord, &word.NativeDescription, &word.Pronunciation, &word.DictionaryID)
	return word, err
}

func UpdateWord(db *sql.DB, word *Word) error {
	_, err := db.Exec(updateWord, word.ForeignWord, word.ForeignDescription, word.NativeWord, word.NativeDescription, word.Pronunciation, word.ID)
	return err
}

func SearchWordByDictionaryID(db *sql.DB, dictionaryID int, searchString string) ([]Word, error) {
	rows, err := db.Query(searchWordByDictionaryID, dictionaryID, searchString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []Word
	for rows.Next() {
		var word Word
		if err := rows.Scan(&word.ID, &word.ForeignWord, &word.ForeignDescription, &word.NativeWord, &word.NativeDescription, &word.Pronunciation, &word.DictionaryID); err != nil {
			return nil, err
		}
		words = append(words, word)
	}

	return words, nil
}

func GetAllWordsNotInUserWord(db *sql.DB, userID int) ([]Word, error) {
	rows, err := db.Query(getAllWordsNotInUserWord, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []Word
	for rows.Next() {
		var word Word
		if err := rows.Scan(&word.ID, &word.ForeignWord, &word.ForeignDescription, &word.NativeWord, &word.NativeDescription, &word.Pronunciation, &word.DictionaryID); err != nil {
			return nil, err
		}
		words = append(words, word)
	}

	return words, nil
}
