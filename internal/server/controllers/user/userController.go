package usercontroller

import (
	"BalkanLinGO/internal/db"
	dbr "BalkanLinGO/internal/db/repository"
	"BalkanLinGO/internal/server/middleware"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	repo *dbr.Queries
}

func New(dbService db.Service) *UserController {
	return &UserController{
		repo: dbService.GetRepositoryRW(),
	}
}

// GetUsers returns all users
func (uc *UserController) GetUsers(c *fiber.Ctx) error {
	users, err := uc.repo.GetAllUsers(c.Context())
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Render("user/userSearch", fiber.Map{"users": users, "title": "User Search", "IsAdmin": c.Locals("is_admin")})
}

// DeleteUser deletes a user by ID
func (uc *UserController) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	// Convert id to integer
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}

	// Call the DeleteUser function from the user model
	err = uc.repo.DeleteUserByID(c.Context(), int64(idInt))
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Status(200).SendString("User deleted")
}

func (uc *UserController) CreateUser(c *fiber.Ctx) error {
	name := c.FormValue("name")

	surname := c.FormValue("surname")
	email := c.FormValue("email")

	// create random string for password

	password := randStringBytes(12)

	err := uc.repo.CreateUser(c.Context(), dbr.CreateUserParams{
		Name:     name,
		Surname:  surname,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	err = middleware.SendEmail(email, password)

	if err != nil {
		return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri kreiranju korisnika!", "link": "/login"})
	} else {
		return c.Render("auth/authBase", fiber.Map{"reset": true, "title": "Reset Password", "resetNotif": true})
	}
}

func (uc *UserController) CreatePass(c *fiber.Ctx, s *session.Store) error {
	password := c.FormValue("password")
	password2 := c.FormValue("password2")
	email := c.FormValue("email")

	if password != password2 {
		return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Lozinke se ne poklapaju!", "link": "/login"})
	}
	hash := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(hash, bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri kreiranju lozinke!", "link": "/login"})
	}

	_, err = uc.repo.UpdatePasswordByEmail(c.Context(), dbr.UpdatePasswordByEmailParams{
		Password: string(hashedPassword),
		Email:    email,
	})
	if err != nil {
		fmt.Println(err)
		return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri kreiranju lozinke!", "link": "/login"})
	}

	err = uc.loginProcedure(c, s, dbr.User{}, email, password)
	if err != nil {
		return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri prijavi!", "link": "/login"})
	}

	return c.Redirect("/dictionary")

}

func (uc *UserController) LoginUser(c *fiber.Ctx, s *session.Store) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	// Check if user exists
	user, err := uc.repo.GetUserByEmail(c.Context(), email)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	fmt.Println(user.Password, password, user.Password == password)
	// Compare passwords
	if user.Password == password {
		return c.Render("auth/createPass", fiber.Map{"email": email, "title": "Create Password", "createPass": true})
	} else {
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Pogrešna lozinka ili korisnik!", "link": "/login"})
		}

		err = uc.loginProcedure(c, s, user, email, password)
		if err != nil {
			return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri prijavi!", "link": "/login"})
		}
		return c.Redirect("/dictionary")
	}

}

func (uc *UserController) LogoutUser(c *fiber.Ctx, s *session.Store) error {
	session, err := s.Get(c)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	err = session.Destroy()
	if err != nil {
		c.Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri odjavi!", "link": "/"})
	}
	c.Response().Header.Set("HX-Redirect", "/login")
	return c.Redirect("/login")
}

func (uc *UserController) SetAdmin(c *fiber.Ctx) error {
	id := c.Params("id")

	// Convert id to integer
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}

	// get user from database
	curUser, err := uc.repo.GetUserByID(c.Context(), int64(c.Locals("user_id").(int64)))
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	if curUser.ID == int64(idInt) {
		return c.Status(500).SendString("Ne možete postaviti sami sebe za administratora!")
	}

	// Call the DeleteUser function from the user model
	user, err := uc.repo.SetAdminByID(c.Context(), int64(idInt))
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	_, err = uc.repo.GetAllUsers(c.Context())
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	// send partial to client but render it first
	return c.Render("user/partials/userRow", fiber.Map{"users": user})

}

func (uc *UserController) ListUsers(c *fiber.Ctx) error {
	email := c.FormValue("email")

	emailNull := sql.NullString{String: email, Valid: false}
	if email != "" {
		emailNull = sql.NullString{String: email, Valid: true}
	} else {
		emailNull = sql.NullString{String: "", Valid: true}
	}

	users, err := uc.repo.GetAllUsersLikeEmail(c.Context(), emailNull)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Render("user/partials/userList", fiber.Map{"users": users})
}

func (uc *UserController) EditUser(c *fiber.Ctx) error {
	if c.Locals("user_id") == nil {
		return c.Status(500).SendString("Error")
	}

	id, ok := c.Locals("user_id").(int64)
	// check if nil
	if !ok || id == 0 {
		return c.Status(500).SendString("Error")
	}

	user, err := uc.repo.GetUserByID(c.Context(), int64(id))
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.Render("user/userEdit", fiber.Map{"user": user})
}

func (uc *UserController) UpdateUser(c *fiber.Ctx) error {

	name := c.FormValue("name")
	surname := c.FormValue("surname")

	token := c.Locals("token").(sql.NullString)
	_, err := uc.repo.UpdateUserByToken(c.Context(), dbr.UpdateUserByTokenParams{
		Name:    name,
		Surname: surname,
		Token:   token,
	})
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Redirect("/user/edit")
}

func (uc *UserController) ResetPass(c *fiber.Ctx) error {
	email := c.FormValue("email")

	user, err := uc.repo.GetUserByEmail(c.Context(), email)
	if err != nil {
		return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Korisnik ne postoji!", "link": "/login"})
	}
	// create new random password
	password := randStringBytes(12)

	err = middleware.SendEmail(email, password)
	if err != nil {
		return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri kreiranju korisnika!", "link": "/login"})
	}
	_, err = uc.repo.UpdatePasswordByEmail(c.Context(), dbr.UpdatePasswordByEmailParams{
		Password: password,
		Email:    email,
	})

	//err = middleware.SendEmail(email, user.Password)

	if err != nil {
		// return last user password
		_, errin := uc.repo.UpdatePasswordByEmail(c.Context(), dbr.UpdatePasswordByEmailParams{
			Password: user.Password,
			Email:    email,
		})
		if errin != nil {
			return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri popravljanju greške!", "link": "/login"})
		}

		return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška pri kreiranju korisnika!", "link": "/login"})
	} else {
		return c.Render("auth/resetPassNotif", fiber.Map{"title": "Reset Password", "reset": true})
	}
}
