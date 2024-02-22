package routes

import (
	usercontroller "BalkanLinGO/controllers/user"
	"BalkanLinGO/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// UsersRouter crates routes for user
func UsersRouter(app *fiber.App, session *session.Store) {
	route := app.Group("/user")

	route.Post("/", func(c *fiber.Ctx) error {
		return middleware.CheckAuth(c, session)
	}, usercontroller.UpdateUser)
	route.Get("/all", usercontroller.GetUsers)
	//route.Delete("/:id", usercontroller.DeleteUser)
	route.Post("/register", usercontroller.CreateUser)
	route.Post("/login", func(c *fiber.Ctx) error { return usercontroller.LoginUser(c, session) })
	route.Get("/logout", func(c *fiber.Ctx) error { return middleware.CheckAuth(c, session) }, func(c *fiber.Ctx) error { return usercontroller.LogoutUser(c, session) })
	route.Get("/getUsers", func(c *fiber.Ctx) error { return middleware.CheckAuth(c, session) }, usercontroller.GetUsers)
	route.Post("/createPass", func(c *fiber.Ctx) error { return usercontroller.CreatePass(c, session) })
	route.Post("/getUsers", usercontroller.ListUsers)
	route.Get("/edit", func(c *fiber.Ctx) error {
		return middleware.CheckAuth(c, session)
	}, usercontroller.EditUser)
	route.Post("/setAdmin/:id", func(c *fiber.Ctx) error {
		return middleware.CheckAuth(c, session)
	}, usercontroller.SetAdmin)
	route.Post("/reset", func(c *fiber.Ctx) error {
		return middleware.CheckAuth(c, session)
	}, usercontroller.ResetPass)

}
