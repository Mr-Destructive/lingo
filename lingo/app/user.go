package app

import (
	"database/sql"
	"lingo/lingo/database"

	"golang.org/x/crypto/bcrypt"
)

func UserExists(db *sql.DB, email, username string) bool {
	userByEmail, _ := UserByEmail(db, email)
	if userByEmail != nil {
		return true
	}
	userByUsername, _ := UserByUsername(db, username)
	if userByUsername != nil {
		return true
	}
	return false
}

func UserByID(db *sql.DB, userID int) (*database.User, error) {
	query := "SELECT id, email, username, password FROM user WHERE id = ?"
	row := db.QueryRow(query, userID)

	user := database.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return nil, err
	}
	return &user, nil
}

func UserByEmail(db *sql.DB, email string) (*database.User, error) {
	query := "SELECT id, email, username, password FROM user WHERE email = ?"
	row := db.QueryRow(query, email)

	user := database.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return nil, err
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
	database.DB.Exec("INSERT INTO user (id, email, username, password) VALUES (?, ?, ?, ?)", newId, newAuthUser.email, newAuthUser.username, newAuthUser.passwordHash)
	return nil
}

func getPasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func (userService) VerifyUser(user database.User) bool {
	dbUser, err := UserByEmail(database.DB, user.Email)
	if err != nil {
		return false
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return false
	}
	return true
}
