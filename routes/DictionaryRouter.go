package routes

import (
	//"BalkanLinGO/controllers"
	"BalkanLinGO/controllers/dictionarycontroller"
	"BalkanLinGO/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func DictionaryRouter(app *fiber.App, session *session.Store) {

	route := app.Group("/dictionary")
	route.Use(func(c *fiber.Ctx) error { return middleware.CheckAuth(c, session) })

	route.Get("/addDictionary", dictionarycontroller.AddDictionary)
	route.Get("/adminEditDict/:id", dictionarycontroller.AdminEditDict)
	route.Post("/adminSaveDict", dictionarycontroller.AdminSaveDict)
	route.Get("/removeDictionary/:id", dictionarycontroller.RemoveDictionary)
	route.Get("/addDictionaryToUser/:id", dictionarycontroller.AddDictionaryToUser)
	route.Get("/dictSearch/:id", dictionarycontroller.SearchDictionary)

	route.Get("/editWord/:id", dictionarycontroller.EditWord)
	route.Post("/editWord/:id", dictionarycontroller.SaveWord)
	route.Get("/addWord/:id", dictionarycontroller.AddWord)
	route.Post("/addWord/:id", dictionarycontroller.SaveWord)
	route.Get("/deleteWord/:id", dictionarycontroller.DeleteWord)

	route.Post("/search/:id", dictionarycontroller.SearchWords)
	route.Post("/search2/:id", dictionarycontroller.SearchWords2)

	route.Get("/adminLocales", dictionarycontroller.AdminLocales)
	route.Get("/editLocale/:id", dictionarycontroller.EditLocale)
	route.Post("/saveLocale", dictionarycontroller.SaveLocale)
}
