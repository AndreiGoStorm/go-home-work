package internalhttp

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/app"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/config"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/logger"
)

type Server struct {
	app    *app.App
	logg   *logger.Logger
	server *http.Server
}

func NewServer(app *app.App, logg *logger.Logger, conf *config.Config) *Server {
	address := net.JoinHostPort(conf.HTTPServer.Host, strconv.Itoa(conf.HTTPServer.Port))
	return &Server{
		server: &http.Server{Addr: address, ReadHeaderTimeout: time.Second * 3},
		logg:   logg,
		app:    app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("hello execution"))
	})

	s.server.Handler = loggingMiddleware(mux, s.logg)
	s.logg.Info("server started on " + s.server.Addr)
	err := s.server.ListenAndServe()
	if err != nil {
		s.logg.Error(err.Error())
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logg.Info("server stopped")
	s.logg.Close()
	s.server.Shutdown(ctx)
	return nil
}
