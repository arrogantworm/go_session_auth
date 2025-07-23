package handler

import (
	"encoding/json"
	"net/http"
	"session-auth/service"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	service *service.Service
	Router  *chi.Mux
}

func NewHandler(service *service.Service) *Handler {
	handler := &Handler{
		service: service,
	}
	handler.NewRouter()
	return handler
}

func (h *Handler) NewRouter() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Timeout(60 * time.Second))

	// Basic CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Route("/api", func(r chi.Router) {
		r.Use(AuthMiddleware(h.service))
		r.Get("/test", h.testEndpoint)
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", h.SignUp)
		r.Post("/signin", h.SignIn)
		r.With(AuthMiddleware(h.service)).Post("/logout", h.LogOut)
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8000/swagger/doc.json"),
	))

	h.Router = r
}

// Common responses

func (h *Handler) sendSuccess(w http.ResponseWriter, message string, status int) error {

	response := SuccessRes{Message: message}

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

func (h *Handler) sendError(w http.ResponseWriter, message string, status int) error {

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
