package server

import (
	usercontroller "BalkanLinGO/internal/server/controllers/user"
	"BalkanLinGO/internal/server/middleware"

	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) RegisterUserRouter() {
	usercontroller := usercontroller.New(s.db)

	// public routes (no auth)
	route := s.Group("/user")

	route.Post("/login", func(c *fiber.Ctx) error { return usercontroller.LoginUser(c, s.session) })
	route.Post("/register", usercontroller.CreateUser)
	route.Post("/createPass", func(c *fiber.Ctx) error { return usercontroller.CreatePass(c, s.session) })

	// protected routes: apply HTMX middleware and authentication
	protected := s.Group("/user")
	protected.Use(func(c *fiber.Ctx) error { return middleware.CheckAuth(c, s.session, s.db) })
	protected.Get("/logout", func(c *fiber.Ctx) error { return usercontroller.LogoutUser(c, s.session) })

	//use htmx for remaining routes
	protected.Use(func(c *fiber.Ctx) error { return middleware.HtmxMiddleware(c) })

	protected.Get("/all", usercontroller.GetUsers)
	//protected.Delete(":id", usercontroller.DeleteUser)
	protected.Post("/getUsers", usercontroller.ListUsers)

	protected.Post("/", usercontroller.UpdateUser)
	protected.Get("/edit", usercontroller.EditUser)
	protected.Post("/setAdmin/:id", usercontroller.SetAdmin)
	protected.Post("/reset", usercontroller.ResetPass)
	protected.Get("/getUsers", usercontroller.GetUsers)

}
