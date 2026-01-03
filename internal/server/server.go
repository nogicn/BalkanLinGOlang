package server

import (
	"BalkanLinGO/internal/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/django/v3"
)

type FiberServer struct {
	*fiber.App
	db      db.Service
	session *session.Store
}

func New(dbUrl string) *FiberServer {
	engine := django.New("./internal/server/web/views", ".html")
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader:      "BalkanLinGO",
			AppName:           "BalkanLinGO",
			Views:             engine,
			ReduceMemoryUsage: false,
		}),

		db:      db.New(dbUrl),
		session: session.New(),
	}

	server.Static("/", "./internal/server/web/public/", fiber.Static{
		Compress: false,
	})

	//app.Use(logger.New())

	return server
}
