package handler

import (
	"context"
	"encoding/json"
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
					sendError(w, "user not authorized", http.StatusUnauthorized)
					return
				}
				sendError(w, err.Error(), http.StatusBadRequest)
				return
			}

			sessionToken := c.Value

			// log.Println(sessionToken)
			sessionId, err := uuid.Parse(sessionToken)
			if err != nil {
				// log.Println(err)
				sendError(w, err.Error(), http.StatusBadRequest)
				return
			}
			// log.Println(sessionId)

			session, err := s.ValidateSession(r.Context(), sessionId)
			if err != nil {
				sendError(w, err.Error(), http.StatusUnauthorized)
				return
			}
			// log.Println(session)

			ctx := context.WithValue(r.Context(), "session", session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func sendError(w http.ResponseWriter, message string, status int) error {

	response := ErrorRes{Message: message}

	resp, err := json.Marshal(response)
	if err != nil {
		return err
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(resp)
	if err != nil {
		return err
	}

	return nil
}
