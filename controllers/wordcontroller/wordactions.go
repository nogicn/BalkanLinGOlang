package wordcontroller

import (
	"BalkanLinGO/db"
	"BalkanLinGO/models/activequestiondb"
	"BalkanLinGO/models/userworddb"
	"BalkanLinGO/models/worddb"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
)

func fillWordList(words []userworddb.Word, finalWords []worddb.Word, n int) []worddb.Word {
	// get 4 random words
	for i := 0; i < n; i++ {
		////fmt.Println("generating ")
		// generate random number between 1 and number of words in dictionary
		random := rand.Intn(len(words))
		duplicate := false
		for j := 0; j < len(finalWords); j++ {
			if finalWords[j].ID == words[random].ID {
				duplicate = true
				break
			}
		}
		if duplicate {
			i--
			continue
		}
		finalWords = append(finalWords, worddb.Word(words[random]))
	}
	return finalWords
}

func createWords(c *fiber.Ctx, idInt int) error {
	words, err := worddb.GetAllWordsNotInUserWord(db.DB, c.Locals("user_id").(int))
	if err != nil {
		return fiber.NewError(500, "Database error")
	}

	for _, word := range words {
		if word.DictionaryID != idInt {
			continue
		}
		currentDate := time.Now().Add(-time.Hour * 24 * 30)
		var trueword = userworddb.UserWord{
			UserID:       c.Locals("user_id").(int),
			WordID:       word.ID,
			Active:       1,
			Delay:        0,
			LastAnswered: currentDate.Format("2006-01-02 15:04:05"),
		}
		err := userworddb.CreateUserWord(db.DB, &trueword)
		if err != nil {
			return fiber.NewError(500, "Greška kod stvaranja riječi!")
		}
	}
	return nil
}

func setActiveQuestion(activequestion *activequestiondb.ActiveQuestion, c *fiber.Ctx, idInt int, typeOf int) error {
	tmpActive := activequestiondb.ActiveQuestion{}

	if *activequestion != (activequestiondb.ActiveQuestion{}) {
		err := activequestiondb.DeleteActiveQuestionByUserID(db.DB, c.Locals("user_id").(int))
		if err != nil {
			return fiber.NewError(500, "Greška kod brisanja aktivne riječi!")
		}
	}

	userWords, _, err := userworddb.GetViableWordsForUserForDictionary(db.DB, c.Locals("user_id").(int), idInt)
	if err != nil {
		return fiber.NewError(500, "Greška kod dohvaćanja riječi! nema rijeci za dict")
	}

	if len(userWords) < 3 {
		return fiber.NewError(404, "Nema više riječi za učenje!")
	}

	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(len(userWords))
	tmpActive.WordID = userWords[random].WordID
	tmpActive.UserID = c.Locals("user_id").(int)
	if typeOf > 4 {
		typeOf = 1
	}
	tmpActive.Type = typeOf
	err = activequestiondb.SetActiveQuestion(db.DB, &tmpActive)
	if err != nil {

		return fiber.NewError(404, "Greška kod stvaranja aktivne riječi!")
	}
	tmpActive, _ = activequestiondb.GetActiveQuestionByUserID(db.DB, c.Locals("user_id").(int))

	*activequestion = tmpActive
	return nil
}
