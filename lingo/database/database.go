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

	statement, err := db.Prepare("INSERT INTO user(username, name, email, password) VALUES(?, ?, ?)")
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

func GetUser(db *sql.DB, userId int64) (*User, error) {
	user := User{}

	row := db.QueryRow("SELECT id, username, email, password FROM user WHERE id = ?", userId)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateLink(db *sql.DB, link *Link) error {
	statement, err := db.Prepare("INSERT INTO links(name, url, user_id) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(link.Name, link.URL, link.UserID)
	if err != nil {
		return err
	}

	return nil
}

func GetLink(db *sql.DB, linkId int) (*Link, error) {
	link := Link{}

	row := db.QueryRow("SELECT * FROM links WHERE id = ?", linkId)
	err := row.Scan(&link.ID, &link.Name, &link.URL, &link.UserID)
	if err != nil {
		return nil, err
	}

	user, err := GetUser(db, link.UserID)
	if err != nil {
		return nil, err
	}
	link.UserID = user.ID

	return &link, nil
}
func GetLinkByName(db *sql.DB, linkName string, userID int) (*Link, error) {
	link := Link{}

	row := db.QueryRow("SELECT * FROM links WHERE name = ? AND user_id = ?", linkName, userID)
	err := row.Scan(&link.ID, &link.Name, &link.URL, &link.UserID)
	if err != nil {
		return nil, err
	}

	user, err := GetUser(db, link.UserID)
	if err != nil {
		return nil, err
	}
	link.UserID = user.ID

	return &link, nil
}

func UpdateLink(db *sql.DB, link *Link) error {
	statement, err := db.Prepare("UPDATE links SET name = ?, url = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(link.Name, link.URL, link.ID)
	if err != nil {
		return err
	}

	return nil
}

func RetrieveLinksFromDB(db *sql.DB, userId *int) ([]Link, error) {
	query := fmt.Sprintf("SELECT id, name, url, user_id FROM links WHERE user_id = %d", *userId)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	links := []Link{}
	for rows.Next() {
		link := Link{}
		err := rows.Scan(&link.ID, &link.Name, &link.URL, &link.UserID)
		if err != nil {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		links = append(links, link)
	}

	return links, nil
}

func DeleteLink(db *sql.DB, link *Link) error {
	statement, err := db.Prepare("DELETE FROM links WHERE id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(link.ID)
	if err != nil {
		return err
	}
	return nil
}

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

func UserByID(db *sql.DB, userID int) (*User, error) {
	query := "SELECT id, email, username, password FROM user WHERE id = ?"
	row := db.QueryRow(query, userID)

	user := User{}
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return nil, err
	}
	return &user, nil
}

func UserByEmail(db *sql.DB, email string) (*User, error) {
	query := "SELECT id, email, username, password FROM user WHERE email = ?"
	row := db.QueryRow(query, email)

	user := User{}
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return nil, err
	}
	return &user, nil
}

func UserByUsername(db *sql.DB, username string) (*User, error) {
	query := "SELECT id, username, password FROM user WHERE username = ?"
	row := db.QueryRow(query, username)

	user := User{}
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user %s not found", username)
	}
	if err != nil {
		return &user, err
	}

	return &user, nil
}

func GetProfile(db *sql.DB, userId int) (*Profile, error) {
	query := "SELECT id, color, avatar FROM profiles WHERE user_id = ?"
	linkQuery := "SELECT id, name, url FROM links WHERE profile_id = ?"
	row := db.QueryRow(query, userId)
	profile := Profile{}
	err := row.Scan(&profile.ID, &profile.Color, &profile.Avatar)
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(linkQuery, profile.ID)
	links := []Link{}
	for rows.Next() {
		link := Link{}
		err := rows.Scan(&link.ID, &link.Name, &link.URL)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	profile.Links = links
	profile.UserID = int64(userId)

	if err != nil {
		return nil, err
	}
	return &profile, nil
}
