package wordcontroller

import (
	"BalkanLinGO/db"
	"BalkanLinGO/models/activequestiondb"
	"BalkanLinGO/models/userworddb"
	"BalkanLinGO/models/worddb"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LearnSession(c *fiber.Ctx) error {
	id := c.Params("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}
	activequestion, err := activequestiondb.GetActiveQuestionByUserID(db.DB, c.Locals("id").(int))
	fmt.Println(activequestion)
	fmt.Println(err)

	if err != nil {
		fmt.Println("No active question")
		userWord, words, err := userworddb.GetViableWordsForUserForDictionary(db.DB, c.Locals("id").(int), idInt)
		fmt.Println(userWord)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja riječi!", "link": "/login"})
		}
		fmt.Println(words)
		for _, word := range words {
			currentDate := time.Now()
			var trueword = userworddb.UserWord{
				UserID:       c.Locals("id").(int),
				WordID:       word.ID,
				Active:       1,
				Delay:        0,
				LastAnswered: currentDate.Format("2006-01-02 15:04:05"),
			}
			err := userworddb.CreateUserWord(db.DB, &trueword)
			if err != nil {
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod stvaranja riječi!", "link": "/login"})
			}
		}

		// check if activequestion exists

		activeword, err := worddb.GetWordByID(db.DB, activequestion.WordID)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja riječi!", "link": "/login"})
		}

		if activeword.DictionaryID != idInt {
			err := activequestiondb.DeleteActiveQuestionByUserID(db.DB, c.Locals("user_id").(int))
			if err != nil {
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod brisanja aktivne riječi!", "link": "/login"})
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
	return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja riječi!", "link": "/login"})
}

func LearnSessionForeignNative(c *fiber.Ctx) error {
	activequestion, err := activequestiondb.GetActiveQuestionByUserID(db.DB, c.Locals("user_id").(int))
	if err != nil {
		return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja aktivne riječi!", "link": "/login"})
	}
	dictid := c.Params("id")
	dictidInt, err := strconv.Atoi(dictid)
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}
	if activequestion != (activequestiondb.ActiveQuestion{}) {
		userWords, words, err := userworddb.GetViableWordsForUserForDictionary(db.DB, c.Locals("user_id").(int), dictidInt)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja riječi!", "link": "/login"})
		}
		if len(userWords) < 4 {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Nema više riječi za učenje!", "link": "/login"})
		}
		// get 4 random words
		for i := 0; i < 4; i++ {
			// generate random number between 1 and number of words in dictionary
			random := rand.Intn(len(userWords))
			duplicate := false
			for j := 0; j < len(userWords); j++ {
				if userWords[j].ID == userWords[random].ID {
					duplicate = true
					break
				}
			}
			if duplicate {
				i--
				continue
			}
			userWords = append(userWords, userWords[random])
		}
		// set random word as active question
		random := rand.Intn(len(userWords))
		activequestion.WordID = userWords[random].ID
		activequestion.Type = 1
		err = activequestiondb.SetActiveQuestion(db.DB, &activequestion)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod stvaranja aktivne riječi!", "link": "/login"})
		}
		// swap values of foreign and native word from userWord
		for i := 0; i < len(words); i++ {
			if words[i].ID == userWords[random].WordID {
				words[i].ForeignWord, words[i].NativeWord = words[i].NativeWord, words[i].ForeignWord
			}
		}
		return c.Render("learnSessionForeignNative", fiber.Map{"words": words, "activequestion": activequestion})

	}
	return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja aktivne riječi!", "link": "/login"})
}
