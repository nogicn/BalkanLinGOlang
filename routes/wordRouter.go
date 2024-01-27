package routes

import (
	//"BalkanLinGO/controllers"

	wordcontroller "BalkanLinGO/controllers/word"
	"BalkanLinGO/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func WordRouter(app *fiber.App, session *session.Store) {

	route := app.Group("/word")
	route.Use(func(c *fiber.Ctx) error { return middleware.CheckAuth(c, session) })
	route.Get("/editWord/:id", wordcontroller.EditWord)
	route.Post("/editWord/:id", wordcontroller.SaveWord)
	route.Get("/addWord/:id", wordcontroller.AddWord)
	route.Post("/addWord/:id", wordcontroller.SaveWord)
	route.Get("/deleteWord/:id", wordcontroller.DeleteWord)

}
