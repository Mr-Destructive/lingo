package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lingo/lingo/database"
	"lingo/lingo/middleware"
	"net/http"
	"strings"
)

func APIHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/v1/login" {
		LoginHandler(w, r)
		return
	} else if r.URL.Path == "/api/v1/whoami" {
		WhoAmIHandler(w, r)
		return
	} else if r.URL.Path == "/api/v1/links" && r.Method == http.MethodGet {
		ListLinksHandler(w, r)
		return
	} else if r.URL.Path == "/api/v1/add/link" && r.Method == http.MethodPost {
		CreateLinkHandler(w, r)
		return
	} else if strings.HasPrefix(r.URL.String(), "/api/v1/edit/link/") && r.Method == http.MethodPut {
		editLinkHandler(w, r)
		return
	} else if strings.HasPrefix(r.URL.Path, "/api/v1/delete/link") && r.Method == http.MethodDelete {
		deleteLinkHandler(w, r)
		return
	} else if r.URL.Path == "/api/v1/profile" && r.Method == http.MethodGet {
		ProfileAPIHandler(w, r)
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
	user, err := database.UserByEmail(database.DB, requestData.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusBadRequest)
		return
	}

	verify := middleware.DefaultUserService.VerifyUser(database.User{Email: requestData.Email, Password: requestData.Password})
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
	user, err := database.UserByID(database.DB, int(session.UserID))
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
	links, err := database.RetrieveLinksFromDB(database.DB, &session.UserID)
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
func CreateLinkHandler(w http.ResponseWriter, r *http.Request) {
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
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var link struct {
		Name   string `json:"name"`
		Url    string `json:"url"`
		UserID int64
	}
	err = json.Unmarshal(body, &link)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	link.UserID = int64(session.UserID)
	linkObj := database.Link{
		Name:   link.Name,
		URL:    link.Url,
		UserID: link.UserID,
	}
	err = database.CreateLink(database.DB, &linkObj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		data := struct {
			message database.Link
		}{
			message: linkObj,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data.message)
	}
}

// edit a particular link for the logged in user
func editLinkHandler(w http.ResponseWriter, r *http.Request) {
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
	pathSegments := strings.Split(r.URL.Path, "/")
	linkName := pathSegments[len(pathSegments)-1]
	if linkName != "" {
		link, err := database.GetLinkByName(database.DB, linkName, session.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if int(session.UserID) != session.UserID {
			http.Error(w, "You are not authorized to edit this link", http.StatusBadRequest)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var linkBody struct {
			Name   string `json:"name"`
			Url    string `json:"url"`
			UserID int64
		}
		err = json.Unmarshal(body, &linkBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		linkBody.UserID = int64(session.UserID)
		link.Name = linkBody.Name
		link.URL = linkBody.Url
		err = database.UpdateLink(database.DB, link)
		editedLink, err := database.GetLink(database.DB, int(link.ID))
		if err == nil {
			data := struct {
				message database.Link
			}{
				message: *editedLink,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(data.message)
		}
	}
}

// delete a particular link for the logged in user
func deleteLinkHandler(w http.ResponseWriter, r *http.Request) {
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
	pathSegments := strings.Split(r.URL.Path, "/")
	linkName := pathSegments[len(pathSegments)-1]
	if linkName != "" {
		link, err := database.GetLinkByName(database.DB, linkName, session.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if int(session.UserID) != session.UserID {
			http.Error(w, "You are not authorized to edit this link", http.StatusBadRequest)
			return
		}
		err = database.DeleteLink(database.DB, link)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			data := struct {
				message string
			}{
				message: "Deleted Successfully",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(data.message)
		}
	}
}
