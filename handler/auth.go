package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	custom_errors "session-auth/errors"
	"session-auth/models"
)

// @Summary SignUp
// @Tags auth
// @Description Registration
// @ID auth-signup
// @Accept json
// @Param input body models.Credentials true "Credentials"
// @Produce json
// @Success 201 {object} SuccessRes
// @Failure 400,500 {object} ErrorRes
// @Router /auth/signup [post]
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
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

// @Summary SignIn
// @Tags auth
// @Description Login
// @ID auth-signin
// @Accept json
// @Param input body models.Credentials true "Credentials"
// @Produce json
// @Success 200 {object} SuccessRes
// @Failure 400,500 {object} ErrorRes
// @Router /auth/signin [post]
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
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

// auth/logout
// @Summary LogOut
// @Tags auth
// @Description Logout
// @ID auth-logout
// @Security SessionAuth
// @Produce json
// @Success 200 {object} SuccessRes
// @Failure 401,500 {object} ErrorRes
// @Router /auth/logout [post]
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
