package middleware

import (
	"context"
	"fmt"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := GetLoggedSession(w, r)
		fmt.Println(session)
		if err == nil {
			ctx := context.WithValue(r.Context(), "user", session.UserID)
			ctx = context.WithValue(r.Context(), "sessionId", session.SessionID)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}
