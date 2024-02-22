package learningcontroller

import (
	"BalkanLinGO/db"
	"BalkanLinGO/models/activequestiondb"
	"BalkanLinGO/models/userworddb"
	"BalkanLinGO/models/worddb"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func LearnSession(c *fiber.Ctx) error {
	id := c.Params("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}

	activequestion, activeerr := activequestiondb.GetActiveQuestionByUserID(db.DB, c.Locals("user_id").(int))
	if activequestion == (activequestiondb.ActiveQuestion{}) {

		err := createWords(c, idInt)
		if err != nil {
			return c.Status(404).Render("forOfor", fiber.Map{"status": "500", "errorText": err, "link": "/dashboard"})
		}
		err = setActiveQuestion(&activequestion, c, idInt, 1)
		if err != nil {
			return c.Status(404).Render("forOfor", fiber.Map{"status": "500", "errorText": err, "link": "/dashboard"})
		}
	}

	if activeerr != nil {
		if !strings.Contains(activeerr.Error(), "no rows in result") {
			return c.Status(404).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja aktivne riječi activequestion", "link": "/dashboard"})
		}
	}

	activeword, err := worddb.GetWordByID(db.DB, activequestion.WordID)
	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(404).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja riječi!", "link": "/dashboard"})
		}
	}
	if activeword.DictionaryID != idInt {
		err := activequestiondb.DeleteActiveQuestionByUserID(db.DB, c.Locals("user_id").(int))
		if err != nil {
			return c.Status(404).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod brisanja aktivne riječi!", "link": "/dashboard"})
		}
		err = createWords(c, idInt)
		if err != nil {
			return c.Status(404).Render("forOfor", fiber.Map{"status": "500", "errorText": err, "link": "/dashboard"})
		}
		err = setActiveQuestion(&activequestion, c, idInt, 1)
		if err != nil {
			return c.Status(404).Render("forOfor", fiber.Map{"status": "500", "errorText": err, "link": "/dashboard"})
		}
	}

	switch activequestion.Type {
	case 1:
		LearnSessionForeignNative(c)

	case 2:
		LearnSessionNativeForeign(c)

	case 3:
		LearnSessionWriting(c)

	case 4:
		LearnSessionPronunciation(c)

	default:
		LearnSessionForeignNative(c)

	}

	return nil
}

func LearnSessionForeignNative(c *fiber.Ctx) error {
	activequestion, err := activequestiondb.GetActiveQuestionByUserID(db.DB, c.Locals("user_id").(int))

	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result") {
			return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja aktivne riječi activequestion", "link": "/dashboard"})
		}
	}

	dictid := c.Params("id")
	dictidInt, err := strconv.Atoi(dictid)
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}

	_, words, err := userworddb.GetViableWordsForUserForDictionary(db.DB, c.Locals("user_id").(int), dictidInt)
	if err != nil {
		return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja riječi! nema rijeci za dict", "link": "/dashboard"})
	}

	finalWords := []worddb.Word{}

	activeword, err := worddb.GetWordByID(db.DB, activequestion.WordID)
	if err != nil {
		return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja aktivne riječi!", "link": "/dashboard"})
	}
	finalWords = append(finalWords, activeword)
	finalWords = fillWordList(words, finalWords, 3)

	for i := 0; i < len(finalWords); i++ {
		finalWords[i].ForeignWord, finalWords[i].NativeWord = finalWords[i].NativeWord, finalWords[i].ForeignWord
		finalWords[i].ForeignDescription, finalWords[i].NativeDescription = finalWords[i].NativeDescription, finalWords[i].ForeignDescription
	}

	rand.Shuffle(len(finalWords), func(i, j int) { finalWords[i], finalWords[j] = finalWords[j], finalWords[i] })
	activeword.ForeignWord, activeword.NativeWord = activeword.NativeWord, activeword.ForeignWord
	activeword.ForeignDescription, activeword.NativeDescription = activeword.NativeDescription, activeword.ForeignDescription
	return c.Render("learn/selectWord", fiber.Map{"words": finalWords, "dictionaryId": dictidInt, "currentWord": activeword, "next": 2})

}

