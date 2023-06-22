package database

import (
	"database/sql"
	"log"
	//"net/url"
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Link struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	URL    string `json:"link"`
	UserID int64  `json:"user_id"`
}

type Session struct {
	UserID    int
	SessionID string
}

type Project struct {
    ID       int64  `json:"id"`
    Name     string `json:"name"`
    Url      string `json:"url"`
    UserID   int64  `json:"user_id"`
}

type Profile struct {
    ID       int64  `json:"id"`
    UserID   int64  `json:"user_id"`
    Color    string `json:"color"`
    Avatar   string `json:"avatar"`
    Links    []Link `json:"links"`
}

var DB *sql.DB

func InitDB(dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	err = runMigrations(db)
	if err != nil {
		return err
	}

	log.Println("Database connection established")

	DB = db
	return nil
}
