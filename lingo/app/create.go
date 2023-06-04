package app

import (
	"html/template"
	"lingo/lingo/database"
	"log"
	"net/http"
)

func AddLinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		templates, err := template.ParseFiles("lingo/templates/addLink.html")
		if err != nil {
			log.Fatal(err)
		}
		err = templates.ExecuteTemplate(w, "addLink.html", nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		name := r.FormValue("name")
		url := r.FormValue("url")

		if name == "" || url == "" || username == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userId, err := UserIdFromUsername(database.DB, username)
		if err != nil {
			log.Fatal(err)
		}

		link := database.Link{
			Name: name,
			URL:  url,
			User: database.User{ID: int64(*userId)},
		}

		err = database.CreateLink(database.DB, &link)
		if err != nil {
			log.Fatal(err)
		}
		http.Redirect(w, r, "/links/"+username, http.StatusFound)
	}
}
