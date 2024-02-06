package dictionarydb

import "database/sql"

const (
	createDictionaryTable = `
        CREATE TABLE IF NOT EXISTS dictionary (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            language_id INTEGER NOT NULL,
            image_link TEXT NOT NULL,
            FOREIGN KEY (language_id) REFERENCES language(id)
        );
    `

	createNewDictionary = `
        INSERT INTO dictionary (name, language_id, image_link) VALUES (@name, @language_id, @imageLink);
    `

	getDictionariesForUser = `
        SELECT dictionary.*, language.flag_icon
        FROM dictionary
        LEFT JOIN language ON dictionary.language_id = language.id
        LEFT JOIN dictionary_user ON dictionary.id = dictionary_user.dictionary_id
        WHERE dictionary_user.dictionary_id = dictionary.id
        AND dictionary_user.user_id = @userID;
    `

	getAllDictionaries = `
        SELECT *
        FROM dictionary;
    `

	getAllDictionariesWithIcons = `
        SELECT dictionary.*, language.flag_icon
        FROM dictionary 
        LEFT JOIN language ON dictionary.language_id = language.id;
    `

	deleteDictionary = `
        DELETE FROM dictionary
        WHERE dictionary.id = @id;
    `

	getDictionariesNotAssignedToUser = `
        SELECT dictionary.*
        FROM dictionary
        WHERE dictionary.id NOT IN (
            SELECT dictionary_id
            FROM dictionary_user
            WHERE user_id = @userID
        );
    `

	getDictionaryByID = `
        SELECT *
        FROM dictionary
        WHERE id = @id;
    `

	updateDictionary = `
        UPDATE dictionary
        SET name = @name, language_id = @language_id, image_link = @imageLink
        WHERE id = @id;
    `
)

type Dictionary struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	LanguageID int    `json:"languageId"`
	ImageLink  string `json:"imageLink"`
	FlagIcon   string `json:"flagIcon"`
}

func CreateDictionaryTable(db *sql.DB) error {
	_, err := db.Exec(createDictionaryTable)
	return err
}

func CreateNewDictionary(db *sql.DB, dictionary *Dictionary) error {
	_, err := db.Exec(createNewDictionary, sql.Named("name", dictionary.Name), sql.Named("language_id", dictionary.LanguageID), sql.Named("imageLink", dictionary.ImageLink))
	return err
}

func GetDictionariesForUser(db *sql.DB, userID int) ([]Dictionary, error) {
	rows, err := db.Query(getDictionariesForUser, sql.Named("userID", userID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dictionaries []Dictionary
	for rows.Next() {
		var dictionary Dictionary
		if err := rows.Scan(&dictionary.ID, &dictionary.Name, &dictionary.LanguageID, &dictionary.ImageLink, &dictionary.FlagIcon); err != nil {
			return nil, err
		}
		dictionaries = append(dictionaries, dictionary)
	}

	return dictionaries, nil
}

func GetAllDictionaries(db *sql.DB) ([]Dictionary, error) {
	rows, err := db.Query(getAllDictionaries)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dictionaries []Dictionary
	for rows.Next() {
		var dictionary Dictionary
		if err := rows.Scan(&dictionary.ID, &dictionary.Name, &dictionary.LanguageID, &dictionary.ImageLink); err != nil {
			return nil, err
		}
		dictionaries = append(dictionaries, dictionary)
	}

	return dictionaries, nil
}

func GetAllDictionariesWithIcons(db *sql.DB) ([]Dictionary, error) {
	rows, err := db.Query(getAllDictionariesWithIcons)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dictionaries []Dictionary
	for rows.Next() {
		var dictionary Dictionary
		if err := rows.Scan(&dictionary.ID, &dictionary.Name, &dictionary.LanguageID, &dictionary.ImageLink, &dictionary.FlagIcon); err != nil {
			return nil, err
		}
		dictionaries = append(dictionaries, dictionary)
	}

	return dictionaries, nil
}

func DeleteDictionary(db *sql.DB, id int) error {
	_, err := db.Exec(deleteDictionary, sql.Named("id", id))
	return err
}

func GetDictionariesNotAssignedToUser(db *sql.DB, userID int) ([]Dictionary, error) {
	rows, err := db.Query(getDictionariesNotAssignedToUser, sql.Named("userID", userID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dictionaries []Dictionary
	for rows.Next() {
		var dictionary Dictionary
		if err := rows.Scan(&dictionary.ID, &dictionary.Name, &dictionary.LanguageID, &dictionary.ImageLink); err != nil {
			return nil, err
		}
		dictionaries = append(dictionaries, dictionary)
	}

	return dictionaries, nil
}

func GetDictionaryByID(db *sql.DB, id int) (Dictionary, error) {
	var dictionary Dictionary
	err := db.QueryRow(getDictionaryByID, sql.Named("id", id)).Scan(&dictionary.ID, &dictionary.Name, &dictionary.LanguageID, &dictionary.ImageLink)
	return dictionary, err
}

func UpdateDictionary(db *sql.DB, dictionary *Dictionary) error {
	_, err := db.Exec(updateDictionary, sql.Named("name", dictionary.Name), sql.Named("language_id", dictionary.LanguageID), sql.Named("imageLink", dictionary.ImageLink), sql.Named("id", dictionary.ID))
	return err
}
