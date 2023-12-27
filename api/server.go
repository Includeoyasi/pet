package api

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	srv http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	srv := &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
