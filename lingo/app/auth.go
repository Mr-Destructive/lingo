package app

import (
	"errors"
	"html/template"
	"lingo/lingo/database"
	"log"
	"net/http"
)

type TemplateData struct {
	Data string
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/auth/login":
		Login(w, r)
	case "/auth/signup":
		Signup(w, r)
	case "/auth/login-form":
		getLoginForm(w, r)
	case "/auth/signup-form":
		getSignupForm(w, r)
	}
}

func getUser(r *http.Request) (database.User, error) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	if user, err := UserByEmail(database.DB, email); user != nil {
		log.Println("User already exists", err)
		errorname := "User already exists"
		userError := errors.New(errorname)
		return database.User{}, userError
	}
	return database.User{
		Email:    email,
		Password: password,
		Username: username,
	}, nil
}

func Login(w http.ResponseWriter, r *http.Request) {

}

func Signup(w http.ResponseWriter, r *http.Request) {
	newUser, err := getUser(r)
	if err != nil {
		t, _ := template.ParseFiles("lingo/templates/signup.html")
		data := TemplateData{Data: "User already exist with the email"}
		t.ExecuteTemplate(w, "signup.html", data)
		return
	}
	err = DefaultUserService.createUser(newUser)
	fileName := "lingo/templates/signup.html"
	if err != nil {
		t, _ := template.ParseFiles(fileName)
		data := TemplateData{Data: "User Signup Failed"}
		t.ExecuteTemplate(w, fileName, data)
		return
	}
	t, err := template.ParseFiles(fileName)
	if err != nil {
		log.Fatal(err)
	}
	data := TemplateData{Data: "New User Signup Success"}
	t.ExecuteTemplate(w, fileName, data)
	return
}

func getLoginForm(w http.ResponseWriter, r *http.Request) {
	templates, err := template.ParseFiles("lingo/templates/login.html")
	if err != nil {
		log.Fatal(err)
	}
	err = templates.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getSignupForm(w http.ResponseWriter, r *http.Request) {
	templates, err := template.ParseFiles("lingo/templates/signup.html")
	if err != nil {
		log.Fatal(err)
	}
	err = templates.ExecuteTemplate(w, "signup.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}
