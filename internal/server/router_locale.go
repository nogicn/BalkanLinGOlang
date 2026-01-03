package server

import (
	localecontroller "BalkanLinGO/internal/server/controllers/locale"
	"BalkanLinGO/internal/server/middleware"

	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) RegisterLocaleRouter() {
	localecontroller := localecontroller.New(s.db)

	route := s.Group("/locale")
	route.Use(func(c *fiber.Ctx) error { return middleware.CheckAuth(c, s.session, s.db) }, middleware.IsAdmin(s.session))

	route.Get("/adminLocales", func(c *fiber.Ctx) error {
		return middleware.HtmxMiddleware(c, "/")

	}, localecontroller.AdminLocales)
	route.Get("/addLocale", func(c *fiber.Ctx) error {
		return middleware.HtmxMiddleware(c, "/")
	}, localecontroller.AddLocale)

	route.Get("/editLocale/:id", func(c *fiber.Ctx) error {
		return middleware.HtmxMiddleware(c, "/")
	}, localecontroller.EditLocale)

	route.Post("/saveLocale", func(c *fiber.Ctx) error {
		return middleware.HtmxMiddleware(c, "/")
	}, localecontroller.SaveLocale)

	route.Get("/deleteLocale/:id", func(c *fiber.Ctx) error {
		return middleware.HtmxMiddleware(c, "/")
	}, localecontroller.DeleteLocale)

}
