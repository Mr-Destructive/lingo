package app

import (
	"html/template"
	"lingo/lingo/database"
	"lingo/lingo/middleware"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func EditLinkHandler(w http.ResponseWriter, r *http.Request) {
	session, err := middleware.GetLoggedSession(w, r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	files := []string{
		"lingo/templates/base.tmpl",
		"lingo/templates/editLink.html",
	}
	templates, err := template.ParseFiles(files...)
	linkFragment := strings.Join(strings.Split(r.URL.String(), "/")[3:], "")
	linkID, err := strconv.ParseInt(linkFragment, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if r.Method == http.MethodGet {
		// Show the edit link form
		link, err := database.GetLink(database.DB, int(linkID))
		if int64(session.UserID) != link.UserID {
			http.Redirect(w, r, "/", http.StatusForbidden)
			return
		}
		data := LinkTemplateData{
			Link: *link,
		}
		err = templates.ExecuteTemplate(w, "base", data.Link)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		url := r.FormValue("url")

		if name == "" || url == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		link, err := database.GetLink(database.DB, int(linkID))

		if int64(session.UserID) != link.UserID {
			http.Redirect(w, r, "/", http.StatusForbidden)
			return
		}
		if err != nil {
			log.Fatal(err)
		}

		link.Name = name
		link.URL = url
		user, err := database.UserByID(database.DB, int(link.UserID))

		err = database.UpdateLink(database.DB, link)
		if err != nil {
			log.Fatal(err)
		}
		http.Redirect(w, r, "/links/"+user.Username, http.StatusFound)
	}
}
