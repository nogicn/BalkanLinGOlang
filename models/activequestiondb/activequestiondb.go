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
        DELETE FROM active_question WHERE user_id = @userId;
    `

	deleteActiveQuestionWordId = `
        DELETE FROM active_question WHERE word_id = @wordId;
    `

	setActiveQuestion = `
    INSERT OR REPLACE INTO active_question (user_id, word_id, type) 
    VALUES (@userId, @wordId, @type);
    
    `

	getActiveQuestion = `
        SELECT * FROM active_question WHERE user_id = @userId;
    `

	increaseActiveQuestionType = `
        UPDATE active_question
        SET type = CASE
            WHEN type >= 1 AND type < 4 THEN type + 1
            ELSE 1 
            END
        WHERE user_id = @userId;
    `
)

func CreateActiveQuestionTable(db *sql.DB) error {
	_, err := db.Exec(createActiveQuestionTable)
	return err
}

func DeleteActiveQuestionByUserID(db *sql.DB, userId int) error {
	_, err := db.Exec(deleteActiveQuestion, userId)
	return err
}

func DeleteActiveQuestionByWordID(db *sql.DB, wordId int) error {
	_, err := db.Exec(deleteActiveQuestionWordId, wordId)
	return err
}

func SetActiveQuestion(db *sql.DB, activeQuestion *ActiveQuestion) error {
	_, err := db.Exec(setActiveQuestion, activeQuestion.UserID, activeQuestion.WordID, activeQuestion.Type)
	return err
}

func GetActiveQuestionByUserID(db *sql.DB, userId int) (ActiveQuestion, error) {
	var activeQuestion ActiveQuestion
	err := db.QueryRow(getActiveQuestion, userId).Scan(&activeQuestion.ID, &activeQuestion.UserID, &activeQuestion.WordID, &activeQuestion.Type)
	return activeQuestion, err
}

func IncreaseActiveQuestionTypeByUserID(db *sql.DB, userId int) error {
	_, err := db.Exec(increaseActiveQuestionType, userId)
	return err
}
