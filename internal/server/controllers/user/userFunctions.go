package usercontroller

import (
	dbr "BalkanLinGO/internal/db/repository"
	"database/sql"
	"fmt"
	"math/rand"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func randStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (uc *UserController) loginProcedure(c *fiber.Ctx, s *session.Store, user dbr.User, email string, password string) error {
	// Create session
	session, err := s.Get(c)
	if err != nil {
		fmt.Println("1", err)
		return &fiber.Error{Code: 500, Message: "Internal Server Error"}
	}
	// create token
	token := randStringBytes(32)
	// Update token in database
	user, err = uc.repo.UpdateTokenByEmail(c.Context(), dbr.UpdateTokenByEmailParams{
		Token: sql.NullString{String: token, Valid: true},
		Email: email,
	})
	if err != nil {
		fmt.Println("2", err)
		return &fiber.Error{Code: 500, Message: "Internal Server Error"}
	}

	// Set user as authenticated
	var isAdmin int
	if user.IsAdmin {
		isAdmin = 1
	}
	session.Set("is_admin", isAdmin)
	session.Set("user_id", int(user.ID))

	c.Locals("user_id", int(user.ID))
	c.Locals("name", user.Name)
	c.Locals("surname", user.Surname)
	c.Locals("email", user.Email)
	c.Locals("is_admin", isAdmin)

	err = session.Save()
	if err != nil {
		fmt.Println("3", err)
		return &fiber.Error{Code: 500, Message: "Internal Server Error"}
	}
	return nil
}
