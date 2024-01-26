package wordcontroller

import (
	"BalkanLinGO/db"
	"BalkanLinGO/models/activequestiondb"
	"BalkanLinGO/models/userworddb"
	"BalkanLinGO/models/worddb"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func createWords(c *fiber.Ctx, idInt int) error {
	words, err := worddb.GetAllWordsNotInUserWord(db.DB, c.Locals("user_id").(int))
	if err != nil {
		return fiber.NewError(500, "Database error")
	}
	////fmt.Println("rijeci", words)
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

	// delete old active question
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
	//fmt.Println("userWords", userWords)
	// set random word as active question
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
		//fmt.Println("AAAAAAAAAAAAAAA")
		return fiber.NewError(404, "Greška kod stvaranja aktivne riječi!")
	}
	tmpActive, _ = activequestiondb.GetActiveQuestionByUserID(db.DB, c.Locals("user_id").(int))
	//fmt.Println("tmpActive\n", tmpActive)
	*activequestion = tmpActive
	return nil
}

func LearnSession(c *fiber.Ctx) error {
	id := c.Params("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}

	activequestion, activeerr := activequestiondb.GetActiveQuestionByUserID(db.DB, c.Locals("user_id").(int))
	if activequestion == (activequestiondb.ActiveQuestion{}) {
		////fmt.Println("No active question")
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
		//fmt.Println("1")
		LearnSessionForeignNative(c)
	case 2:
		//fmt.Println("2")
		LearnSessionNativeForeign(c)
	case 3:
		//fmt.Println("3")
		LearnSessionWriting(c)
	case 4:
		//fmt.Println("4")
		LearnSessionPronunciation(c)
	default:
		LearnSessionForeignNative(c)

	}

	return nil
}

func LearnSessionForeignNative(c *fiber.Ctx) error {
	activequestion, err := activequestiondb.GetActiveQuestionByUserID(db.DB, c.Locals("user_id").(int))
	////fmt.Println(activequestion, err, "activequestion")
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
	// do the same thing but without setting the active queston and only randoming 3 words
	_, words, err := userworddb.GetViableWordsForUserForDictionary(db.DB, c.Locals("user_id").(int), dictidInt)
	if err != nil {
		return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja riječi! nema rijeci za dict", "link": "/dashboard"})
	}

	finalWords := []worddb.Word{}
	// add active question to final words
	activeword, err := worddb.GetWordByID(db.DB, activequestion.WordID)
	if err != nil {
		return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja aktivne riječi!", "link": "/dashboard"})
	}
	finalWords = append(finalWords, activeword)
	finalWords = fillWordList(words, finalWords, 3)
	// swap values of foreign and native word from userWord
	////fmt.Println("rijeci", finalWords)
	for i := 0; i < len(finalWords); i++ {
		finalWords[i].ForeignWord, finalWords[i].NativeWord = finalWords[i].NativeWord, finalWords[i].ForeignWord
		finalWords[i].ForeignDescription, finalWords[i].NativeDescription = finalWords[i].NativeDescription, finalWords[i].ForeignDescription
	}
	////fmt.Println("rijeci", finalWords)
	// randomize word order in finalWords
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(finalWords), func(i, j int) { finalWords[i], finalWords[j] = finalWords[j], finalWords[i] })
	activeword.ForeignWord, activeword.NativeWord = activeword.NativeWord, activeword.ForeignWord
	activeword.ForeignDescription, activeword.NativeDescription = activeword.NativeDescription, activeword.ForeignDescription
	return c.Render("learnSession", fiber.Map{"words": finalWords, "dictionaryId": dictidInt, "currentWord": activeword, "next": 2})

}

