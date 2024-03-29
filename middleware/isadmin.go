package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// IsAdmin checks if user is authenticated
func IsAdmin(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		//db := db.DB
		// get session
		session, err := store.Get(c)
		if err != nil {
			return c.Redirect("/login")
		}
		// get user from db
		/*user, err := userdb.GetUserByID(db, session.Get("user_id").(int))
		if err != nil {
			return c.Redirect("/login")
		}

		if user.IsAdmin == 0 {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Nemate pristup ovoj stranici!", "link": "/dashboard"})
		}*/
		// check if is admin from session
		if session.Get("is_admin") == nil {
			return c.Redirect("/login")
		}

		return c.Next()
	}
}
