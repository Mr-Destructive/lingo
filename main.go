package main

import (
	"fmt"
	"lingo/lingo/database"
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

	fmt.Fprintf(w, "Hello!")
}

func main() {
	http.HandleFunc("/", helloHandler)
	fmt.Println("hello")
	db, err := database.Connect("lingo.db")
	HandleError(err)
	fmt.Println(db)
	fmt.Printf("Starting server at port 8000\n")
	err = http.ListenAndServe(":8000", nil)
	HandleError(err)
}
