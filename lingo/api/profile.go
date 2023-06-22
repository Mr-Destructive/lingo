package api

import (
	"encoding/json"
	"fmt"
	"lingo/lingo/database"
	"lingo/lingo/middleware"
	"net/http"
)

func ProfileAPIHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/v1/profile" && r.Method == http.MethodGet {
		ProfileHandler(w, r)
		return
	}
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	session, err := middleware.GetLoggedSession(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userId := session.UserID
	profile, err := database.GetProfile(database.DB, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	data := struct {
		message database.Profile
	}{
		message: *profile,
	}
	json.NewEncoder(w).Encode(data.message)

}
