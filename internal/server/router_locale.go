package server

import (
	localecontroller "BalkanLinGO/internal/server/controllers/locale"
	"BalkanLinGO/internal/server/middleware"

	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) RegisterLocaleRouter() {
	lc := localecontroller.New(s.db)

	route := s.Group("/locale")
	route.Use(
		func(c *fiber.Ctx) error { return middleware.CheckAuth(c, s.session, s.db) },
		middleware.IsAdmin(s.session),
		func(c *fiber.Ctx) error { return middleware.HtmxMiddleware(c) })

	route.Get("/adminLocales", lc.AdminLocales)
	route.Get("/addLocale", lc.AddLocale)
	route.Get("/editLocale/:id", lc.EditLocale)
	route.Post("/saveLocale", lc.SaveLocale)
	route.Get("/deleteLocale/:id", lc.DeleteLocale)

}
