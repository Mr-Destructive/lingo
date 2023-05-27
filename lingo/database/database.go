package database

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func Connect(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	err = runMigrations(db)
	if err != nil {
		return nil, err
	}
	fmt.Println("pinging")

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
	path := "lingo/database/" + "migrations"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			path := filepath.Join(path, file.Name())
			sql, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			_, err = db.Exec(string(sql))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func CreateUser(user *User) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	nil_user := User{}
	if err != nil {
		log.Fatal("failed to hash password")
		return nil_user, errors.New("failed to create user")
	}

	db, err := sql.Open("sqlite3", "lingo.db")
	if err != nil {
		log.Printf("failed to open database: %w", err)
		return nil_user, errors.New("failed to create user")
	}
	defer db.Close()

	statement, err := db.Prepare("INSERT INTO users(username, email, password) VALUES(?, ?, ?)")
	if err != nil {
		log.Printf("failed to prepare insert statement %w", err)
		return nil_user, errors.New("failed to create user")
	}
	defer statement.Close()

	result, err := statement.Exec(user.Username, user.Email, hashedPassword)
	if err != nil {
		log.Fatal("failed to insert user")
		return nil_user, errors.New("failed to create user")
	}

	user.ID, err = result.LastInsertId()
	if err != nil {
		log.Fatal("failed to get last insert ID")
		return nil_user, errors.New("failed to create user")
	}

	return *user, nil
}