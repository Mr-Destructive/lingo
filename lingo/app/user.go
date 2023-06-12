package app

import (
	"database/sql"
	"lingo/lingo/database"

	"golang.org/x/crypto/bcrypt"
)

func UserByUsername(db *sql.DB, username string) (*database.User, error) {
	query := "SELECT id, username, password FROM users WHERE username = ?"
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
