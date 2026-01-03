package activequestiondb

import "database/sql"

const (
	createActiveQuestionTable = `
        CREATE TABLE IF NOT EXISTS active_question (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER ,
            word_id INTEGER,
            type INTEGER NOT NULL DEFAULT 1, -- (1, 2, 3)
            FOREIGN KEY (user_id) REFERENCES user(id),
            FOREIGN KEY (word_id) REFERENCES word(id)
            UNIQUE (user_id, word_id)
        );
    `

	deleteActiveQuestion = `
        DELETE FROM active_question WHERE user_id = @userID;
    `

	deleteActiveQuestionWordID = `
        DELETE FROM active_question WHERE word_id = @wordId;
    `

	setActiveQuestion = `
    INSERT OR REPLACE INTO active_question (user_id, word_id, type) 
    VALUES (@userID, @wordId, @type);
    
    `

	getActiveQuestion = `
        SELECT * FROM active_question WHERE user_id = @userID;
    `

	increaseActiveQuestionType = `
        UPDATE active_question
        SET type = CASE
            WHEN type >= 1 AND type < 4 THEN type + 1
            ELSE 1 
            END
        WHERE user_id = @userID;
    `
)

type ActiveQuestion struct {
	ID     int `json:"id"`
	UserID int `json:"userId"`
	WordID int `json:"wordId"`
	Type   int `json:"type"`
}

func CreateActiveQuestionTable(db *sql.DB) error {
	_, err := db.Exec(createActiveQuestionTable)
	return err
}

func DeleteActiveQuestionByUserID(db *sql.DB, userID int) error {
	_, err := db.Exec(deleteActiveQuestion, sql.Named("userID", userID))
	return err
}

func DeleteActiveQuestionByWordID(db *sql.DB, wordID int) error {
	_, err := db.Exec(deleteActiveQuestionWordID, sql.Named("wordId", wordID))
	return err
}

func SetActiveQuestion(db *sql.DB, activeQuestion *ActiveQuestion) error {
	_, err := db.Exec(setActiveQuestion, sql.Named("userID", activeQuestion.UserID), sql.Named("wordId", activeQuestion.WordID), sql.Named("type", activeQuestion.Type))
	return err
}

func GetActiveQuestionByUserID(db *sql.DB, userID int) (ActiveQuestion, error) {
	var activeQuestion ActiveQuestion
	err := db.QueryRow(getActiveQuestion, sql.Named("userID", userID)).Scan(&activeQuestion.ID, &activeQuestion.UserID, &activeQuestion.WordID, &activeQuestion.Type)
	return activeQuestion, err
}

func IncreaseActiveQuestionTypeByUserID(db *sql.DB, userID int) error {
	_, err := db.Exec(increaseActiveQuestionType, sql.Named("userID", userID))
	return err
}
