package server

import (
	//"BalkanLinGO/controllers"

	dictionarycontroller "BalkanLinGO/internal/server/controllers/dictionary"
	learningcontroller "BalkanLinGO/internal/server/controllers/learning"
	"BalkanLinGO/internal/server/middleware"

	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) RegisterDictionaryRouter() {
	dc := dictionarycontroller.New(s.db)
	lc := learningcontroller.New(s.db)

	route := s.Group("/dictionary")
	route.Use(func(c *fiber.Ctx) error { return middleware.CheckAuth(c, s.session, s.db) }, func(c *fiber.Ctx) error {
		return middleware.HtmxMiddleware(c)
	})

	route.Get("/", dc.Dashboard)
	route.Get("/dictSearch/:dictId", dc.SearchDictionary)
	route.Get("/adminEditDict/:id", dc.AdminEditDict)

	route.Get("/addDictionary", dc.AddDictionary)
	route.Post("/adminSaveDict", dc.AdminSaveDict)
	route.Get("/removeDictionary/:id", dc.RemoveDictionary)
	route.Get("/addDictionaryToUser/:id", dc.AddDictionaryToUser)
	route.Post("/search/:id", dc.SearchWords)

	route.Post("/checkWord/:answer", lc.CheckAnswer)
	route.Post("/checkWriting/:answer", lc.CheckWritingAnswer)
	route.Post("/checkListening/:answer", lc.CheckListeningAnswer)
}