func LearnSessionNativeForeign(c *fiber.Ctx) error {
	activequestion, err := activequestiondb.GetActiveQuestionByUserID(db.DB, c.Locals("user_id").(int))
	////fmt.Println(activequestion, err)
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

	// do the same thing but without setting the active queston and only randoming 3 words
	userWords, words, err := userworddb.GetViableWordsForUserForDictionary(db.DB, c.Locals("user_id").(int), dictidInt)
	if err != nil {
		return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja riječi! nema rijeci za dict", "link": "/dashboard"})
	}
	//////fmt.Println("rijeci", words)
	////fmt.Println("userWords", userWords)
	if len(userWords) < 4 {
		return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Nema više riječi za učenje!", "link": "/dashboard"})
	}
	finalWords := []worddb.Word{}
	// add active question to final words
	activeword, err := worddb.GetWordByID(db.DB, activequestion.WordID)
	if err != nil {
		return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja aktivne riječi!", "link": "/dashboard"})
	}
	finalWords = append(finalWords, activeword)
	finalWords = fillWordList(words, finalWords, 3)
	// randomize word order in finalWords
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(finalWords), func(i, j int) { finalWords[i], finalWords[j] = finalWords[j], finalWords[i] })
	return c.Render("learnSession", fiber.Map{"words": finalWords, "dictionaryId": dictidInt, "currentWord": activeword, "next": 3})

}

func LearnSessionWriting(c *fiber.Ctx) error {

	activequestion, err := activequestiondb.GetActiveQuestionByUserID(db.DB, c.Locals("user_id").(int))
	////fmt.Println(activequestion, err)
	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result") {
			return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja aktivne riječi activequestion", "link": "/dashboard"})
		}
	}
	id := c.Params("id")
	dictidInt, err := strconv.Atoi(id)
	fmt.Println("ASDASDASDAS", dictidInt)
	if err != nil {
		return c.Status(500).SendString("Invalid ID")
	}

	// get ative word
	activeword, err := worddb.GetWordByID(db.DB, activequestion.WordID)
	if err != nil {
		return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja aktivne riječi!", "link": "/dashboard"})
	}
	return c.Render("writeWord", fiber.Map{"word": activeword, "dictionaryId": dictidInt, "next": 4})

}

func LearnSessionPronunciation(c *fiber.Ctx) error {
	activequestion, err := activequestiondb.GetActiveQuestionByUserID(db.DB, c.Locals("user_id").(int))
	////fmt.Println(activequestion, err)
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
	// get ative word
	activeword, err := worddb.GetWordByID(db.DB, activequestion.WordID)
	if err != nil {
		return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja aktivne riječi!", "link": "/dashboard"})
	}
	return c.Render("sayWord", fiber.Map{"word": activeword, "dictionaryId": dictidInt, "next": 1})

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
		// set current word as correct
		err = userworddb.SetNewDelayForUser(db.DB, c.Locals("user_id").(int), activequestion.WordID, 1)
		if err != nil {
			return c.Status(500).SendString("Greška kod postavljanja nove riječi!")
		}
		correct = true
	} else {
		// set current word as correct
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
	return c.Render("partials/word", fiber.Map{"word": activeword, "correct": correct})

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
		// set current word as correct
		err = userworddb.SetNewDelayForUser(db.DB, c.Locals("user_id").(int), activequestion.WordID, 1)
		if err != nil {
			return c.Status(500).SendString("Greška kod postavljanja nove riječi!")
		}
		correct = true
	} else {
		// set current word as correct
		err = userworddb.SetNewDelayForUser(db.DB, c.Locals("user_id").(int), activequestion.WordID, 0)
		if err != nil {
			return c.Status(500).SendString("Greška kod postavljanja nove riječi!")
		}
		correct = false
	}
	setActiveQuestion(&activequestion, c, activeword.DictionaryID, activequestion.Type+1)
	return c.Render("partials/word", fiber.Map{"word": activeword, "correct": correct})

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

	// random a valute 0 to 100
	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(100)
	var correct bool
	if random > 50 {
		// set current word as correct
		err = userworddb.SetNewDelayForUser(db.DB, c.Locals("user_id").(int), activequestion.WordID, 1)
		if err != nil {
			return c.Status(500).SendString("Greška kod postavljanja nove riječi!")
		}
		correct = true
	} else {
		// set current word as correct
		err = userworddb.SetNewDelayForUser(db.DB, c.Locals("user_id").(int), activequestion.WordID, 0)
		if err != nil {
			return c.Status(500).SendString("Greška kod postavljanja nove riječi!")
		}
		correct = false
	}
	setActiveQuestion(&activequestion, c, activeword.DictionaryID, activequestion.Type+1)
	return c.Render("partials/word", fiber.Map{"word": activeword, "correct": correct})

}
