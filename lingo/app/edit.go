package app

import (
	"fmt"
	"html/template"
	"lingo/lingo/database"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func EditLinkHandler(w http.ResponseWriter, r *http.Request) {
	templates, err := template.ParseFiles("lingo/templates/editLink.html")
	linkFragment := strings.Join(strings.Split(r.URL.String(), "/")[3:], "")
	linkID, err := strconv.ParseInt(linkFragment, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if r.Method == http.MethodGet {
		// Show the edit link form
		link, err := database.GetLink(database.DB, int(linkID))
		data := LinkTemplateData{
			Link: *link,
		}
		err = templates.ExecuteTemplate(w, "editLink.html", data.Link)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println(linkID)
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
		if err != nil {
			log.Fatal(err)
		}

		link.Name = name
		link.URL = url

		err = database.UpdateLink(database.DB, link)
		if err != nil {
			log.Fatal(err)
		}
		http.Redirect(w, r, "/links/"+link.User.Username, http.StatusFound)
	}
}
