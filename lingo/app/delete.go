package app

import (
	"fmt"
	"html/template"
	"lingo/lingo/database"
	"lingo/lingo/middleware"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func DeleteLinkHandler(w http.ResponseWriter, r *http.Request) {
	session, err := middleware.GetLoggedSession(w, r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	files := []string{
		"lingo/templates/base.tmpl",
		"lingo/templates/deleteLink.html",
	}
	templates, err := template.ParseFiles(files...)
	linkFragment := strings.Join(strings.Split(r.URL.String(), "/")[3:], "")
	linkID, err := strconv.ParseInt(linkFragment, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	link, err := database.GetLink(database.DB, int(linkID))
	if r.Method == http.MethodGet {
		if int64(session.UserID) != link.UserID {
			http.Redirect(w, r, "/", http.StatusForbidden)
			return
		}
		data := LinkTemplateData{
			Link: *link,
		}
		err = templates.ExecuteTemplate(w, "base", data.Link)
		return
	} else if r.Method == http.MethodPost {
		if int64(session.UserID) != link.UserID {
			http.Redirect(w, r, "/", http.StatusForbidden)
			return
		}
		query := fmt.Sprintf("DELETE FROM links WHERE id = %d", linkID)
		_, err = database.DB.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	}

	http.Redirect(w, r, "/links", http.StatusFound)
}
