package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	custom_errors "session-auth/errors"
	"session-auth/models"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		var badRequestError *custom_errors.BadRequestError
		if errors.As(err, &badRequestError) {
			h.sendError(w, err.Error(), http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

		return
	}

	ctx := r.Context()

	if err := h.service.RegisterUser(ctx, &credentials); err != nil {
		var badRequestError *custom_errors.BadRequestError
		if errors.As(err, &badRequestError) {
			h.sendError(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.sendSuccess(w, "user successfully registered", http.StatusCreated)
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		var badRequestError *custom_errors.BadRequestError
		if errors.As(err, &badRequestError) {
			h.sendError(w, err.Error(), http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

		return
	}

	ctx := r.Context()

	session, err := h.service.Authenticate(ctx, &credentials)
	if err != nil {
		var badRequestError *custom_errors.BadRequestError
		if errors.As(err, &badRequestError) {
			h.sendError(w, err.Error(), http.StatusBadRequest)
		} else {
			h.sendError(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.Token,
		Expires:  session.ExpiresAt,
		Path:     "/",
		HttpOnly: true,
	})

	h.sendSuccess(w, "user logged in", http.StatusOK)
}

func (h *Handler) LogOut(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value("session").(*models.Session)
	if session.Username == "" {
		h.sendError(w, "user logged out", http.StatusBadRequest)
		return
	}

	if err := h.service.Logout(r.Context(), session.Token); err != nil {
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendSuccess(w, "user logged out", http.StatusOK)
}
