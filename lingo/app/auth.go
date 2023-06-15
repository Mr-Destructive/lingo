package app

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"lingo/lingo/database"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type TemplateData struct {
	Data string
}

type Session struct {
	UserID    int
	SessionID string
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

func getUserData(r *http.Request) (database.User, error) {
	email := r.FormValue("email")
	user, err := UserByEmail(database.DB, email)
	if err != nil {
		return database.User{}, err
	}
	return *user, nil
}

func getUser(r *http.Request) (database.User, error) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	if user, _ := UserByEmail(database.DB, email); &user == nil {
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
	verified := DefaultUserService.VerifyUser(userForm)
	fileName := "profile.html"
	if verified {
		data := "Login Success"
		RenderTemplate(w, data, fileName)
		session, err := GetSession(w, r, int(user.ID))
		if err != nil {
			log.Fatal(err)
		}
		user, err := UserByID(database.DB, int(session.UserID))
		data = fmt.Sprintf("Welcome, %s!", user.Username)
		RenderTemplate(w, data, fileName)
		http.Redirect(w, r, "/links", http.StatusFound)
		return
	}
	data := "Login Failed"
	RenderTemplate(w, data, fileName)
	return
}

func RenderTemplate(w http.ResponseWriter, data, fileName string) {
	d := TemplateData{Data: data}
	filepath := fmt.Sprintf("lingo/templates/%s", fileName)
	t, err := template.ParseFiles(filepath)
	if err != nil {
		log.Fatal(err)
	}
	t.ExecuteTemplate(w, fileName, d)
}

const sessionIDChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
const sessionIDLength = 32

func NewSession(w http.ResponseWriter, r *http.Request, userID int) Session {
	userSession, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		userSession = &http.Cookie{
			Name:  "session_id",
			Value: "",
		}
	}
	if userSession.Value == "" {
		var sessionID = make([]byte, sessionIDLength)
		source := rand.NewSource(time.Now().UnixNano())
		rand_source := rand.New(source)
		for i := range sessionID {
			sessionID[i] = sessionIDChars[rand_source.Intn(len(sessionIDChars))]
		}

		cookie := http.Cookie{
			Name:  "session_id",
			Value: string(sessionID),
		}
		http.SetCookie(w, &cookie)

		db := database.DB
		stmt, err := db.Prepare("INSERT INTO sessions (id, user_id, session_id) VALUES (?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		_, err = stmt.Exec(userID, userID, string(sessionID))
		if err != nil {
			log.Fatal(err)
		}

		session := Session{
			UserID: userID,
		}
		return session
	}
	return Session{}
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

func GetSession(w http.ResponseWriter, r *http.Request, userID int) (Session, error) {
	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		db := database.DB
		var sessionID string
		row := db.QueryRow("SELECT session_id FROM sessions WHERE user_id = ?", userID)
		err = row.Scan(&sessionID)
		if err == sql.ErrNoRows {
			return NewSession(w, r, userID), nil
		} else if err != nil {
			return Session{}, err
		}

		cookie = &http.Cookie{
			Name:  "session_id",
			Value: sessionID,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	} else if err != nil {
		return Session{}, err
	}

	db := database.DB
	var sessionID string
	row := db.QueryRow("SELECT session_id FROM sessions WHERE session_id = ?", cookie.Value)
	err = row.Scan(&sessionID)
	if err == sql.ErrNoRows {
		return NewSession(w, r, userID), nil
	} else if err != nil {
		return Session{}, err
	}

	return Session{
		UserID:    userID,
		SessionID: sessionID,
	}, nil
}

func GetLoggedSession(w http.ResponseWriter, r *http.Request) (Session, error) {
	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		return Session{}, err
	} else if err != nil {
		return Session{}, err
	}

	db := database.DB
	var userID int
	var sessionID string
	row := db.QueryRow("SELECT user_id, session_id FROM sessions WHERE session_id = ?", cookie.Value)
	err = row.Scan(&userID, &sessionID)
	if err == sql.ErrNoRows {
		return Session{}, err
	} else if err != nil {
		return Session{}, err
	}

	return Session{
		UserID:    userID,
		SessionID: sessionID,
	}, nil
}
