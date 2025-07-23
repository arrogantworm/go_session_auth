package handler

import (
	"net/http"
	"session-auth/models"
)

func (h *Handler) testEndpoint(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value("session").(*models.Session)
	// log.Println(session.Username)

	// h.sendSuccess(w, "working", http.StatusOK)
	h.sendSuccess(w, session.Username, http.StatusOK)
}
