package server

import (
	usercontroller "BalkanLinGO/internal/server/controllers/user"
	"BalkanLinGO/internal/server/middleware"

	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) RegisterUserRouter() {
	uc := usercontroller.New(s.db)

	// public routes (no auth)
	route := s.Group("/user")

	route.Post("/login", func(c *fiber.Ctx) error { return uc.LoginUser(c, s.session) })
	route.Post("/register", uc.CreateUser)
	route.Post("/createPass", func(c *fiber.Ctx) error { return uc.CreatePass(c, s.session) })

	// protected routes: apply HTMX middleware and authentication
	protected := s.Group("/user")
	protected.Use(func(c *fiber.Ctx) error { return middleware.CheckAuth(c, s.session, s.db) })
	protected.Get("/logout", func(c *fiber.Ctx) error { return uc.LogoutUser(c, s.session) })

	//use htmx for remaining routes
	protected.Use(func(c *fiber.Ctx) error { return middleware.HtmxMiddleware(c) })

	protected.Get("/all", uc.GetUsers)
	//protected.Delete(":id", usercontroller.DeleteUser)
	protected.Post("/getUsers", uc.ListUsers)

	protected.Post("/", uc.UpdateUser)
	protected.Get("/edit", uc.EditUser)
	protected.Post("/setAdmin/:id", uc.SetAdmin)
	protected.Post("/reset", uc.ResetPass)
	protected.Get("/getUsers", uc.GetUsers)

}