func LearnSessionNativeForeign(c *fiber.Ctx) error {
	activequestion, err := activequestiondb.GetActiveQuestionByUserID(db.DB, c.Locals("user_id").(int))

	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result") {
			return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja aktivne riječi activequestion", "link": "/dashboard"})
		}
	}

	dictid := c.Params("id")
	dictidInt, err := strconv.Atoi(dictid)
	if err != nil {
		return c.Status(404).Status(400).SendString("Invalid ID")
	}

	userWords, words, err := userworddb.GetViableWordsForUserForDictionary(db.DB, c.Locals("user_id").(int), dictidInt)
	if err != nil {
		return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja riječi! nema rijeci za dict", "link": "/dashboard"})
	}

	if len(userWords) < 4 {
		return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Nema više riječi za učenje!", "link": "/dashboard"})
	}
	finalWords := []worddb.Word{}

	activeword, err := worddb.GetWordByID(db.DB, activequestion.WordID)
	if err != nil {
		return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja aktivne riječi!", "link": "/dashboard"})
	}
	finalWords = append(finalWords, activeword)
	finalWords = fillWordList(words, finalWords, 3)

	rand.Shuffle(len(finalWords), func(i, j int) { finalWords[i], finalWords[j] = finalWords[j], finalWords[i] })
	return c.Render("learn/selectWord", fiber.Map{"words": finalWords, "dictionaryId": dictidInt, "currentWord": activeword, "next": 3})

}

func LearnSessionWriting(c *fiber.Ctx) error {

	activequestion, err := activequestiondb.GetActiveQuestionByUserID(db.DB, c.Locals("user_id").(int))

	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result") {
			return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja aktivne riječi activequestion", "link": "/dashboard"})
		}
	}
	id := c.Params("id")
	dictidInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(500).SendString("Invalid ID")
	}

	activeword, err := worddb.GetWordByID(db.DB, activequestion.WordID)
	if err != nil {
		return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja aktivne riječi!", "link": "/dashboard"})
	}
	return c.Render("learn/writeWord", fiber.Map{"word": activeword, "dictionaryId": dictidInt, "next": 4})

}

func LearnSessionPronunciation(c *fiber.Ctx) error {
	activequestion, err := activequestiondb.GetActiveQuestionByUserID(db.DB, c.Locals("user_id").(int))

	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result") {
			return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja aktivne riječi activequestion", "link": "/dashboard"})
		}
	}
	id := c.Params("id")
	dictidInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(500).SendString("Invalid ID")
	}

	activeword, err := worddb.GetWordByID(db.DB, activequestion.WordID)
	if err != nil {
		return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja aktivne riječi!", "link": "/dashboard"})
	}
	return c.Render("learn/sayWord", fiber.Map{"word": activeword, "dictionaryId": dictidInt, "next": 1})

}

func CheckAnswer(c *fiber.Ctx) error {
	answer := c.Params("answer")
	answer = strings.ToLower(answer)

	activequestion, err := activequestiondb.GetActiveQuestionByUserID(db.DB, c.Locals("user_id").(int))
	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result") {
			return c.Status(500).SendString("Greška kod dohvaćanja aktivne riječi activequestion")
		}
	}
	if activequestion == (activequestiondb.ActiveQuestion{}) {
		return c.Redirect("/dashboard")
	}

	activeword, err := worddb.GetWordByID(db.DB, activequestion.WordID)
	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(500).SendString("Greška kod dohvaćanja riječi!")
		}
	}

	answerInt, err := strconv.Atoi(answer)
	if err != nil {
		return c.Status(500).SendString("Invalid ID")
	}

	if activequestion.Type == 1 {
		activeword.ForeignWord, activeword.NativeWord = activeword.NativeWord, activeword.ForeignWord
		activeword.ForeignDescription, activeword.NativeDescription = activeword.NativeDescription, activeword.ForeignDescription
	}
	var correct bool

	if answerInt == activequestion.WordID {

		err = userworddb.SetNewDelayForUser(db.DB, c.Locals("user_id").(int), activequestion.WordID, 1)
		if err != nil {
			return c.Status(500).SendString("Greška kod postavljanja nove riječi!")
		}
		correct = true
	} else {

		err = userworddb.SetNewDelayForUser(db.DB, c.Locals("user_id").(int), activequestion.WordID, 0)
		if err != nil {
			return c.Status(500).SendString("Greška kod postavljanja nove riječi!")
		}
		activeword, err = worddb.GetWordByID(db.DB, answerInt)
		if err != nil {
			if !strings.Contains(err.Error(), "no rows in result set") {
				return c.Status(500).SendString("Greška kod dohvaćanja riječi!")
			}
		}
		correct = false
	}
	setActiveQuestion(&activequestion, c, activeword.DictionaryID, activequestion.Type+1)
	return c.Render("learn/partials/word", fiber.Map{"word": activeword, "correct": correct})

}

