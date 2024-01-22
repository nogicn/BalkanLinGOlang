package routes

import (
	//"BalkanLinGO/controllers"
	"BalkanLinGO/controllers/dictionarycontroller"
	"BalkanLinGO/controllers/wordcontroller"
	"BalkanLinGO/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// create fiber route for users

func IndexRouter(app *fiber.App, session *session.Store) {

	route := app.Group("/")
	route.Get("/", func(c *fiber.Ctx) error {
		return c.Render("home", fiber.Map{"title": "Home"})
	})

	route.Get("/register", func(c *fiber.Ctx) error {
		return c.Render("register", fiber.Map{"title": "Register"})
	})

	route.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{"title": "Login"})
	})

	route.Get("/learn", func(c *fiber.Ctx) error {
		return c.Render("learnSession", fiber.Map{"title": "Learn"})
	})

	route.Get("/dict", func(c *fiber.Ctx) error {
		return c.Render("dictSearch", fiber.Map{"title": "Dictionary"})
	})

	route.Get("/dashboard", func(c *fiber.Ctx) error {
		return middleware.CheckAuth(c, session)
	}, dictionarycontroller.Dashboard)

	route.Get("/learnSession/:id", func(c *fiber.Ctx) error {
		return middleware.CheckAuth(c, session)
	}, wordcontroller.LearnSession)
}
