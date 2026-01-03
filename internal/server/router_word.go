package server

import (
	//"BalkanLinGO/controllers"

	wordcontroller "BalkanLinGO/internal/server/controllers/word"
	"BalkanLinGO/internal/server/middleware"

	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) RegisterWordRouter() {
	wordcontroller := wordcontroller.New(s.db)

	route := s.Group("/word")
	route.Use(func(c *fiber.Ctx) error { return middleware.HtmxMiddleware(c) })
	route.Use(func(c *fiber.Ctx) error { return middleware.CheckAuth(c, s.session, s.db) })
	route.Get("/editWord/:id", wordcontroller.EditWord)
	route.Post("/editWord/:id", wordcontroller.SaveWord)
	route.Get("/addWord/:id", wordcontroller.AddWord)
	route.Post("/addWord/:id", wordcontroller.SaveWord)
	route.Delete("/deleteWord/:wordId/:dictId", wordcontroller.DeleteWord)

	//route.Post("/fillWordData/:id", wordcontroller.FillWordData)
	route.Post("/createPronunciation/:id", wordcontroller.CreatePronunciation)

}
