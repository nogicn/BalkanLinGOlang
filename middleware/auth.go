package middleware

import (
	"BalkanLinGO/db"
	"BalkanLinGO/models/userdb"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// check if user is authenticated by checking if there is a session and comparing it to the database

func CheckAuth(c *fiber.Ctx, store *session.Store) error {

	db := db.DB
	// get session
	session, _ := store.Get(c)
	// get user with session id 1
	user, err := userdb.GetUserById(db, 1)

	c.Locals("id", user.ID)
	c.Locals("name", user.Name)
	c.Locals("surname", user.Surname)
	c.Locals("email", user.Email)
	c.Locals("is_admin", user.IsAdmin)
	c.Locals("token", user.Token)

	session.Set("id", user.ID)
	session.Set("name", user.Name)
	session.Set("surname", user.Surname)
	session.Set("email", user.Email)
	session.Set("is_admin", user.IsAdmin)
	session.Set("token", user.Token)
	err = session.Save()
	return c.Next()

	if err != nil {
		fmt.Println(err)
		// if htmx, send redirect
		if c.Get("HX-Request") == "true" {
			// set header
			c.Set("HX-Redirect", "/login")
			return c.SendStatus(401)
		}
		return c.Status(401).Render("forOfor", fiber.Map{"status": "401", "errorText": "Niste prijavljeni!", "link": "/login"})
	}
	if session.Get("id") == nil {
		fmt.Println(err)
		if c.Get("HX-Request") == "true" {
			// set header
			c.Set("HX-Redirect", "/login")
			return c.SendStatus(401)
		}
		return c.Status(401).Render("forOfor", fiber.Map{"status": "401", "errorText": "Niste prijavljeni!", "link": "/login"})
	}
	// get user from session
	_, err = userdb.GetUserById(db, session.Get("id").(int))
	if err != nil {
		fmt.Println(err)
		if c.Get("HX-Request") == "true" {
			// set header
			c.Set("HX-Redirect", "/login")
			return c.SendStatus(401)
		}
		return c.Status(401).Render("forOfor", fiber.Map{"status": "401", "errorText": "Niste prijavljeni!", "link": "/login"})
	}

	// set locals
	c.Locals("id", session.Get("id").(int))
	c.Locals("name", session.Get("name").(string))
	c.Locals("surname", session.Get("surname").(string))
	c.Locals("email", session.Get("email").(string))
	c.Locals("is_admin", session.Get("is_admin").(int))
	c.Locals("token", session.Get("token").(string))

	return c.Next()
}
