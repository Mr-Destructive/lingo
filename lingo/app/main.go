package app

import (
	"database/sql"
	"fmt"
	"html/template"
	"lingo/lingo/database"
	"log"
	"net/http"
	"strings"
)

type LinksTemplateData struct {
	Links []database.Link
}

func LinksHandler(w http.ResponseWriter, r *http.Request) {
	templates, err := template.ParseFiles("lingo/templates/test.html")
	if err != nil {
		log.Fatal(err)
	}
	username := strings.Split(r.URL.Path, "/links/")
	fmt.Println(username)
	var userId *int
	if len(username) > 1 {
		userId, err = UserIdFromUsername(database.DB, username[1])
		fmt.Println("UU: ", &userId)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println(userId)
	links, err := retrieveLinksFromDB(database.DB, userId)
	if err != nil {
		log.Fatal(err)
	}
	data := LinksTemplateData{
		Links: links,
	}

	err = templates.ExecuteTemplate(w, "test.html", data)
	if err != nil {
		log.Fatal(err)
	}
}

func retrieveLinksFromDB(db *sql.DB, userId *int) ([]database.Link, error) {
	query := "SELECT id, name, url, user_id FROM links"
	fmt.Println("U:", userId)
	if userId != nil {
		query = fmt.Sprintf("SELECT id, name, url, user_id FROM links WHERE user_id = %d", *userId)
	}
	fmt.Println(query)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	links := []database.Link{}
	for rows.Next() {
		link := database.Link{}
		var userID int64
		err := rows.Scan(&link.ID, &link.Name, &link.URL, &userID)
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

func UserIdFromUsername(db *sql.DB, username string) (*int, error) {
	query := fmt.Sprintf("SELECT id FROM user WHERE username = '%s';", username)
	row := db.QueryRow(query)
	var userID int
	err := row.Scan(&userID)
	fmt.Println(userID)
	if err != nil {
		if userID == 0 {
			return nil, nil
		}
		return nil, err
	}
	return &userID, nil
}
