package app

import (
	"fmt"
	"html/template"
	"lingo/lingo/database"
	"lingo/lingo/middleware"
	"log"
	"net/http"
)

func AddLinkHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"lingo/templates/base.tmpl",
		"lingo/templates/addLink.html",
	}
	if r.Method == http.MethodGet {
		templates, err := template.ParseFiles(files...)
		if err != nil {
			log.Fatal(err)
		}
		err = templates.ExecuteTemplate(w, "base", nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		name := r.FormValue("name")
		url := r.FormValue("url")

		if name == "" || url == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		session, err := middleware.GetLoggedSession(w, r)
		user, err := database.UserByID(database.DB, session.UserID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(user)

		link := database.Link{
			Name:   name,
			URL:    url,
			UserID: user.ID,
		}

		err = database.CreateLink(database.DB, &link)
		if err != nil {
			log.Fatal(err)
		}
		http.Redirect(w, r, "/links/"+name, http.StatusFound)
	}
}
