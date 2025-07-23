package server

import (
	"log"
	"net/http"
	"session-auth/handler"
	"strings"
	"time"
)

type Server struct {
	srv *http.Server
}

func NewServer(port string, handler *handler.Handler, timeout time.Duration) *Server {
	srv := http.Server{
		Addr:           ":" + port,
		Handler:        handler.Router,
		ReadTimeout:    timeout,
		WriteTimeout:   timeout,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	return &Server{
		srv: &srv,
	}
}

func (s *Server) Run() error {
	log.Printf("Запускаю сервер по адресу http://localhost:%s\n", strings.Split(s.srv.Addr, ":")[1])
	return s.srv.ListenAndServe()
}
func (s *Server) Close() error {
	return s.srv.Close()
}
