package routes

import (
	//"BalkanLinGO/controllers"

	learningcontroller "BalkanLinGO/controllers/learning"
	"BalkanLinGO/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// create fiber route for users

func IndexRouter(app *fiber.App, session *session.Store) {

	route := app.Group("/")

	route.Get("/register", func(c *fiber.Ctx) error {
		return c.Render("auth/register", fiber.Map{"title": "Register"})
	})

	route.Get("/login",
	func(c *fiber.Ctx) error {
		return c.Render("auth/login", fiber.Map{"title": "Login"})
	})

	route.Get("/reset", func(c *fiber.Ctx) error {
		return c.Render("auth/resetPass", fiber.Map{"title": "Reset"})
	})

	route.Get("/dict", func(c *fiber.Ctx) error {
		return c.Render("dictSearch", fiber.Map{"title": "Dictionary"})
	})

	route.Get("/error", func(c *fiber.Ctx) error {
		return c.Render("forOfor", fiber.Map{"status": "404", "errorText": c.Locals("errorText"), "link": "/", "auth": c.Locals("user_id") != nil})
	})
	
	route.Get("/",
	func(c *fiber.Ctx) error {
		return middleware.CheckAuth(c, session)
	},
	func(c *fiber.Ctx) error {
		return c.Render("dashboard", fiber.Map{"title": "Dashboard", "hx_link": "dictionary", "IsAdmin": c.Locals("is_admin")})
	})

	route.Get("/learnSession/:id",
	func(c *fiber.Ctx) error {
		return middleware.CheckAuth(c, session)
	},
	func(c *fiber.Ctx) error {
		return middleware.HtmxMiddleware(c)
	},
	func(c *fiber.Ctx) error {
		return learningcontroller.LearnSession(c)
	})
}
