package httpserver

import (
	"context"
	"net/http"
	"time"
)

// Создание абстракции над структурой сервер
type Server struct {
	httpServer *http.Server
}

// Запуск сервера
func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    5 * time.Minute,
		WriteTimeout:   5 * time.Minute,
	}

	return s.httpServer.ListenAndServe()
}

// Остановка сервера
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
