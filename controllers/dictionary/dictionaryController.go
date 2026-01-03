package dictionarycontroller

import (
	"BalkanLinGO/db"
	"BalkanLinGO/models/activequestiondb"
	"BalkanLinGO/models/dictionarydb"
	"BalkanLinGO/models/dictionaryuserdb"
	"BalkanLinGO/models/languagedb"
	"BalkanLinGO/models/worddb"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Dashboard(c *fiber.Ctx) error {
	// get user from locals

	isAdmin := c.Locals("is_admin").(int)
	var dictionaries []dictionarydb.Dictionary
	var err error
	if isAdmin == 0 {
		id := c.Locals("user_id").(int)
		dictionaries, err = dictionarydb.GetDictionariesForUser(db.DB, id)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/"})
		}

	} else {
		dictionaries, err = dictionarydb.GetAllDictionariesWithIcons(db.DB)
		if err != nil {
			fmt.Println(err)
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/"})
		}

	}
	c.Response().Header.Set("HX-Replace-Url", "/dictionary")
	return c.Render("dictionary/dictionaryShow", fiber.Map{"dictionaries": dictionaries, "IsAdmin": c.Locals("is_admin")})

}

func AddDictionary(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(int)
	if isAdmin == 0 {
		// get all dictionaries not assigned to user
		id := c.Locals("user_id").(int)
		dictionaries, err := dictionarydb.GetDictionariesNotAssignedToUser(db.DB, id)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/"})
		}
		return c.Render("dictionary/addDictionary", fiber.Map{"dictionaries": dictionaries, "IsAdmin": c.Locals("is_admin")})
	} else {
		// get all dictionaries
		dictionaries, err := dictionarydb.GetAllDictionaries(db.DB)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/"})
		}
		languages, err := languagedb.GetAllLanguages(db.DB)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju jezika!", "link": "/"})
		}
		return c.Render("dictionary/dictionaryAddAdmin", fiber.Map{"dictionaries": dictionaries, "IsAdmin": c.Locals("is_admin"), "languages": languages})
	}
}

func AddDictionaryToUser(c *fiber.Ctx) error {
	// get user from locals
	id := c.Locals("user_id").(int)
	dictID := c.Params("id")
	// convert to int
	dictIDInt, err := strconv.Atoi(dictID)
	if err != nil {
		return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dodavanju rečnika!", "link": "/"})
	}
	err = dictionaryuserdb.AddDictionaryToUser(db.DB, id, dictIDInt)
	if err != nil {
		return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dodavanju rečnika!", "link": "/"})
	}
	return Dashboard(c)
}

func AdminEditDict(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(int)
	if isAdmin == 0 {
		return c.Render("forOfor", fiber.Map{"status": "401", "errorText": "Nemate pristup!", "link": "/"})
	} else {
		id := c.Params("id")
		// convert to int
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/"})
		}

		dict, err := dictionarydb.GetDictionaryByID(db.DB, idInt)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/"})
		}

		languages, err := languagedb.GetAllLanguages(db.DB)

		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju jezika!", "link": "/"})
		}
		return c.Render("dictionary/dictionaryAddAdmin", fiber.Map{"dictionary": dict, "IsAdmin": c.Locals("is_admin"), "languages": languages})
	}
}

func AdminSaveDict(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(int)
	if isAdmin == 0 {
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
		if id == "" {
			err = dictionarydb.CreateNewDictionary(db.DB, &dictionarydb.Dictionary{Name: description, LanguageID: langID, ImageLink: imageLink})
			if err != nil {
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri kreiranju rečnika!", "link": "/"})
			}
		} else {
			idInt, err := strconv.Atoi(id)
			if err != nil {
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/"})
			}
			err = dictionarydb.UpdateDictionary(db.DB, &dictionarydb.Dictionary{ID: idInt, Name: description, LanguageID: langID, ImageLink: imageLink})
			if err != nil {
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri kreiranju rečnika!", "link": "/"})
			}
		}
		return Dashboard(c)
	}
}

func RemoveDictionary(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(int)
	if isAdmin == 0 {
		dictID := c.Params("id")
		// convert to int
		dictIDInt, err := strconv.Atoi(dictID)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/"})
		}
		dictionaryuserdb.DeleteDictionaryFromUser(db.DB, c.Locals("user_id").(int), dictIDInt)

	} else {
		id := c.Params("id")
		// convert to int
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/"})
		}

		allwords, err := worddb.GetWordsByDictionaryID(db.DB, idInt)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri brisanju rečnika!", "link": "/"})
		}
		for _, word := range allwords {
			err = activequestiondb.DeleteActiveQuestionByWordID(db.DB, word.ID)
			if err != nil {
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri brisanju rečnika!", "link": "/"})
			}
		}
		err = dictionarydb.DeleteDictionary(db.DB, idInt)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri brisanju rečnika!", "link": "/"})
		}

	}
	return Dashboard(c)
}

func SearchDictionary(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(int)
	id := c.Params("dictId")
	// convert to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/"})
	}

	if isAdmin == 0 {
		return c.Render("forOfor", fiber.Map{"status": "401", "errorText": "Nemate pristup!", "link": "/"})
	} else {
		dictionaries, err := dictionarydb.GetDictionaryByID(db.DB, idInt)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/"})
		}
		words, err := worddb.GetWordsByDictionaryID(db.DB, dictionaries.ID)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju reči!", "link": "/"})
		}

		return c.Render("dictionary/dictSearch", fiber.Map{"dictionary": dictionaries, "IsAdmin": c.Locals("is_admin"), "words": words})
	}
}

func SearchWords(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(int)
	id := c.Params("id")
	word := c.FormValue("word")
	// convert to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/"})
	}

	if isAdmin == 0 {
		return c.Render("forOfor", fiber.Map{"status": "401", "errorText": "Nemate pristup!", "link": "/"})
	} else {
		words, err := worddb.SearchWordByDictionaryID(db.DB, idInt, word)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/"})
		}
		return c.Render("word/partials/wordsList", fiber.Map{"words": words})

	}
}
