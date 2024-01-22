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
		// get user from session
		user := userdb.User{}
		if err := db.QueryRow("SELECT id, name, surname, email, is_admin FROM user WHERE id = ?", session).Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.IsAdmin); err != nil {
			return c.Redirect("/login")
		}

		if user.IsAdmin == 0 {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Nemate pristup ovoj stranici!", "link": "/dashboard"})
		}

		// set locals
		c.Locals("name", user.Name)
		c.Locals("surname", user.Surname)
		c.Locals("email", user.Email)
		c.Locals("is_admin", user.IsAdmin)

		return c.Next()
	}
}
