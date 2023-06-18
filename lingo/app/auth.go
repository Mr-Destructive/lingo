package app

import (
	"errors"
	"fmt"
	"html/template"
	"lingo/lingo/database"
	"lingo/lingo/middleware"
	"log"
	"net/http"
)

type TemplateData struct {
	Data string
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/auth/login":
		http.Handle("/auth/login", middleware.AuthMiddleware(http.HandlerFunc(Login)))
	case "/auth/signup":
		Signup(w, r)
	case "/auth/login-form":
		getLoginForm(w, r)
	case "/auth/signup-form":
		getSignupForm(w, r)
	case "/auth/logout":
		http.Handle("/auth/logout", middleware.AuthMiddleware(http.HandlerFunc(LogOut)))
	}
}

func getUserData(r *http.Request) (database.User, error) {
	email := r.FormValue("email")
	user, err := database.UserByEmail(database.DB, email)
	if err != nil {
		return database.User{}, err
	}
	return *user, nil
}

func getUser(r *http.Request) (database.User, error) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	if user, _ := database.UserByEmail(database.DB, email); &user == nil {
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
	user, err := getUserData(r)
	if err != nil {
		log.Fatal(err)
	}
	userForm, err := getUser(r)
	if err != nil {
		log.Fatal(err)
	}
	verified := middleware.DefaultUserService.VerifyUser(userForm)
	//userCtx := r.Context().Value("user")
	//fileName := "profile.html"
	if verified {
		session, err := middleware.GetSession(w, r, int(user.ID))
		if err != nil {
		}
		setSessionCookie(w, session.SessionID)
		user, err := database.UserByID(database.DB, int(session.UserID))
		if err != nil {
		}
		RenderTemplate(w, fmt.Sprintf("Welcome, %s!", user.Username), "profile.html")
	} else {
		RenderTemplate(w, "Login failed", "login.html")
	}
}

func LogOut(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("lingo_session")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := session.Value
	middleware.DeleteSession(w, r, sessionToken)
	http.Redirect(w, r, "/links", http.StatusSeeOther)
}

func RenderTemplate(w http.ResponseWriter, data, fileName string) {
	d := TemplateData{Data: data}
	files := []string{
		"lingo/templates/base.tmpl",
		"lingo/templates/" + fileName,
	}
	//filepath := fmt.Sprintf("lingo/templates/%s", fileName)
	t, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
	}
	t.ExecuteTemplate(w, "base", d)
}

func Signup(w http.ResponseWriter, r *http.Request) {
	newUser, err := getUser(r)
	if err != nil {
		t, _ := template.ParseFiles("lingo/templates/signup.html")
		data := TemplateData{Data: "User already exist with the email"}
		t.ExecuteTemplate(w, "signup.html", data)
		return
	}
	err = middleware.DefaultUserService.CreateUser(newUser)
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
	files := []string{
		"lingo/templates/base.tmpl",
		"lingo/templates/login.html",
	}
	templates, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
	}
	err = templates.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getSignupForm(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"lingo/templates/base.tmpl",
		"lingo/templates/signup.html",
	}
	templates, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
	}
	err = templates.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func setSessionCookie(w http.ResponseWriter, sessionID string) {
	cookie := http.Cookie{
		Name:     "lingo_session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}
