package userdb

import (
	"database/sql"
)

const (
	createUserTable = `
	CREATE TABLE IF NOT EXISTS user (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        surname TEXT NOT NULL,
        email TEXT NOT NULL,
        password TEXT NOT NULL,
        is_admin INTEGER DEFAULT 0,
        token TEXT DEFAULT N,
        UNIQUE (email)
    );`

	createUser = `INSERT INTO user (name, surname, email, password) VALUES (@name, @surname, @email, @password);`

	createAdmin = `
		INSERT INTO user (name, surname, email, password, is_admin) VALUES (@name, @surname, @email, @password, 1);
	`

	loginEmailPassword = `
		SELECT * FROM user WHERE email = @email AND password = @password;
	`

	getAllUsers = `
		SELECT * FROM user;
	`

	getUserByToken = `
		SELECT * FROM user WHERE token = @token;
	`

	getUserById = `
		SELECT * FROM user WHERE id = @id;
	`

	getUserByEmail = `
		SELECT * from user WHERE email = @email;
	`

	updateTokenByEmail = `
		UPDATE user SET token = @token WHERE email = @email RETURNING *;
	`

	updateTokenById = `
		UPDATE user SET token = @token WHERE id = @id RETURNING *;
	`

	updatePasswordByEmail = `
		UPDATE user SET password = @password WHERE email = @email RETURNING *;
	`

	updateUserByToken = `
		UPDATE user SET name = @name, surname = @surname WHERE token = @token RETURNING *;
	`
	deleteUserById = `
		DELETE FROM user WHERE id = ?;
	`
	setAdminByEmail = `
    	UPDATE user SET is_admin = not is_admin WHERE email = @email RETURNING *;
	`
	setAdminById = `
		UPDATE user SET is_admin = not is_admin WHERE id = @id RETURNING *;
	`
	getAllUsersLikeEmail = `
		SELECT * FROM user WHERE email LIKE @email;
	`
)

type User struct {
	ID       int            `json:"id"`
	Name     string         `json:"name"`
	Surname  string         `json:"surname"`
	Email    string         `json:"email"`
	Password string         `json:"password"`
	IsAdmin  int            `json:"is_admin"`
	Token    sql.NullString `json:"token"`
}

func CreateUserTable(dbase *sql.DB) error {
	_, err := dbase.Exec(createUserTable)
	return err
}

func CreateUser(dbase *sql.DB, u *User) error {
	_, err := dbase.Exec(createUser, u.Name, u.Surname, u.Email, u.Password)
	return err
}

func CreateAdmin(dbase *sql.DB, u *User) error {
	_, err := dbase.Exec(createAdmin, u.Name, u.Surname, u.Email, u.Password)
	return err
}

func SetAdminByEmail(dbase *sql.DB, email string) (User, error) {
	var user User
	err := dbase.QueryRow(setAdminByEmail, email).Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &user.IsAdmin, &user.Token)
	return user, err
}

func SetAdminById(dbase *sql.DB, id int) (User, error) {
	var user User
	err := dbase.QueryRow(setAdminById, id).Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &user.IsAdmin, &user.Token)
	return user, err
}

func LoginEmailPassword(dbase *sql.DB, u *User) (User, error) {
	var user User
	err := dbase.QueryRow(loginEmailPassword, u.Email, u.Password).Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &user.IsAdmin, &user.Token)
	return user, err
}

func GetAllUsers(dbase *sql.DB) ([]User, error) {
	var users []User
	rows, err := dbase.Query(getAllUsers)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &user.IsAdmin, &user.Token)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func GetUserByToken(dbase *sql.DB, token string) (User, error) {
	var user User
	err := dbase.QueryRow(getUserByToken, token).Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &user.IsAdmin, &user.Token)
	return user, err
}

func GetUserById(dbase *sql.DB, id int) (User, error) {
	var user User
	err := dbase.QueryRow(getUserById, id).Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &user.IsAdmin, &user.Token)
	return user, err
}

func GetUserByEmail(dbase *sql.DB, email string) (User, error) {
	var user User
	err := dbase.QueryRow(getUserByEmail, email).Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &user.IsAdmin, &user.Token)

	return user, err
}

func UpdateTokenByEmail(dbase *sql.DB, email string, token string) (User, error) {
	var user User = User{}
	err := dbase.QueryRow(updateTokenByEmail, token, email).Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &user.IsAdmin, &user.Token)
	return user, err
}

func UpdateTokenById(dbase *sql.DB, id int, token string) (User, error) {
	var user User
	err := dbase.QueryRow(updateTokenById, id, token).Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &user.IsAdmin, &user.Token)
	return user, err
}

func UpdatePasswordByEmail(dbase *sql.DB, email string, password string) (User, error) {
	var user User
	err := dbase.QueryRow(updatePasswordByEmail, password, email).Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &user.IsAdmin, &user.Token)
	return user, err
}

func UpdateUserByToken(dbase *sql.DB, name, surname, token string) (User, error) {
	var user User
	err := dbase.QueryRow(updateUserByToken, name, surname, token).Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &user.IsAdmin, &user.Token)
	return user, err
}

func DeleteUserById(dbase *sql.DB, id int) error {
	_, err := dbase.Exec(deleteUserById, id)
	return err
}

func GetAllUsersLikeEmail(dbase *sql.DB, email string) ([]User, error) {
	var users []User
	rows, err := dbase.Query(getAllUsersLikeEmail, email+"%")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &user.IsAdmin, &user.Token)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
