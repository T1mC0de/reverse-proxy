package server

import (
	"build-your-own-reverse-proxy/src/internal/config"
	"net/http"
	"context"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    cfg.Server.Port,
			Handler: handler,
		},
	}
}

func (s *Server) StartServer() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) StopServer(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}