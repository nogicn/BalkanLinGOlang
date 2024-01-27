package user_controller

import (
	"BalkanLinGO/db"
	"BalkanLinGO/middleware"
	"BalkanLinGO/models/userdb"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
)

// get all users from database
func GetUsers(c *fiber.Ctx) error {
	users, err := userdb.GetAllUsers(db.DB)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Render("userSearch", fiber.Map{"users": users, "title": "User Search", "IsAdmin": c.Locals("is_admin")})
}

// DeleteUser deletes a user by ID
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	// Convert id to integer
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}

	// Call the DeleteUser function from the user model
	err = userdb.DeleteUserById(db.DB, idInt)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Status(200).SendString("User deleted")
}

func CreateUser(c *fiber.Ctx) error {
	name := c.FormValue("name")

	surname := c.FormValue("surname")
	email := c.FormValue("email")

	// create random string for password

	password := randStringBytes(12)

	newuser := userdb.User{
		Name:     name,
		Surname:  surname,
		Email:    email,
		Password: password,
		IsAdmin:  0,
	}
	err := userdb.CreateUser(db.DB, &newuser)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	err = middleware.SendEmail(email)

	if err != nil {
		return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri kreiranju korisnika!", "link": "/login"})
	} else {
		return c.Render("resetPassNotif", fiber.Map{})
	}
}

func LoginUser(c *fiber.Ctx, s *session.Store) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	// Check if user exists
	user, err := userdb.GetUserByEmail(db.DB, email)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	fmt.Println(user)
	// Compare passwords
	if user.Password == password {
		c.Render("createPass", fiber.Map{"email": email})
	} else {
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Pogrešna lozinka ili korisnik!", "link": "/login"})
		}

		// Create session
		session, err := s.Get(c)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		// create token
		token := randStringBytes(32)
		// Update token in database
		user, err = userdb.UpdateTokenByEmail(db.DB, email, token)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		// Set user as authenticated
		session.Set("is_admin", user.IsAdmin)
		session.Set("user_id", user.ID)
		session.Set("token", user.Token)

		c.Locals("user_id", user.ID)
		c.Locals("name", user.Name)
		c.Locals("surname", user.Surname)
		c.Locals("email", user.Email)
		c.Locals("is_admin", user.IsAdmin)
		c.Locals("token", user.Token)

		err = session.Save()
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.Redirect("/dashboard")
	}

	return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Pogrešna lozinka ili korisnik!", "link": "/login"})
}

func LogoutUser(c *fiber.Ctx, s *session.Store) error {
	session, err := s.Get(c)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	err = session.Destroy()
	if err != nil {
		c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri odjavi!", "link": "/"})
	}
	return c.Redirect("/")
}

func SetAdmin(c *fiber.Ctx) error {
	id := c.Params("id")

	// Convert id to integer
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}

	// get user from database
	curUser, err := userdb.GetUserById(db.DB, c.Locals("user_id").(int))
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	if curUser.ID == idInt {
		return c.Status(500).SendString("Ne možete postaviti sami sebe za administratora!")
	}

	// Call the DeleteUser function from the user model
	user, err := userdb.SetAdminById(db.DB, idInt)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	_, err = userdb.GetAllUsers(db.DB)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	// send partial to client but render it first
	return c.Render("partials/userRow", fiber.Map{"users": user})

}

func ListUsers(c *fiber.Ctx) error {
	users, err := userdb.GetAllUsersLikeEmail(db.DB, c.FormValue("email"))
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Render("partials/userList", fiber.Map{"users": users})
}

func EditUser(c *fiber.Ctx) error {
	id := c.Locals("user_id").(int)
	// check if nil
	if id == 0 {
		return c.Status(500).SendString("Error")
	}

	user, err := userdb.GetUserById(db.DB, id)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.Render("userEdit", fiber.Map{"user": user})
}

func UpdateUser(c *fiber.Ctx) error {

	name := c.FormValue("name")
	surname := c.FormValue("surname")

	_, err := userdb.UpdateUserByToken(db.DB, name, surname, c.Locals("token").(string))
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Redirect("/user/edit")
}
