package languagedb

import "database/sql"

const (
	createLanguageTable = `
        CREATE TABLE IF NOT EXISTS language (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            shorthand TEXT NOT NULL UNIQUE,
            flag_icon TEXT NOT NULL
        );
    `

	createNewLanguage = `
        INSERT INTO language (name, shorthand, flag_icon) VALUES (@name, @shorthand, @flagIcon);
    `

	getAllLanguages = `
        SELECT *
        FROM language;
    `

	getShorthands = `
        SELECT shorthand
        FROM language;
    `

	getLanguageById = `
        SELECT *
        FROM language
        WHERE id = @id;
    `

	deleteLanguageById = `
        DELETE FROM language
        WHERE id = @id;
    `

	updateLanguage = `
        UPDATE language
        SET name = @name,
            shorthand = @shorthand,
            flag_icon = @flagIcon
        WHERE id = @id;
    `
)

type Language struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Shorthand string `json:"shorthand"`
	FlagIcon  string `json:"flagIcon"`
}

func CreateLanguageTable(db *sql.DB) error {
	_, err := db.Exec(createLanguageTable)
	return err
}

func CreateLanguage(db *sql.DB, language *Language) error {
	_, err := db.Exec(createNewLanguage, language.Name, language.Shorthand, language.FlagIcon)
	return err
}

func GetAllLanguages(db *sql.DB) ([]Language, error) {
	rows, err := db.Query(getAllLanguages)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var languages []Language
	for rows.Next() {
		var language Language
		if err := rows.Scan(&language.ID, &language.Name, &language.Shorthand, &language.FlagIcon); err != nil {
			return nil, err
		}
		languages = append(languages, language)
	}

	return languages, nil
}

func GetShorthands(db *sql.DB) ([]string, error) {
	rows, err := db.Query(getShorthands)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shorthands []string
	for rows.Next() {
		var shorthand string
		if err := rows.Scan(&shorthand); err != nil {
			return nil, err
		}
		shorthands = append(shorthands, shorthand)
	}

	return shorthands, nil
}

func GetLanguageById(db *sql.DB, id int) (Language, error) {
	var language Language
	err := db.QueryRow(getLanguageById, id).Scan(&language.ID, &language.Name, &language.Shorthand, &language.FlagIcon)
	return language, err
}

func DeleteLanguageById(db *sql.DB, id int) error {
	_, err := db.Exec(deleteLanguageById, id)
	return err
}

func UpdateLanguage(db *sql.DB, language *Language) error {
	_, err := db.Exec(updateLanguage, language.Name, language.Shorthand, language.FlagIcon, language.ID)
	return err
}
