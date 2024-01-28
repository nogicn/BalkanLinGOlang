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

	return c.Render("user/userSearch", fiber.Map{"users": users, "title": "User Search", "IsAdmin": c.Locals("is_admin")})
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

	err = middleware.SendEmail(email, password)

	if err != nil {
		return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri kreiranju korisnika!", "link": "/login"})
	} else {
		return c.Render("auth/resetPassNotif", fiber.Map{})
	}
}

func CreatePass(c *fiber.Ctx, s *session.Store) error {
	password := c.FormValue("password")
	password2 := c.FormValue("password2")
	email := c.FormValue("email")

	if password != password2 {
		return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Lozinke se ne poklapaju!", "link": "/login"})
	}
	hash := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(hash, bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri kreiranju lozinke!", "link": "/login"})
	}

	_, err = userdb.UpdatePasswordByEmail(db.DB, email, string(hashedPassword))
	if err != nil {
		fmt.Println(err)
		return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri kreiranju lozinke!", "link": "/login"})
	}

	err = loginProcedure(c, s, userdb.User{}, email, password)
	if err != nil {
		return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri prijavi!", "link": "/login"})
	}

	return c.Redirect("/dashboard")

}

func LoginUser(c *fiber.Ctx, s *session.Store) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	// Check if user exists
	user, err := userdb.GetUserByEmail(db.DB, email)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	fmt.Println(user.Password, password, user.Password == password)
	// Compare passwords
	if user.Password == password {
		return c.Render("auth/createPass", fiber.Map{"email": email})
	} else {
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Pogrešna lozinka ili korisnik!", "link": "/login"})
		}

		err = loginProcedure(c, s, user, email, password)
		if err != nil {
			return c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri prijavi!", "link": "/login"})
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
	return c.Render("user/partials/userRow", fiber.Map{"users": user})

}

func ListUsers(c *fiber.Ctx) error {
	users, err := userdb.GetAllUsersLikeEmail(db.DB, c.FormValue("email"))
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Render("user/partials/userList", fiber.Map{"users": users})
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
	return c.Render("user/userEdit", fiber.Map{"user": user})
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
