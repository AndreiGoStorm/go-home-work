package internalhttp

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	ep "github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/api/pb/event"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/app"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/config"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/logger"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const tcpPort = "8080"

type Server struct {
	app        *app.App
	logg       *logger.Logger
	address    string
	listener   net.Listener
	gwServer   *http.Server
	middleware *Middleware
}

func NewServer(app *app.App, logg *logger.Logger, conf *config.Config) *Server {
	lis, err := net.Listen("tcp", net.JoinHostPort(conf.HTTPServer.Host, tcpPort))
	if err != nil {
		logg.Error("Failed to start listener: ", err)
	}

	return &Server{
		app:        app,
		logg:       logg,
		address:    net.JoinHostPort(conf.HTTPServer.Host, strconv.Itoa(conf.HTTPServer.Port)),
		listener:   lis,
		middleware: NewMiddleware(logg),
	}
}

func (s *Server) Start(ctx context.Context) error {
	// Create a gRPC server object
	gs := grpc.NewServer()
	ep.RegisterEventServiceServer(gs, s.app)
	go func() {
		s.logg.Info("Serving gRPC on " + s.listener.Addr().String())
		if err := gs.Serve(s.listener); err != nil {
			s.logg.Error("Failed to stop listener: ", err)
		}
	}()

	// Create a client connection to the gRPC server we just started
	conn, err := grpc.NewClient(
		s.listener.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		s.logg.Error("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	err = ep.RegisterEventServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		s.logg.Error("Failed to register gateway:", err)
	}

	s.gwServer = &http.Server{
		Addr:              s.address,
		Handler:           s.middleware.loggingMiddleware(gwmux),
		ReadHeaderTimeout: time.Second * 3,
	}
	s.logg.Info("Serving gRPC-Gateway on " + s.gwServer.Addr)
	if err := s.gwServer.ListenAndServe(); err != nil {
		s.logg.Error("Failed to stop gRPC-Gateway:", err)
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logg.Info("server stopped")
	s.listener.Close()
	s.gwServer.Shutdown(ctx)
	return nil
}
