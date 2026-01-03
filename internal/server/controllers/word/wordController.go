package wordcontroller

import (
	"BalkanLinGO/internal/db"
	dbr "BalkanLinGO/internal/db/repository"
	"BalkanLinGO/internal/server/middleware"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type WordController struct {
	repo *dbr.Queries
}

func New(dbService db.Service) *WordController {
	return &WordController{
		repo: dbService.GetRepository(),
	}
}

func (wc *WordController) EditWord(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(bool)
	if isAdmin == false {
		return c.Render("forOfor", fiber.Map{"status": "401", "errorText": "Nemate pristup!", "link": "/dashboard"})
	} else {
		id := c.Params("id")
		// convert to int
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/dashboard"})
		}
		word, err := wc.repo.GetWordByID(c.Context(), int64(idInt))
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/dashboard"})
		}
		dictionary, err := wc.repo.GetDictionaryByID(c.Context(), word.DictionaryID)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/dashboard"})
		}

		return c.Render("word/editWord", fiber.Map{"word": word, "IsAdmin": c.Locals("is_admin"), "dictionary": dictionary, "wordactionurl": "editWord"})
	}
}

func (wc *WordController) SaveWord(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(bool)
	if isAdmin == false {
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
			err = wc.repo.CreateWord(c.Context(), dbr.CreateWordParams{
				ForeignWord:        foreignWord,
				ForeignDescription: foreignDescription,
				NativeWord:         nativeWord,
				NativeDescription:  nativeDescription,
				Pronunciation:      pronunciation,
				DictionaryID:       int64(dictIDInt),
			})
			if err != nil {
				fmt.Println(err)
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri kreiranju reči!", "link": "/dashboard"})
			}
		} else {
			idInt, err := strconv.Atoi(id)
			if err != nil {
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/dashboard"})
			}
			err = wc.repo.UpdateWord(c.Context(), dbr.UpdateWordParams{
				ID:                 int64(idInt),
				ForeignWord:        foreignWord,
				ForeignDescription: foreignDescription,
				NativeWord:         nativeWord,
				NativeDescription:  nativeDescription,
				Pronunciation:      pronunciation,
			})
			if err != nil {
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri kreiranju reči!", "link": "/dashboard"})
			}
		}
		return c.Redirect("/dictionary/dictSearch/" + dictID)
	}
}

func (wc *WordController) AddWord(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(bool)
	if isAdmin == false {
		return c.Render("forOfor", fiber.Map{"status": "401", "errorText": "Nemate pristup!", "link": "/dashboard"})
	} else {
		dictID := c.Params("id")
		// convert to int
		dictIDInt, err := strconv.Atoi(dictID)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/dashboard"})
		}
		dictionary, err := wc.repo.GetDictionaryByID(c.Context(), int64(dictIDInt))
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/dashboard"})
		}

		return c.Render("word/editWord", fiber.Map{"IsAdmin": c.Locals("is_admin"), "dictionary": dictionary, "wordactionurl": "addWord"})
	}
}

func (wc *WordController) DeleteWord(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(bool)
	if isAdmin == false {
		return c.Render("forOfor", fiber.Map{"status": "401", "errorText": "Nemate pristup!", "link": "/dashboard"})
	} else {
		id := c.Params("wordId")
		// convert to int
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/dashboard"})
		}
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/dashboard"})
		}

		err = wc.repo.DeleteWordByID(c.Context(), int64(idInt))
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri brisanju reči!", "link": "/dashboard"})
		}

		dictId := c.Params("dictId")
		return c.Redirect("/dictionary/dictSearch/" + dictId)
	}
}

func (wc *WordController) CreatePronunciation(c *fiber.Ctx) error {
	id := c.FormValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/"})
	}
	// get all data from form as word
	dictionaryID := c.Params("id")
	// convert to int
	dictionaryIDInt, err := strconv.Atoi(dictionaryID)
	if err != nil {
		return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/"})
	}
	word := dbr.Word{
		ID:                 int64(idInt),
		Foreignword:        c.FormValue("foreignWord"),
		Foreigndescription: c.FormValue("foreignDescription"),
		Nativeword:         c.FormValue("nativeWord"),
		Nativedescription:  c.FormValue("nativeDescription"),
		Pronunciation:      c.FormValue("pronunciation"),
	}

	dictionary, err := wc.repo.GetDictionaryByID(c.Context(), int64(dictionaryIDInt))
	if err != nil {
		return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/"})
	}

	filename := randStringBytes(32) + ".mp3"
	err = middleware.GenerateSpeech(word.Foreignword, filename)
	if err != nil {
		return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri generisanju izgovora!", "link": "/"})
	}

	word.Pronunciation = filename
	return c.Render("word/partials/wordsEdit", fiber.Map{"word": word, "dictionary": dictionary})
}
