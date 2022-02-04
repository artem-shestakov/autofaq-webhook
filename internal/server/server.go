package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/artem-shestakov/autofaq-webhook/internal/apperror"
	"github.com/gorilla/mux"
)

type Server struct {
	HttpServer http.Server
	Addr       string
	Router     *mux.Router
	Errc       chan *apperror.Error
	Infoc      chan string
}

func NewServer(addr string, errc chan *apperror.Error, infoc chan string) (*Server, *mux.Router) {
	router := mux.NewRouter()
	srv := http.Server{
		Addr:    addr,
		Handler: router,
	}
	return &Server{
		HttpServer: srv,
		Errc:       errc,
		Infoc:      infoc,
	}, router
}

func (s *Server) Run() {
	s.Infoc <- fmt.Sprintf("Server starting on %s", s.HttpServer.Addr)
	fmt.Println(">")
	err := s.HttpServer.ListenAndServe()
	fmt.Println(">>")
	if err != nil {
		s.Errc <- apperror.NewError("Can't start server", err.Error(), "0000", err)
	}
	s.Infoc <- "Server is running"
}

func (s *Server) Stop(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	err := s.HttpServer.Shutdown(ctx)
	if err != nil {
		s.Errc <- apperror.NewError("Shutdown server error", err.Error(), "0000", err)
	}
	cancel()
	s.Infoc <- "Server stoped"
}
