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
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/dashboard"})
		}

	} else {
		dictionaries, err = dictionarydb.GetAllDictionariesWithIcons(db.DB)
		if err != nil {
			fmt.Println(err)
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/dashboard"})
		}

	}
	return c.Render("dashboard", fiber.Map{"dictionaries": dictionaries, "IsAdmin": c.Locals("is_admin")})

}

func AddDictionary(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(int)
	if isAdmin == 0 {
		// get all dictionaries not assigned to user
		id := c.Locals("user_id").(int)
		dictionaries, err := dictionarydb.GetDictionariesNotAssignedToUser(db.DB, id)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/dashboard"})
		}
		return c.Render("dictionary/addDictionary", fiber.Map{"dictionaries": dictionaries, "IsAdmin": c.Locals("is_admin")})
	} else {
		// get all dictionaries
		dictionaries, err := dictionarydb.GetAllDictionaries(db.DB)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/dashboard"})
		}
		languages, err := languagedb.GetAllLanguages(db.DB)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju jezika!", "link": "/dashboard"})
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
		return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dodavanju rečnika!", "link": "/dashboard"})
	}
	err = dictionaryuserdb.AddDictionaryToUser(db.DB, id, dictIDInt)
	if err != nil {
		return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dodavanju rečnika!", "link": "/dashboard"})
	}
	return c.Redirect("/dashboard")
}

func AdminEditDict(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(int)
	if isAdmin == 0 {
		return c.Render("forOfor", fiber.Map{"status": "401", "errorText": "Nemate pristup!", "link": "/dashboard"})
	} else {
		id := c.Params("id")
		// convert to int
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/dashboard"})
		}

		dict, err := dictionarydb.GetDictionaryByID(db.DB, idInt)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/dashboard"})
		}

		languages, err := languagedb.GetAllLanguages(db.DB)

		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju jezika!", "link": "/dashboard"})
		}
		return c.Render("dictionary/dictionaryAddAdmin", fiber.Map{"dictionary": dict, "IsAdmin": c.Locals("is_admin"), "languages": languages})
	}
}

func AdminSaveDict(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(int)
	if isAdmin == 0 {
		return c.Render("forOfor", fiber.Map{"status": "401", "errorText": "Nemate pristup!", "link": "/dashboard"})
	} else {

		// get data from form
		description := c.FormValue("description")
		imageLink := c.FormValue("imageLink")
		// convert to int

		langID, err := strconv.Atoi(c.FormValue("langId"))
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/dashboard"})
		}
		id := c.FormValue("id")
		if id == "" {
			err = dictionarydb.CreateNewDictionary(db.DB, &dictionarydb.Dictionary{Name: description, LanguageID: langID, ImageLink: imageLink})
			if err != nil {
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri kreiranju rečnika!", "link": "/dashboard"})
			}
		} else {
			idInt, err := strconv.Atoi(id)
			if err != nil {
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/dashboard"})
			}
			err = dictionarydb.UpdateDictionary(db.DB, &dictionarydb.Dictionary{ID: idInt, Name: description, LanguageID: langID, ImageLink: imageLink})
			if err != nil {
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri kreiranju rečnika!", "link": "/dashboard"})
			}
		}
		return c.Redirect("/dashboard")
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
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/dashboard"})
		}
		dictionaryuserdb.DeleteDictionaryFromUser(db.DB, c.Locals("user_id").(int), dictIDInt)

	} else {
		id := c.Params("id")
		// convert to int
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/dashboard"})
		}

		allwords, err := worddb.GetWordsByDictionaryID(db.DB, idInt)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri brisanju rečnika!", "link": "/dashboard"})
		}
		for _, word := range allwords {
			err = activequestiondb.DeleteActiveQuestionByWordID(db.DB, word.ID)
			if err != nil {
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri brisanju rečnika!", "link": "/dashboard"})
			}
		}
		err = dictionarydb.DeleteDictionary(db.DB, idInt)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri brisanju rečnika!", "link": "/dashboard"})
		}

	}
	return c.Redirect("/dashboard")
}

func SearchDictionary(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(int)
	id := c.Params("id")
	// convert to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/dashboard"})
	}

	if isAdmin == 0 {
		return c.Render("forOfor", fiber.Map{"status": "401", "errorText": "Nemate pristup!", "link": "/dashboard"})
	} else {
		dictionaries, err := dictionarydb.GetDictionaryByID(db.DB, idInt)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/dashboard"})
		}
		words, err := worddb.GetWordsByDictionaryID(db.DB, dictionaries.ID)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju reči!", "link": "/dashboard"})
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
		return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/dashboard"})
	}

	if isAdmin == 0 {
		return c.Render("forOfor", fiber.Map{"status": "401", "errorText": "Nemate pristup!", "link": "/dashboard"})
	} else {
		words, err := worddb.SearchWordByDictionaryID(db.DB, idInt, word)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/dashboard"})
		}
		return c.Render("word/partials/wordsList", fiber.Map{"words": words})

	}
}
