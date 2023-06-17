package middleware

import (
	"database/sql"
	"lingo/lingo/database"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func GetSession(w http.ResponseWriter, r *http.Request, userID int) (database.Session, error) {
	session := NewSession(w, r, userID)
	cookie := http.Cookie{
		Name:  "lingo_session",
		Value: session.SessionID,
		Path:  "/",
	}
	http.SetCookie(w, &cookie)
	return session, nil
}

func GetLoggedSession(w http.ResponseWriter, r *http.Request) (database.Session, error) {
	cookie, err := r.Cookie("lingo_session")
	if err != nil {
        if err == http.ErrNoCookie {
            return database.Session{}, nil
        }
		return database.Session{}, err
	}
	sessionID := cookie.Value
	session, err := GetSessionByID(sessionID)
	if err != nil {
		return database.Session{}, err
	}
	cookie = &http.Cookie{
		Name:     "lingo_session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	}
	http.SetCookie(w, cookie)
	return session, nil
}

const sessionIDChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
const sessionIDLength = 32

func NewSession(w http.ResponseWriter, r *http.Request, userID int) database.Session {
	db := database.DB
	var sessionID []byte
	row := db.QueryRow("SELECT session_id FROM sessions WHERE user_id = ?", userID)
	err := row.Scan(&sessionID)
	if err == sql.ErrNoRows {
		sessionID = make([]byte, sessionIDLength)
		source := rand.NewSource(time.Now().UnixNano())
		rand_source := rand.New(source)
		for i := range sessionID {
			sessionID[i] = sessionIDChars[rand_source.Intn(len(sessionIDChars))]
		}
		stmt, err := db.Prepare("INSERT INTO sessions (id, user_id, session_id) VALUES (?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		_, err = stmt.Exec(userID, userID, string(sessionID))
		if err != nil {
			log.Fatal(err)
		}
	}
	cookie := http.Cookie{
		Name:     "lingo_session",
		Value:    string(sessionID),
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	}
	http.SetCookie(w, &cookie)

	session := database.Session{
		UserID:    userID,
		SessionID: string(sessionID),
	}
	return session
}

func DeleteSession(w http.ResponseWriter, r *http.Request, sessionID string) {
	db := database.DB
	stmt, err := db.Prepare("DELETE FROM sessions WHERE session_id = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(sessionID)
	if err != nil {
		log.Fatal(err)
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "lingo_session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}

func GetSessionByID(sessionID string) (database.Session, error) {
	db := database.DB
	var userID int
	row := db.QueryRow("SELECT user_id FROM sessions WHERE session_id = ?", sessionID)
	err := row.Scan(&userID)
	session := database.Session{
		UserID:    userID,
		SessionID: sessionID,
	}
	if err != nil {
		return database.Session{}, err
	}
	return session, nil
}
