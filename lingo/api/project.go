package api

import (
	"encoding/json"
	"fmt"
	"lingo/lingo/database"
	"lingo/lingo/middleware"
	"net/http"
)

func ProjectAPIHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/v1/projects" && r.Method == http.MethodGet {
		ProjectHandler(w, r)
		return
	}
}

func ProjectHandler(w http.ResponseWriter, r *http.Request) {
	session, err := middleware.GetLoggedSession(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userId := session.UserID
	project, err := database.GetProject(database.DB, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	data := struct {
		message database.Project
	}{
		message: *project,
	}
	json.NewEncoder(w).Encode(data.message)

}
