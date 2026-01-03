package routes

import (
	usercontroller "BalkanLinGO/controllers/user"
	"BalkanLinGO/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func UsersRouter(app *fiber.App, session *session.Store) {
	// public routes (no auth)
	route := app.Group("/user")

	route.Post("/login", func(c *fiber.Ctx) error { return usercontroller.LoginUser(c, session) })
	route.Post("/register", usercontroller.CreateUser)
	route.Post("/createPass", func(c *fiber.Ctx) error { return usercontroller.CreatePass(c, session) })

	// protected routes: apply HTMX middleware and authentication
	protected := app.Group("/user")
	protected.Use(func(c *fiber.Ctx) error { return middleware.CheckAuth(c, session) })
	protected.Get("/logout", func(c *fiber.Ctx) error { return usercontroller.LogoutUser(c, session) })

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
