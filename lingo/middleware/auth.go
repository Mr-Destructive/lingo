package middleware

import (
	"context"
	"net/http"
)

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
