package routes

import (
	localecontroller "BalkanLinGO/controllers/locale"
	"BalkanLinGO/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func LocaleRouter(app *fiber.App, session *session.Store) {

	route := app.Group("/locale")
	route.Use(func(c *fiber.Ctx) error { return middleware.CheckAuth(c, session) }, middleware.IsAdmin(session))

	route.Get("/adminLocales", localecontroller.AdminLocales)
	route.Get("/editLocale/:id", localecontroller.EditLocale)
	route.Get("/deleteLocale/:id", localecontroller.DeleteLocale)
	route.Get("/addLocale", localecontroller.AddLocale)
	route.Post("/saveLocale", localecontroller.SaveLocale)

}
