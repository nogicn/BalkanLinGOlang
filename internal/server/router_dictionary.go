package server

import (
	//"BalkanLinGO/controllers"

	dictionarycontroller "BalkanLinGO/internal/server/controllers/dictionary"
	learningcontroller "BalkanLinGO/internal/server/controllers/learning"
	"BalkanLinGO/internal/server/middleware"

	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) RegisterDictionaryRouter() {
	dictController := dictionarycontroller.New(s.db)
	learningcontroller := learningcontroller.New(s.db)

	route := s.Group("/dictionary")
	route.Use(func(c *fiber.Ctx) error { return middleware.CheckAuth(c, s.session, s.db) }, func(c *fiber.Ctx) error {
		return middleware.HtmxMiddleware(c)
	})

	route.Get("/", dictController.Dashboard)
	route.Get("/dictSearch/:dictId", dictController.SearchDictionary)
	route.Get("/adminEditDict/:id", dictController.AdminEditDict)

	route.Get("/addDictionary", dictController.AddDictionary)
	route.Post("/adminSaveDict", dictController.AdminSaveDict)
	route.Get("/removeDictionary/:id", dictController.RemoveDictionary)
	route.Get("/addDictionaryToUser/:id", dictController.AddDictionaryToUser)
	route.Post("/search/:id", dictController.SearchWords)

	route.Post("/checkWord/:answer", learningcontroller.CheckAnswer)
	route.Post("/checkWriting/:answer", learningcontroller.CheckWritingAnswer)
	route.Post("/checkListening/:answer", learningcontroller.CheckListeningAnswer)
}
