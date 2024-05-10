package server

import (
	"context"
	"fmt"
	"net/http"
	"xkcd/internal/handlers"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler handlers.HTTPHandler) error {
	s.httpServer = &http.Server{
		Addr: ":" + port,
	}
	handler.Init()
	fmt.Println("Server running at port " + port)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
