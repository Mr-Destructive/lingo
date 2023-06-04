package app

import (
	"database/sql"
	"html/template"
	"lingo/lingo/database"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)


const Store = sessions.NewCookieStore([]byte("secret-key"))

func UserByUsername(db *sql.DB, username string) (*database.User, error) {
	query := "SELECT id, username, password FROM users WHERE username = ?"
	row := db.QueryRow(query, username)

	user := database.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return &user, err
	}

	return &user, nil
}

func CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		templates, err := template.ParseFiles("lingo/templates/login.html")
		if err != nil {
			log.Fatal(err)
		}
		err = templates.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		username := r.FormValue("username")
		password := r.FormValue("password")
		if username == "" || password == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := UserByUsername(database.DB, username)
		if err != nil {
			log.Fatal(err)
		}
		if user.ID == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		session, err := store.Get(r, "session")
		if err != nil {
			log.Fatal(err)
		}
		session.Values["user_id"] = user.ID
		err = session.Save(r, w)
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, "/links", http.StatusFound)
	}
}
