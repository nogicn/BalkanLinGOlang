package server

import (
	//"BalkanLinGO/controllers"

	wordcontroller "BalkanLinGO/internal/server/controllers/word"
	"BalkanLinGO/internal/server/middleware"

	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) RegisterWordRouter() {
	wc := wordcontroller.New(s.db)

	route := s.Group("/word")
	route.Use(func(c *fiber.Ctx) error { return middleware.HtmxMiddleware(c) })
	route.Use(func(c *fiber.Ctx) error { return middleware.CheckAuth(c, s.session, s.db) })
	route.Get("/editWord/:id", wc.EditWord)
	route.Post("/editWord/:id", wc.SaveWord)
	route.Get("/addWord/:id", wc.AddWord)
	route.Post("/addWord/:id", wc.SaveWord)
	route.Delete("/deleteWord/:wordId/:dictId", wc.DeleteWord)

	//route.Post("/fillWordData/:id", wordcontroller.FillWordData)
	route.Post("/createPronunciation/:id", wc.CreatePronunciation)

}
