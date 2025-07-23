package handler

import (
	"context"
	"net/http"
	"session-auth/service"

	"github.com/google/uuid"
)

func AuthMiddleware(s *service.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			c, err := r.Cookie("session_token")
			// log.Println(c)
			if err != nil {
				// log.Println(err)
				if err == http.ErrNoCookie {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			sessionToken := c.Value

			// log.Println(sessionToken)
			sessionId, err := uuid.Parse(sessionToken)
			if err != nil {
				// log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			// log.Println(sessionId)

			session, err := s.ValidateSession(r.Context(), sessionId)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// log.Println(session)

			ctx := context.WithValue(r.Context(), "session", session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
