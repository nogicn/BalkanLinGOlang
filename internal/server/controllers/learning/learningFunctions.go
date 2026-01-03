package learningcontroller

import (
	dbr "BalkanLinGO/internal/db/repository"
	"context"
	"database/sql"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
)

func fillWordList(rows []dbr.GetViableWordsForUserForDictionaryRow, finalWords []dbr.Word, n int) []dbr.Word {
	// get 4 random words
	for i := 0; i < n; i++ {
		////fmt.Println("generating ")
		// generate random number between 1 and number of words in dictionary
		random := rand.Intn(len(rows))
		duplicate := false
		for j := 0; j < len(finalWords); j++ {
			if finalWords[j].ID == rows[random].Word.ID {
				duplicate = true
				break
			}
		}
		if duplicate {
			i--
			continue
		}
		finalWords = append(finalWords, rows[random].Word)
	}
	return finalWords
}

func (lc *LearnController) createWords(c *fiber.Ctx, idInt int) error {
	words, err := lc.repo.GetAllWordsNotInUserWord(c.Context(), int64(c.Locals("user_id").(int64)))
	if err != nil {
		return fiber.NewError(500, "Database error")
	}

	for _, word := range words {
		if word.DictionaryID != int64(idInt) {
			continue
		}
		currentDate := time.Now().Add(-time.Hour * 24 * 30)
		err := lc.repo.CreateUserWord(c.Context(), dbr.CreateUserWordParams{
			UserID:       int64(c.Locals("user_id").(int64)),
			WordID:       word.ID,
			Active:       1,
			Delay:        sql.NullInt64{Int64: 0, Valid: true},
			LastAnswered: sql.NullString{String: currentDate.Format("2006-01-02 15:04:05"), Valid: true},
		})
		if err != nil {
			return fiber.NewError(500, "Greška kod stvaranja riječi!")
		}
	}
	return nil
}

func (lc *LearnController) setActiveQuestion(activequestion *dbr.ActiveQuestion, c *fiber.Ctx, idInt int, typeOf int) error {
	ctx := context.Background()
	tmpActive := dbr.ActiveQuestion{}

	if *activequestion != (dbr.ActiveQuestion{}) {
		err := lc.repo.DeleteActiveQuestionByUserID(ctx, sql.NullInt64{Int64: int64(c.Locals("user_id").(int64)), Valid: true})
		if err != nil {
			return fiber.NewError(500, "Greška kod brisanja aktivne riječi!")
		}
	}

	rows, err := lc.repo.GetViableWordsForUserForDictionary(ctx, dbr.GetViableWordsForUserForDictionaryParams{
		UserID:       int64(c.Locals("user_id").(int64)),
		DictionaryID: int64(idInt),
	})
	if err != nil {
		return fiber.NewError(500, "Greška kod dohvaćanja riječi! nema rijeci za dict")
	}

	if len(rows) < 3 {
		return fiber.NewError(404, "Nema više riječi za učenje!")
	}

	random := rand.Intn(len(rows))
	if typeOf > 4 {
		typeOf = 1
	}
	err = lc.repo.SetActiveQuestion(ctx, dbr.SetActiveQuestionParams{
		UserID: sql.NullInt64{Int64: int64(c.Locals("user_id").(int64)), Valid: true},
		WordID: sql.NullInt64{Int64: rows[random].UserWord.WordID, Valid: true},
		Type:   int64(typeOf),
	})
	if err != nil {

		return fiber.NewError(404, "Greška kod stvaranja aktivne riječi!")
	}
	tmpActive, _ = lc.repo.GetActiveQuestionByUserID(ctx, sql.NullInt64{Int64: int64(c.Locals("user_id").(int64)), Valid: true})

	*activequestion = tmpActive
	return nil
}
