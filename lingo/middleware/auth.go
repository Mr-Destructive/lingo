package middleware

import (
	"context"
	"lingo/lingo/database"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type authUser struct {
	email        string
	username     string
	passwordHash string
}

var DefaultUserService userService

type userService struct {
}

func (userService) VerifyUser(user database.User) bool {
	dbUser, err := database.UserByEmail(database.DB, user.Email)
	if err != nil {
		return false
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return false
	}
	return true
}

func getPasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func (userService) CreateUser(newUser database.User) error {
	passwordHash, err := getPasswordHash(newUser.Password)
	if err != nil {
		return err
	}
	newAuthUser := authUser{
		email:        newUser.Email,
		username:     newUser.Username,
		passwordHash: passwordHash,
	}
	row := database.DB.QueryRow("SELECT id FROM user ORDER BY id DESC LIMIT 1")
	var lastId int
	err = row.Scan(&lastId)
	if err != nil {
		panic(err)
	}
	newId := lastId + 1
	database.DB.Exec("INSERT INTO user (id, email, username, password) VALUES (?, ?, ?, ?)", newId, newAuthUser.email, newAuthUser.username, newAuthUser.passwordHash)
	return nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := GetLoggedSession(w, r)
		if err == nil {
			ctx := context.WithValue(r.Context(), "user", session.UserID)
			ctx = context.WithValue(r.Context(), "sessionId", session.SessionID)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}
