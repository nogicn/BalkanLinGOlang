package dictionarycontroller

import (
	"BalkanLinGO/db"
	"BalkanLinGO/home"
	"BalkanLinGO/models/activequestiondb"
	"BalkanLinGO/models/dictionarydb"
	"BalkanLinGO/models/dictionaryuserdb"
	"BalkanLinGO/models/languagedb"
	"BalkanLinGO/models/worddb"
	"fmt"
	"strconv"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func Dashboard(c *fiber.Ctx) error {
	// get user from locals

	isAdmin := c.Locals("is_admin").(int)
	var dictionaries []dictionarydb.Dictionary
	var err error
	if isAdmin == 0 {
		id := c.Locals("id").(int)
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
		id := c.Locals("id").(int)
		dictionaries, err := dictionarydb.GetDictionariesNotAssignedToUser(db.DB, id)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/dashboard"})
		}
		return c.Render("addDictionary", fiber.Map{"dictionaries": dictionaries, "IsAdmin": c.Locals("is_admin")})
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
		return c.Render("dictionaryAddAdmin", fiber.Map{"dictionaries": dictionaries, "IsAdmin": c.Locals("is_admin"), "languages": languages})
	}
}

func AddDictionaryToUser(c *fiber.Ctx) error {
	// get user from locals
	id := c.Locals("id").(int)
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

		dict, err := dictionarydb.GetDictionaryById(db.DB, idInt)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/dashboard"})
		}

		languages, err := languagedb.GetAllLanguages(db.DB)

		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju jezika!", "link": "/dashboard"})
		}
		return c.Render("dictionaryAddAdmin", fiber.Map{"dictionary": dict, "IsAdmin": c.Locals("is_admin"), "languages": languages})
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

		langId, err := strconv.Atoi(c.FormValue("langId"))
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/dashboard"})
		}
		id := c.FormValue("id")
		if id == "" {
			err = dictionarydb.CreateNewDictionary(db.DB, &dictionarydb.Dictionary{Name: description, LanguageID: langId, ImageLink: imageLink})
			if err != nil {
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri kreiranju rečnika!", "link": "/dashboard"})
			}
		} else {
			idInt, err := strconv.Atoi(id)
			if err != nil {
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/dashboard"})
			}
			err = dictionarydb.UpdateDictionary(db.DB, &dictionarydb.Dictionary{ID: idInt, Name: description, LanguageID: langId, ImageLink: imageLink})
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
		dictionaryuserdb.DeleteDictionaryFromUser(db.DB, c.Locals("id").(int), dictIDInt)

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
		dictionaries, err := dictionarydb.GetDictionaryById(db.DB, idInt)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju rečnika!", "link": "/dashboard"})
		}
		words, err := worddb.GetWordsByDictionaryID(db.DB, dictionaries.ID)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri dohvatanju reči!", "link": "/dashboard"})
		}

		return c.Render("dictSearch", fiber.Map{"dictionary": dictionaries, "IsAdmin": c.Locals("is_admin"), "words": words})
	}
}

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

		return c.Render("editWord", fiber.Map{"word": word, "IsAdmin": c.Locals("is_admin"), "dictionary": dictionary})
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

		return c.Render("addWord", fiber.Map{"IsAdmin": c.Locals("is_admin"), "dictionary": dictionary})
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
		return c.Render("partials/wordsList", fiber.Map{"words": words})

	}
}

func SearchWords2(c *fiber.Ctx) error {
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

		handler := adaptor.HTTPHandler(templ.Handler(home.Words(words)))

		return handler(c)

	}
}

func AdminLocales(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(int)
	if isAdmin == 0 {
		return c.Render("forOfor", fiber.Map{"status": "401", "errorText": "Nemate pristup!", "link": "/dashboard"})
	} else {
		languages, err := languagedb.GetAllLanguages(db.DB)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/dashboard"})
		}
		return c.Render("dictionaryLocales", fiber.Map{"languages": languages, "IsAdmin": c.Locals("is_admin")})
	}
}

func EditLocale(c *fiber.Ctx) error {
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
		language, err := languagedb.GetLanguageById(db.DB, idInt)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/dashboard"})
		}

		return c.Render("localeAddAdmin", fiber.Map{"locale": language, "IsAdmin": c.Locals("is_admin")})
	}
}

func SaveLocale(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(int)
	if isAdmin == 0 {
		return c.Render("forOfor", fiber.Map{"status": "401", "errorText": "Nemate pristup!", "link": "/dashboard"})
	} else {

		// get data from form
		name := c.FormValue("name")
		shorthand := c.FormValue("shorthand")
		flagIcon := c.FormValue("flagIcon")
		id := c.FormValue("id")
		if id == "" {
			err := languagedb.CreateLanguage(db.DB, &languagedb.Language{Name: name, Shorthand: shorthand, FlagIcon: flagIcon})
			if err != nil {
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri kreiranju jezika!", "link": "/dashboard"})
			}
		} else {
			idInt, err := strconv.Atoi(id)
			if err != nil {
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/dashboard"})
			}
			err = languagedb.UpdateLanguage(db.DB, &languagedb.Language{ID: idInt, Name: name, Shorthand: shorthand, FlagIcon: flagIcon})
			if err != nil {
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri kreiranju jezika!", "link": "/dashboard"})
			}
		}
		return c.Redirect("/dictionary/adminLocales")
	}
}
