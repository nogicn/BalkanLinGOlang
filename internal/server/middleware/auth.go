package middleware

import (
	"BalkanLinGO/internal/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// check if user is authenticated by checking if there is a session and comparing it to the database

func CheckAuth(c *fiber.Ctx, s *session.Store, DB db.Service) error {

	repo := DB.GetRepository()
	// get session
	session, err := s.Get(c)

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	/*if session == nil || session.Get("user_id") == nil {
			//return c.Status(401).Render("forOfor", fiber.Map{"status": "401", "errorText": "Niste prijavljeni!", "link": "/login", "auth": false})
			return c.Redirect("/login")
		}

	userid := session.Get("user_id").(int)
	user, _ := userdb.GetUserByID(repo, userid)*/

	user, _ := repo.GetUserByID(c.Context(), 1)

	c.Locals("user_id", user.ID)
	c.Locals("name", user.Name)
	c.Locals("surname", user.Surname)
	c.Locals("email", user.Email)
	c.Locals("is_admin", user.IsAdmin)
	c.Locals("token", user.Token)

	session.Set("user_id", user.ID)
	session.Set("is_admin", user.IsAdmin)
	err = session.Save()
	if err != nil {
		return err
	}
	return c.Next()

}
