package locale_controller

import (
	"BalkanLinGO/db"
	"BalkanLinGO/models/languagedb"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

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
		return c.Redirect("/locale/adminLocales")
	}
}

func DeleteLocale(c *fiber.Ctx) error {
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
		err = languagedb.DeleteLanguageById(db.DB, idInt)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri brisanju jezika!", "link": "/dashboard"})

		}
		return c.Redirect("/locale/adminLocales")
	}
}

func AddLocale(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(int)
	if isAdmin == 0 {
		return c.Render("forOfor", fiber.Map{"status": "401", "errorText": "Nemate pristup!", "link": "/dashboard"})

	} else {
		return c.Render("localeAddAdmin", fiber.Map{"IsAdmin": c.Locals("is_admin")})
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
