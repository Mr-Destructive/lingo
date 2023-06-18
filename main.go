package main

import (
	"fmt"
	"lingo/lingo/api"
	"lingo/lingo/app"
	"lingo/lingo/database"
	"lingo/lingo/middleware"
	"log"
	"net/http"
)

func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Lingo!")
}

func main() {
	db := database.InitDB("lingo.db")
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/links/", app.LinksHandler)
	http.Handle("/add/link/", middleware.AuthMiddleware(http.HandlerFunc(app.AddLinkHandler)))
	http.HandleFunc("/edit/link/", app.EditLinkHandler)
	http.HandleFunc("/delete/link/", app.DeleteLinkHandler)
	http.HandleFunc("/auth/", app.AuthHandler)

	http.HandleFunc("/api/v1/", api.APIHandler)
	fmt.Println(db)
	fmt.Printf("Starting server at port 8000\n")
	err := http.ListenAndServe(":8000", nil)
	HandleError(err)
}
