package main

import (
	"BalkanLinGO/db"
	"BalkanLinGO/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	//"github.com/gofiber/template/html/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/django/v3"
	"github.com/joho/godotenv"
)

//templ generate

func main() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//engine := html.New("./views", ".html")
	engine := django.New("./views", ".html")
	//engine.Reload(true) // Optional. Default: false
	//engine.Debug(true)  // Optional. Default: false

	/*Store := memory.New(memory.Config{
		GCInterval: 6000 * time.Second,
	})*/
	session := session.New()
	app := fiber.New(fiber.Config{
		Views:                engine,
		ReduceMemoryUsage:    true,
		CompressedFileSuffix: ".fiber.gz",
		//Prefork:           true,
	})
	// debug
	app.Use(logger.New())

	// add store to ap

	// create locals
	/*app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,

		Next: func(c *fiber.Ctx) bool {
			// change headers to accept gzip if the request contains both gzip and br
			if c.Request().Header.Peek("Accept-Encoding") != nil {
				bits := c.Request().Header.Peek("Accept-Encoding")
				if bits != nil {
					if bytes.Contains(bits, []byte("gzip")) && bytes.Contains(bits, []byte("br")) {
						c.Request().Header.Set("Accept-Encoding", "gzip")
					}

				}
				return true // set to true to disable
			}
			return true
		},
	}))*/
	//app.Use(pprof.New())

	//app.Use(cache.New())
	db.Init()
	//defer db.DB.Close()
	app.Static("/", "./public", fiber.Static{
		Compress: false,
	})

	routes.UsersRouter(app, session)
	routes.IndexRouter(app, session)
	routes.DictionaryRouter(app, session)

	/*app.Get("/test", func(c *fiber.Ctx) error {
		handler := adaptor.HTTPHandler(templ.Handler(home.Home()))

		return handler(c)
	})*/

	log.Fatal(app.Listen(":3000"))

}
