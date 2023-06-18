package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lingo/lingo/app"
	"lingo/lingo/database"
	"lingo/lingo/middleware"
	"net/http"
)

func APIHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/v1/login" {
		LoginHandler(w, r)
		return
	} else if r.URL.Path == "/api/v1/whoami" {
		WhoAmIHandler(w, r)
		return
	} else if r.URL.Path == "/api/v1/links" {
		ListLinksHandler(w, r)
		return
	} else if r.URL.Path == "/api/v1/edit/link" {
		return
	} else if r.URL.Path == "/api/v1/delete/link" {
		return
	}
}

// login route handler
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var requestData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := app.UserByEmail(database.DB, requestData.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusBadRequest)
		return
	}

	verify := app.DefaultUserService.VerifyUser(database.User{Email: requestData.Email, Password: requestData.Password})
	if !verify {
		http.Error(w, "Invalid email or password", http.StatusBadRequest)
		return
	}

	session, err := middleware.GetSession(w, r, int(user.ID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "lingo_sesssion",
		Value:    session.SessionID,
		HttpOnly: true,
	})
	return
}

// who am I handler
func WhoAmIHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("lingo_sesssion")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	session, err := middleware.GetSessionByID(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := app.UserByID(database.DB, int(session.UserID))
	w.Header().Set("Content-Type", "application/json")
	data := struct {
		message string
	}{
		message: fmt.Sprintf("Welcome, %s!", user.Username),
	}
	json.NewEncoder(w).Encode(data.message)
}

// list all the links for the logged in user
func ListLinksHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("lingo_sesssion")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	session, err := middleware.GetSessionByID(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	links, err := app.RetrieveLinksFromDB(database.DB, &session.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	data := struct {
		message []database.Link
	}{
		message: links,
	}
	json.NewEncoder(w).Encode(data.message)
}

// create a new link for the logged in user

// edit a particular link for the logged in user

// delete a particular link for the logged in user
