package localecontroller

import (
	"BalkanLinGO/internal/db"
	dbr "BalkanLinGO/internal/db/repository"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type LocaleController struct {
	repo *dbr.Queries
}

func New(dbService db.Service) *LocaleController {
	return &LocaleController{
		repo: dbService.GetRepositoryRW(),
	}
}

func (lc *LocaleController) SaveLocale(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(bool)
	if isAdmin == false {
		return c.Render("forOfor", fiber.Map{"status": "401", "errorText": "Nemate pristup!", "link": "/dashboard"})
	} else {

		// get data from form
		name := c.FormValue("name")
		shorthand := c.FormValue("shorthand")
		flagIcon := c.FormValue("flagIcon")
		id := c.FormValue("id")
		if id == "" {
			err := lc.repo.CreateLanguage(c.Context(), dbr.CreateLanguageParams{Name: name, Shorthand: shorthand, FlagIcon: flagIcon})
			if err != nil {
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri kreiranju jezika!", "link": "/dashboard"})
			}
		} else {
			idInt, err := strconv.Atoi(id)
			if err != nil {
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/dashboard"})
			}
			err = lc.repo.UpdateLanguage(c.Context(), dbr.UpdateLanguageParams{ID: int64(idInt), Name: name, Shorthand: shorthand, FlagIcon: flagIcon})
			if err != nil {
				return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri kreiranju jezika!", "link": "/dashboard"})
			}
		}
		return c.Redirect("/locale/adminLocales")
	}
}

func (lc *LocaleController) DeleteLocale(c *fiber.Ctx) error {
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
		err = lc.repo.DeleteLanguageByID(c.Context(), int64(idInt))
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri brisanju jezika!", "link": "/dashboard"})

		}
		languages, err := lc.repo.GetAllLanguages(c.Context())
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/dashboard"})
		}
		return c.Render("locale/showLocales", fiber.Map{"IsAdmin": c.Locals("is_admin"), "languages": languages})
	}
}

func (lc *LocaleController) AddLocale(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(bool)
	if isAdmin == false {
		return c.Render("forOfor", fiber.Map{"status": "401", "errorText": "Nemate pristup!", "link": "/dashboard"})

	} else {
		return c.Render("locale/localeAddAdmin", fiber.Map{"IsAdmin": c.Locals("is_admin")})
	}
}

func (lc *LocaleController) AdminLocales(c *fiber.Ctx) error {
	// get user from locals
	isAdmin := c.Locals("is_admin").(bool)
	if isAdmin == false {
		return c.Render("forOfor", fiber.Map{"status": "401", "errorText": "Nemate pristup!", "link": "/dashboard"})
	} else {
		languages, err := lc.repo.GetAllLanguages(c.Context())
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/dashboard"})
		}
		return c.Render("locale/showLocales", fiber.Map{"languages": languages, "IsAdmin": c.Locals("is_admin")})
	}
}

func (lc *LocaleController) EditLocale(c *fiber.Ctx) error {
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
		language, err := lc.repo.GetLanguageByID(c.Context(), int64(idInt))
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška, nije int!", "link": "/dashboard"})
		}

		return c.Render("locale/localeAddAdmin", fiber.Map{"locale": language, "IsAdmin": c.Locals("is_admin")})
	}
}
