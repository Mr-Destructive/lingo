package app

import (
	"database/sql"
	"fmt"
	"html/template"
	"lingo/lingo/database"
	"lingo/lingo/middleware"
	"log"
	"net/http"
	"strings"
)

type LinksTemplateData struct {
	Links []database.Link
	User  database.User
}

type LinkTemplateData struct {
	Link database.Link
}

func LinksHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"lingo/templates/base.tmpl",
		"lingo/templates/links.html",
	}
	templates, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
	}
	username := strings.Split(r.URL.Path, "/links/")
	var userId *int
	session, err := middleware.GetLoggedSession(w, r)
	if len(username) > 1 && username[1] != "" {
		userId, err = UserIdFromUsername(database.DB, username[1])
		if err != nil || userId == nil {
			return
		}
	} else {
		if err != nil {
			log.Fatal(err)
		}
		userId = &session.UserID
	}
	loggedUser := int64(session.UserID)
	links, err := retrieveLinksFromDB(database.DB, userId)
	for i, link := range links {
		if loggedUser != link.UserID {
			id := 0
			links[i].ID = int64(id)
		}
	}
	if err == sql.ErrNoRows {
		return
	}
	user, err := UserByID(database.DB, *userId)
	if err != nil {
		log.Fatal(err)
	}

	data := LinksTemplateData{
		Links: links,
		User:  *user,
	}

	err = templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Fatal(err)
	}
}

func retrieveLinksFromDB(db *sql.DB, userId *int) ([]database.Link, error) {
	query := fmt.Sprintf("SELECT id, name, url, user_id FROM links WHERE user_id = %d", *userId)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	links := []database.Link{}
	for rows.Next() {
		link := database.Link{}
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

func UserIdFromUsername(db *sql.DB, username string) (*int, error) {
	query := fmt.Sprintf("SELECT id FROM user WHERE username = '%s';", username)
	row := db.QueryRow(query)
	var userID int
	err := row.Scan(&userID)
	if err != nil {
		if userID == 0 {
			return nil, nil
		}
		return nil, err
	}
	return &userID, nil
}
