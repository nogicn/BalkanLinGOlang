package wordcontroller

import (
	"BalkanLinGO/db"
	"BalkanLinGO/models/activequestiondb"
	"BalkanLinGO/models/userworddb"
	"BalkanLinGO/models/worddb"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LearnSession(c *fiber.Ctx) error {
	id := c.Params("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}
	activequestion, err := activequestiondb.GetActiveQuestionByUserID(db.DB, c.Locals("user_id").(int))
	//fmt.Println(activequestion)
	//fmt.Println(err)

	// check if activequestion exists
	if activequestion == (activequestiondb.ActiveQuestion{}) {
		//fmt.Println("No active question")
		words, err := worddb.GetAllWordsNotInUserWord(db.DB, c.Locals("user_id").(int))

		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja riječi!", "link": "/dashboard"})
		}
		//fmt.Println("rijeci", words)
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
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod stvaranja riječi!", "link": "/dashboard"})
			}
		}

		// check if activequestion exists

		activeword, err := worddb.GetWordByID(db.DB, activequestion.WordID)
		if !strings.Contains(err.Error(), "no rows in result set") {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja riječi!", "link": "/dashboard"})

		}

		if activeword.DictionaryID != idInt {
			err := activequestiondb.DeleteActiveQuestionByUserID(db.DB, c.Locals("user_id").(int))
			if err != nil {
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod brisanja aktivne riječi!", "link": "/dashboard"})
			}
			LearnSessionForeignNative(c)
		}
		switch activequestion.Type {
		case 1:
			LearnSessionForeignNative(c)
		case 2:
			//LearnSessionNativeForeign(c)
		case 3:
			//LearnSessionWriting(c)
		case 4:
			//LearnSessionPronunciation(c)
		default:
			//LearnSessionNativeForeign(c)

		}
	} else {
		LearnSessionForeignNative(c)
	}
	return nil
}

func LearnSessionForeignNative(c *fiber.Ctx) error {
	activequestion, err := activequestiondb.GetActiveQuestionByUserID(db.DB, c.Locals("user_id").(int))
	//fmt.Println(activequestion, err)
	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result") {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja aktivne riječi activequestion", "link": "/dashboard"})
		}
	}

	dictid := c.Params("id")
	dictidInt, err := strconv.Atoi(dictid)
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}
	if activequestion == (activequestiondb.ActiveQuestion{}) {
		userWords, words, err := userworddb.GetViableWordsForUserForDictionary(db.DB, c.Locals("user_id").(int), dictidInt)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja riječi! nema rijeci za dict", "link": "/dashboard"})
		}
		//fmt.Println("rijeci", words)
		//fmt.Println("userWords", userWords)
		if len(userWords) < 4 {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Nema više riječi za učenje!", "link": "/dashboard"})
		}
		finalWords := []worddb.Word{}
		// get 4 random words
		for i := 0; i < 4; i++ {
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
		// set random word as active question
		random := rand.Intn(len(finalWords))
		activequestion.WordID = finalWords[random].ID
		activequestion.UserID = c.Locals("user_id").(int)
		activequestion.Type = 1
		err = activequestiondb.SetActiveQuestion(db.DB, &activequestion)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod stvaranja aktivne riječi!", "link": "/dashboard"})
		}
		// swap values of foreign and native word from userWord
		for i := 0; i < len(finalWords); i++ {
			if finalWords[i].ID == userWords[random].WordID {
				finalWords[i].ForeignWord, finalWords[i].NativeWord = finalWords[i].NativeWord, finalWords[i].ForeignWord
			}
		}
		return c.Render("learnSession", fiber.Map{"words": finalWords, "dictionaryId": dictidInt, "currentWord": finalWords[random]})

	} else {
		// do the same thing but without setting the active queston and only randoming 3 words
		userWords, words, err := userworddb.GetViableWordsForUserForDictionary(db.DB, c.Locals("user_id").(int), dictidInt)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja riječi! nema rijeci za dict", "link": "/dashboard"})
		}
		////fmt.Println("rijeci", words)
		//fmt.Println("userWords", userWords)
		if len(userWords) < 4 {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Nema više riječi za učenje!", "link": "/dashboard"})
		}
		finalWords := []worddb.Word{}
		// add active question to final words
		activeword, err := worddb.GetWordByID(db.DB, activequestion.WordID)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja aktivne riječi!", "link": "/dashboard"})
		}
		finalWords = append(finalWords, activeword)
		// get 4 random words
		for i := 0; i < 3; i++ {
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
		var random int
		// swap values of foreign and native word from userWord
		for i := 0; i < len(finalWords); i++ {
			if finalWords[i].ID == userWords[random].WordID {
				finalWords[i].ForeignWord, finalWords[i].NativeWord = finalWords[i].NativeWord, finalWords[i].ForeignWord
			}
		}
		// randomize word order in finalWords
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(finalWords), func(i, j int) { finalWords[i], finalWords[j] = finalWords[j], finalWords[i] })
		return c.Render("learnSession", fiber.Map{"words": finalWords, "dictionaryId": dictidInt, "currentWord": activeword})

	}

}
