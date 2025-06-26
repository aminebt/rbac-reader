package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
)

type Server struct {
	logger   *slog.Logger
	listener net.Listener
	server   *http.Server
}

func NewServer(logger *slog.Logger) (*Server, error) {
	port := 9090 // TBD read from config
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %w", err)
	}

	svc := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		//Handler:           rbacapi.NewRouter(api), // api.Router() is a function that returns an http.Handler (implemented with chi router)
		IdleTimeout:       90 * time.Second, // matches http.DefaultTransport keep-alive timeout
		ReadTimeout:       32 * time.Second,
		ReadHeaderTimeout: 32 * time.Second,
		WriteTimeout:      32 * time.Second,
		ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	svr := &Server{
		logger:   logger,
		listener: listener,
		server:   svc,
	}

	return svr, nil
}

func (s *Server) SetHandler(handler http.Handler) {
	s.server.Handler = handler
}

func (s *Server) SetupWithManager(mgr ctrl.Manager) error {
	err := mgr.Add(s)
	if err != nil {
		return err
	}
	return nil
}

// Start starts the server
// Blocks until the context is cancelled
func (s *Server) Start(ctx context.Context) error {
	serverShutdown := make(chan struct{})
	go func() {
		<-ctx.Done()
		s.logger.Info("shutting down server")
		if err := s.server.Shutdown(context.Background()); err != nil {
			s.logger.Error("error shutting down server", "error", err)
		}
		close(serverShutdown)
	}()

	s.logger.Info("starting server", "addr", s.server.Addr)
	if err := s.server.Serve(s.listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	<-serverShutdown
	return nil
}

// NeedLeaderElection returns false because this app does not need leader election from the manager
func (s *Server) NeedLeaderElection() bool {
	return false
}
