package server

import (
	//"BalkanLinGO/controllers"

	learningcontroller "BalkanLinGO/internal/server/controllers/learning"
	"BalkanLinGO/internal/server/middleware"

	"github.com/gofiber/fiber/v2"
)

// create fiber route for users

func (s *FiberServer) RegisterIndexRouter() {
	lc := learningcontroller.New(s.db)

	route := s.Group("/")

	route.Get("/register", func(c *fiber.Ctx) error {
		return c.Render("auth/authBase", fiber.Map{"title": "Register", "register": true})
	})

	route.Get("/login",
		func(c *fiber.Ctx) error {
			return c.Render("auth/authBase", fiber.Map{"title": "Login", "login": true})
		})

	route.Get("/reset", func(c *fiber.Ctx) error {
		return c.Render("auth/authBase", fiber.Map{"title": "Reset", "reset": true})
	})

	route.Get("/error", func(c *fiber.Ctx) error {
		return c.Render("forOfor", fiber.Map{"status": "404", "errorText": c.Locals("errorText"), "link": "/", "auth": c.Locals("user_id") != nil})
	})

	route.Get("/",
		func(c *fiber.Ctx) error {
			return middleware.CheckAuth(c, s.session, s.db)
		},
		func(c *fiber.Ctx) error {
			return c.Render("dashboard", fiber.Map{"title": "Dashboard", "hx_link": "dictionary", "IsAdmin": c.Locals("is_admin")})
		})

	route.Get("/learnSession/:id",
		func(c *fiber.Ctx) error {
			return middleware.CheckAuth(c, s.session, s.db)
		},
		func(c *fiber.Ctx) error {
			return middleware.HtmxMiddleware(c)
		},
		func(c *fiber.Ctx) error {
			return lc.LearnSession(c)
		})
}
