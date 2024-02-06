package dictionaryuserdb

import "database/sql"

const (
	createDictionaryUserTable = `
        CREATE TABLE IF NOT EXISTS dictionary_user (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            dictionary_id INTEGER NOT NULL,
            FOREIGN KEY (user_id) REFERENCES user(id),
            FOREIGN KEY (dictionary_id) REFERENCES dictionary(id)
        );
    `

	addDictionaryToUser = `
        INSERT INTO dictionary_user (user_id, dictionary_id) VALUES (@userId, @dictionaryId);
    `

	deleteDictionaryFromUser = `
        DELETE FROM dictionary_user WHERE user_id = @userId AND dictionary_id = @dictionaryId;
    `

	getDictionaryByID = `
        SELECT * FROM dictionary_user WHERE id = @id;
    `

	getUserDictionaries = `
        SELECT * FROM dictionary_user WHERE user_id = @userId AND dictionary_id = @dictionaryId;
    `
)

type DictionaryUser struct {
	ID           int `json:"id"`
	UserID       int `json:"userId"`
	DictionaryID int `json:"dictionaryId"`
}

func CreateDictionaryUserTable(db *sql.DB) error {
	_, err := db.Exec(createDictionaryUserTable)
	return err
}

func AddDictionaryToUser(db *sql.DB, userID, dictionaryID int) error {
	_, err := db.Exec(addDictionaryToUser, userID, dictionaryID)
	return err
}

func DeleteDictionaryFromUser(db *sql.DB, userID, dictionaryID int) error {
	_, err := db.Exec(deleteDictionaryFromUser, userID, dictionaryID)
	return err
}

func GetDictionaryUserByID(db *sql.DB, id int) (DictionaryUser, error) {
	var dictionaryUser DictionaryUser
	err := db.QueryRow(getDictionaryByID, id).Scan(&dictionaryUser.ID, &dictionaryUser.UserID, &dictionaryUser.DictionaryID)
	return dictionaryUser, err
}

func GetUserDictionaries(db *sql.DB, userID, dictionaryID int) ([]DictionaryUser, error) {
	rows, err := db.Query(getUserDictionaries, userID, dictionaryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dictionaryUsers []DictionaryUser
	for rows.Next() {
		var dictionaryUser DictionaryUser
		if err := rows.Scan(&dictionaryUser.ID, &dictionaryUser.UserID, &dictionaryUser.DictionaryID); err != nil {
			return nil, err
		}
		dictionaryUsers = append(dictionaryUsers, dictionaryUser)
	}

	return dictionaryUsers, nil
}
