package middleware

import (
	"BalkanLinGO/db"
	"BalkanLinGO/models/userdb"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// check if user is authenticated
func IsAdmin(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := db.DB
		// get session
		session, err := store.Get(c)
		if err != nil {
			return c.Redirect("/login")
		}
		// get user from db
		user, err := userdb.GetUserById(db, session.Get("user_id").(int))
		if err != nil {
			return c.Redirect("/login")
		}

		if user.IsAdmin == 0 {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Nemate pristup ovoj stranici!", "link": "/dashboard"})
		}

		return c.Next()
	}
}
