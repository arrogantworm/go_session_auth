package handler

import (
	"net/http"
	"session-auth/models"
)

// api/test
// @Summary Test endpoint
// @Tags test
// @Description Secure endpoint
// @ID test-endpoint
// @Security SessionAuth
// @Produce json
// @Success 200 {object} SuccessRes
// @Failure 401 {object} ErrorRes
// @Router /api/test [get]
func (h *Handler) testEndpoint(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value("session").(*models.Session)
	// log.Println(session.Username)

	// h.sendSuccess(w, "working", http.StatusOK)
	h.sendSuccess(w, session.Username, http.StatusOK)
}
