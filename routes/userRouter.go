package routes

import (
	usercontroller "BalkanLinGO/controllers/user"
	"BalkanLinGO/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// create fiber route for users
func UsersRouter(app *fiber.App, session *session.Store) {
	route := app.Group("/user")
	route.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})
	route.Post("/", func(c *fiber.Ctx) error {
		return middleware.CheckAuth(c, session)
	}, usercontroller.UpdateUser)
	route.Get("/all", usercontroller.GetUsers)
	//route.Delete("/:id", usercontroller.DeleteUser)
	// TODO createpass i reset
	route.Post("/create", usercontroller.CreateUser)
	route.Post("/login", func(c *fiber.Ctx) error { return usercontroller.LoginUser(c, session) })
	route.Get("/logout", func(c *fiber.Ctx) error { return middleware.CheckAuth(c, session) }, func(c *fiber.Ctx) error { return usercontroller.LogoutUser(c, session) })
	route.Get("/getUsers", func(c *fiber.Ctx) error { return middleware.CheckAuth(c, session) }, usercontroller.GetUsers)
	route.Post("/getUsers", usercontroller.ListUsers)
	route.Get("/edit", func(c *fiber.Ctx) error {
		return middleware.CheckAuth(c, session)
	}, usercontroller.EditUser)
	route.Post("/setAdmin/:id", func(c *fiber.Ctx) error {
		return middleware.CheckAuth(c, session)
	}, usercontroller.SetAdmin)
}
