package app

import (
	"math/rand"
	"net/http"
)

func randomString(n int) string {
    var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
    s := make([]rune, n)
    for i := range s {
        s[i] = letters[rand.Intn(len(letters))]
    }
    return string(s)
}

func SessionStart(w http.ResponseWriter, r *http.Request) (sessionID string) {
	// Get session ID from cookie
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		// No session ID - generate one
		sessionID = randomString(32)
		cookie = &http.Cookie{
			Name:  "session",
			Value: sessionID,
		}
		http.SetCookie(w, cookie)
	} else {
		sessionID = cookie.Value
	}
	return
}

func SessionGet(r *http.Request) map[string]interface{} {
	sessionID := SessionStart(w, r)
	if sessions[sessionID] == nil {
		sessions[sessionID] = make(map[string]interface{})
	}
	return sessions[sessionID]
}

func SessionSave(w http.ResponseWriter, r *http.Request) {
	sessionID := SessionStart(w, r)
	cookie := &http.Cookie{
		Name:  "session",
		Value: sessionID,
	}
	http.SetCookie(w, cookie)
}

func SessionDestroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == nil {
		sessionID := cookie.Value
		delete(sessions, sessionID)
		cookie = &http.Cookie{
			Name:   "session",
			Value:  "",
			MaxAge: -1,
		}
		http.SetCookie(w, cookie)
	}
}