func CheckWritingAnswer(c *fiber.Ctx) error {
	answer := c.FormValue("foreignWord")
	answer = strings.ToLower(answer)

	activequestion, err := activequestiondb.GetActiveQuestionByUserID(db.DB, c.Locals("user_id").(int))
	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result") {
			return c.Status(500).SendString("Greška kod dohvaćanja aktivne riječi activequestion")
		}
	}
	if activequestion == (activequestiondb.ActiveQuestion{}) {
		return c.Redirect("/dashboard")
	}

	activeword, err := worddb.GetWordByID(db.DB, activequestion.WordID)
	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(500).SendString("Greška kod dohvaćanja riječi!")
		}
	}

	var correct bool
	activeword.ForeignWord = strings.ToLower(activeword.ForeignWord)
	fmt.Println(answer, activeword.ForeignWord)
	if answer == activeword.ForeignWord {

		err = userworddb.SetNewDelayForUser(db.DB, c.Locals("user_id").(int), activequestion.WordID, 1)
		if err != nil {
			return c.Status(500).SendString("Greška kod postavljanja nove riječi!")
		}
		correct = true
	} else {

		err = userworddb.SetNewDelayForUser(db.DB, c.Locals("user_id").(int), activequestion.WordID, 0)
		if err != nil {
			return c.Status(500).SendString("Greška kod postavljanja nove riječi!")
		}
		correct = false
	}
	setActiveQuestion(&activequestion, c, activeword.DictionaryID, activequestion.Type+1)
	return c.Render("learn/partials/writeWordAnswer", fiber.Map{"word": activeword, "correct": correct})

}

func CheckListeningAnswer(c *fiber.Ctx) error {
	activequestion, err := activequestiondb.GetActiveQuestionByUserID(db.DB, c.Locals("user_id").(int))
	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result") {
			return c.Status(500).SendString("Greška kod dohvaćanja aktivne riječi activequestion")
		}
	}
	if activequestion == (activequestiondb.ActiveQuestion{}) {
		return c.Redirect("/dashboard")
	}

	activeword, err := worddb.GetWordByID(db.DB, activequestion.WordID)
	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(500).SendString("Greška kod dohvaćanja riječi!")
		}
	}

	random := rand.Intn(100)
	var correct bool
	if random > 50 {

		err = userworddb.SetNewDelayForUser(db.DB, c.Locals("user_id").(int), activequestion.WordID, 1)
		if err != nil {
			return c.Status(500).SendString("Greška kod postavljanja nove riječi!")
		}
		correct = true
	} else {

		err = userworddb.SetNewDelayForUser(db.DB, c.Locals("user_id").(int), activequestion.WordID, 0)
		if err != nil {
			return c.Status(500).SendString("Greška kod postavljanja nove riječi!")
		}
		correct = false
	}
	setActiveQuestion(&activequestion, c, activeword.DictionaryID, activequestion.Type+1)
	return c.Render("learn/partials/sayWordAnswer", fiber.Map{"word": activeword, "correct": correct})

}
