package word_controller

import (
	"BalkanLinGO/db"
	"BalkanLinGO/middleware"
	"BalkanLinGO/models/activequestiondb"
	"BalkanLinGO/models/dictionarydb"
	"BalkanLinGO/models/worddb"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func EditWord(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(int)
	if isAdmin == 0 {
		return c.Render("forOfor", fiber.Map{"status": "401", "errorText": "Nemate pristup!", "link": "/dashboard"})
	} else {
		id := c.Params("id")
		// convert to int
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/dashboard"})
		}
		word, err := worddb.GetWordByID(db.DB, idInt)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/dashboard"})
		}
		dictionary, err := dictionarydb.GetDictionaryById(db.DB, word.DictionaryID)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/dashboard"})
		}

		return c.Render("word/editWord", fiber.Map{"word": word, "IsAdmin": c.Locals("is_admin"), "dictionary": dictionary})
	}
}

func SaveWord(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(int)
	if isAdmin == 0 {
		return c.Render("forOfor", fiber.Map{"status": "401", "errorText": "Nemate pristup!", "link": "/dashboard"})
	} else {
		// get data from form
		foreignWord := c.FormValue("foreignWord")
		foreignDescription := c.FormValue("foreignDescription")
		nativeWord := c.FormValue("nativeWord")
		nativeDescription := c.FormValue("nativeDescription")
		pronunciation := c.FormValue("pronunciation")
		dictID := c.Params("id")
		id := c.FormValue("id")
		// convert to int
		dictIDInt, err := strconv.Atoi(dictID)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/dashboard"})
		}
		if id == "" {
			err = worddb.CreateWord(db.DB, &worddb.Word{ForeignWord: foreignWord, ForeignDescription: foreignDescription, NativeWord: nativeWord, NativeDescription: nativeDescription, Pronunciation: pronunciation, DictionaryID: dictIDInt})
			if err != nil {
				fmt.Println(err)
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri kreiranju reči!", "link": "/dashboard"})
			}
		} else {
			idInt, err := strconv.Atoi(id)
			if err != nil {
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/dashboard"})
			}
			err = worddb.UpdateWord(db.DB, &worddb.Word{ID: idInt, ForeignWord: foreignWord, ForeignDescription: foreignDescription, NativeWord: nativeWord, NativeDescription: nativeDescription, Pronunciation: pronunciation, DictionaryID: dictIDInt})
			if err != nil {
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri kreiranju reči!", "link": "/dashboard"})
			}
		}
		return c.Redirect("/dictionary/dictSearch/" + dictID)
	}
}

func AddWord(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(int)
	if isAdmin == 0 {
		return c.Render("forOfor", fiber.Map{"status": "401", "errorText": "Nemate pristup!", "link": "/dashboard"})
	} else {
		dictID := c.Params("id")
		// convert to int
		dictIDInt, err := strconv.Atoi(dictID)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/dashboard"})
		}
		dictionary, err := dictionarydb.GetDictionaryById(db.DB, dictIDInt)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/dashboard"})
		}

		return c.Render("word/addWord", fiber.Map{"IsAdmin": c.Locals("is_admin"), "dictionary": dictionary})
	}
}

func DeleteWord(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(int)
	if isAdmin == 0 {
		return c.Render("forOfor", fiber.Map{"status": "401", "errorText": "Nemate pristup!", "link": "/dashboard"})
	} else {
		id := c.Params("id")
		// convert to int
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/dashboard"})
		}
		word, err := worddb.GetWordByID(db.DB, idInt)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/dashboard"})
		}
		dictID := word.DictionaryID

		err = activequestiondb.DeleteActiveQuestionByWordID(db.DB, idInt)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri brisanju reči!", "link": "/dashboard"})
		}
		err = worddb.DeleteWordByID(db.DB, idInt)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri brisanju reči!", "link": "/dashboard"})
		}
		return c.Redirect("/dictionary/dictSearch/" + strconv.Itoa(dictID))
	}
}

func CreatePronunciation(c *fiber.Ctx) error {
	// get all data from form as word
	wordid := c.Params("id")
	// convert to int
	wordidInt, err := strconv.Atoi(wordid)
	if err != nil {
		return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/"})
	}
	word := worddb.Word{
		ID:                 wordidInt,
		ForeignWord:        c.FormValue("foreignWord"),
		ForeignDescription: c.FormValue("foreignDescription"),
		NativeWord:         c.FormValue("nativeWord"),
		NativeDescription:  c.FormValue("nativeDescription"),
		Pronunciation:      c.FormValue("pronunciation"),
	}
	id := c.FormValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/"})
	}

	dictionary, err := dictionarydb.GetDictionaryById(db.DB, idInt)
	if err != nil {
		return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/"})
	}

	filename := randStringBytes(32) + ".mp3"
	err = middleware.GenerateSpeech(word.ForeignWord, filename)
	if err != nil {
		return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri generisanju izgovora!", "link": "/"})
	}

	word.Pronunciation = filename
	return c.Render("word/partials/wordsEdit", fiber.Map{"word": word, "dictionary": dictionary})
}
