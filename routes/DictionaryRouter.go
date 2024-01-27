package routes

import (
	//"BalkanLinGO/controllers"
	dictionarycontroller "BalkanLinGO/controllers/dictionary"
	learningcontroller "BalkanLinGO/controllers/learning"
	wordcontroller "BalkanLinGO/controllers/word"
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

	route.Get("/editWord/:id", wordcontroller.EditWord)
	route.Post("/editWord/:id", wordcontroller.SaveWord)
	route.Get("/addWord/:id", wordcontroller.AddWord)
	route.Post("/addWord/:id", wordcontroller.SaveWord)
	route.Get("/deleteWord/:id", wordcontroller.DeleteWord)

	route.Post("/search/:id", dictionarycontroller.SearchWords)

	route.Post("checkWord/:answer", learningcontroller.CheckAnswer)
	route.Post("checkWriting/:answer", learningcontroller.CheckWritingAnswer)
	route.Post("checkListening/:answer", learningcontroller.CheckListeningAnswer)

}
