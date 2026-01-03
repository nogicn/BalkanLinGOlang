package dictionarycontroller

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"BalkanLinGO/internal/db"
	dbr "BalkanLinGO/internal/db/repository"

	"github.com/gofiber/fiber/v2"
)

type DictionaryController struct {
	repo *dbr.Queries
}

func New(dbService db.Service) *DictionaryController {
	return &DictionaryController{
		repo: dbService.GetRepository(),
	}
}

func (dc *DictionaryController) Dashboard(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(bool)
	var dictionaries interface{}
	var err error
	ctx := context.Background()

	if isAdmin == false {
		id := c.Locals("user_id").(int64)
		dictionaries, err = dc.repo.GetDictionariesForUser(ctx, int64(id))
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/"})
		}

	} else {
		dictionaries, err = dc.repo.GetAllDictionariesWithIcons(ctx)
		if err != nil {
			fmt.Println(err)
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/"})
		}

	}
	c.Response().Header.Set("HX-Replace-Url", "/dictionary")
	return c.Render("dictionary/dictionaryShow", fiber.Map{"dictionaries": dictionaries, "IsAdmin": c.Locals("is_admin")})

}

func (dc *DictionaryController) AddDictionary(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(bool)
	ctx := context.Background()

	if isAdmin == false {
		// get all dictionaries not assigned to user
		id := c.Locals("user_id").(int64)
		dictionaries, err := dc.repo.GetDictionariesNotAssignedToUser(ctx, int64(id))
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/"})
		}
		return c.Render("dictionary/addDictionary", fiber.Map{"dictionaries": dictionaries, "IsAdmin": c.Locals("is_admin")})
	} else {
		// get all dictionaries
		dictionaries, err := dc.repo.GetAllDictionaries(ctx)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/"})
		}
		languages, err := dc.repo.GetAllLanguages(ctx)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju jezika!", "link": "/"})
		}
		return c.Render("dictionary/dictionaryAddAdmin", fiber.Map{"dictionaries": dictionaries, "IsAdmin": c.Locals("is_admin"), "languages": languages})
	}
}

func (dc *DictionaryController) AddDictionaryToUser(c *fiber.Ctx) error {
	// get user from locals
	id := c.Locals("user_id").(int64)
	dictID := c.Params("id")
	// convert to int
	dictIDInt, err := strconv.Atoi(dictID)
	if err != nil {
		return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dodavanju rečnika!", "link": "/"})
	}
	ctx := context.Background()
	err = dc.repo.AddDictionaryToUser(ctx, dbr.AddDictionaryToUserParams{
		UserID:       int64(id),
		DictionaryID: int64(dictIDInt),
	})
	if err != nil {
		return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dodavanju rečnika!", "link": "/"})
	}
	return dc.Dashboard(c)
}

func (dc *DictionaryController) AdminEditDict(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(bool)
	if isAdmin == false {
		return c.Render("forOfor", fiber.Map{"status": "401", "errorText": "Nemate pristup!", "link": "/"})
	} else {
		id := c.Params("id")
		// convert to int
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/"})
		}

		ctx := context.Background()
		dict, err := dc.repo.GetDictionaryByID(ctx, int64(idInt))
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/"})
		}

		languages, err := dc.repo.GetAllLanguages(ctx)

		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju jezika!", "link": "/"})
		}
		return c.Render("dictionary/dictionaryAddAdmin", fiber.Map{"dictionary": dict, "IsAdmin": c.Locals("is_admin"), "languages": languages})
	}
}

func (dc *DictionaryController) AdminSaveDict(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(bool)
	if isAdmin == false {
		return c.Render("forOfor", fiber.Map{"status": "401", "errorText": "Nemate pristup!", "link": "/"})
	} else {

		// get data from form
		description := c.FormValue("description")
		imageLink := c.FormValue("imageLink")
		// convert to int

		langID, err := strconv.Atoi(c.FormValue("langId"))
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/"})
		}
		id := c.FormValue("id")
		ctx := context.Background()

		if id == "" {
			err = dc.repo.CreateDictionary(ctx, dbr.CreateDictionaryParams{
				Name:       description,
				LanguageID: int64(langID),
				ImageLink:  imageLink,
			})
			if err != nil {
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri kreiranju rečnika!", "link": "/"})
			}
		} else {
			idInt, err := strconv.Atoi(id)
			if err != nil {
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/"})
			}
			err = dc.repo.UpdateDictionary(ctx, dbr.UpdateDictionaryParams{
				ID:         int64(idInt),
				Name:       description,
				LanguageID: int64(langID),
				ImageLink:  imageLink,
			})
			if err != nil {
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri kreiranju rečnika!", "link": "/"})
			}
		}
		return dc.Dashboard(c)
	}
}

func (dc *DictionaryController) RemoveDictionary(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(bool)
	ctx := context.Background()

	if isAdmin == false {
		dictID := c.Params("id")
		// convert to int
		dictIDInt, err := strconv.Atoi(dictID)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/"})
		}
		dc.repo.DeleteDictionaryFromUser(ctx, dbr.DeleteDictionaryFromUserParams{
			UserID:       int64(c.Locals("user_id").(int64)),
			DictionaryID: int64(dictIDInt),
		})

	} else {
		id := c.Params("id")
		// convert to int
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/"})
		}

		allwords, err := dc.repo.GetWordsByDictionaryID(ctx, int64(idInt))
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri brisanju rečnika!", "link": "/"})
		}
		for _, word := range allwords {
			err = dc.repo.DeleteActiveQuestionByWordID(ctx, sql.NullInt64{Int64: word.ID, Valid: true})
			if err != nil {
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri brisanju rečnika!", "link": "/"})
			}
		}
		err = dc.repo.DeleteDictionary(ctx, int64(idInt))
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri brisanju rečnika!", "link": "/"})
		}

	}
	return dc.Dashboard(c)
}

func (dc *DictionaryController) SearchDictionary(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(bool)
	id := c.Params("dictId")
	// convert to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/"})
	}

	if isAdmin == false {
		return c.Render("forOfor", fiber.Map{"status": "401", "errorText": "Nemate pristup!", "link": "/"})
	} else {
		ctx := context.Background()
		dictionaries, err := dc.repo.GetDictionaryByID(ctx, int64(idInt))
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/"})
		}
		words, err := dc.repo.GetWordsByDictionaryID(ctx, dictionaries.ID)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju reči!", "link": "/"})
		}

		return c.Render("dictionary/dictSearch", fiber.Map{"dictionary": dictionaries, "IsAdmin": c.Locals("is_admin"), "words": words})
	}
}

func (dc *DictionaryController) SearchWords(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(bool)
	id := c.Params("id")
	word := c.FormValue("word")
	// convert to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/"})
	}

	if isAdmin == false {
		return c.Render("forOfor", fiber.Map{"status": "401", "errorText": "Nemate pristup!", "link": "/"})
	} else {
		ctx := context.Background()
		words, err := dc.repo.SearchWordByDictionaryID(ctx, dbr.SearchWordByDictionaryIDParams{
			DictionaryID: int64(idInt),
			SearchTerm:   sql.NullString{String: word, Valid: true},
		})
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/"})
		}
		return c.Render("word/partials/wordsList", fiber.Map{"words": words})

	}
}
