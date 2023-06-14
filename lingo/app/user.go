package app

import (
	"database/sql"
	"fmt"
	"lingo/lingo/database"

	"golang.org/x/crypto/bcrypt"
)

func UserByEmail(db *sql.DB, email string) (*database.User, error) {
	query := "SELECT id, email, username, password FROM user WHERE email = ?"
	row := db.QueryRow(query, email)

	user := database.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return &user, err
	}
	return &user, nil
}

func UserByUsername(db *sql.DB, username string) (*database.User, error) {
	query := "SELECT id, username, password FROM user WHERE username = ?"
	row := db.QueryRow(query, username)

	user := database.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return &user, err
	}

	return &user, nil
}

func CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

type authUser struct {
	email        string
	username     string
	passwordHash string
}

var DefaultUserService userService

type userService struct {
}

func (userService) createUser(newUser database.User) error {
	passwordHash, err := getPasswordHash(newUser.Password)
	if err != nil {
		return err
	}
	newAuthUser := authUser{
		email:        newUser.Email,
		username:     newUser.Username,
		passwordHash: passwordHash,
	}
	row := database.DB.QueryRow("SELECT id FROM user ORDER BY id DESC LIMIT 1")
	var lastId int
	err = row.Scan(&lastId)
	if err != nil {
		panic(err)
	}
	newId := lastId + 1
	fmt.Println(newId)
	database.DB.Exec("INSERT INTO user (id, email, username, password) VALUES (?, ?, ?, ?)", newId, newAuthUser.email, newAuthUser.username, newAuthUser.passwordHash)
	return nil
}

func getPasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
