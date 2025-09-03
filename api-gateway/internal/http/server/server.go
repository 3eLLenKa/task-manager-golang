package server

import (
	"context"
	"fmt"
	"net/http"
)

type Server struct {
	Server *http.Server
}

func New(addr string, handler http.Handler) *Server {
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return &Server{
		Server: server,
	}
}

func (s *Server) MustRun() {
	if err := s.Run(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func (s *Server) Run() error {
	if err := s.Server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) {
	if err := s.Server.Shutdown(ctx); err != nil {
		fmt.Print(err)
	}
}
