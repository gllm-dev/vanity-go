package rest

import (
	"context"
	"errors"
	"fmt"
	"go.gllm.dev/vanity-go/internal/adapters/handlers/rest/gohdl"
	"go.gllm.dev/vanity-go/internal/adapters/handlers/rest/healthzhdl"
	"go.gllm.dev/vanity-go/internal/services/gosvc"
	"log/slog"
	"net/http"
)

// Server represents the HTTP server for serving Go vanity import paths.
type Server struct {
	// server is the underlying HTTP server instance.
	server *http.Server
	// config holds the server configuration.
	config *Config
	// svc is the service that provides the logic for handling requests.
	svc *gosvc.Service
}

// New creates a new Server instance with the provided configuration and service.
func New(
	cfg *Config,
	svc *gosvc.Service,
) *Server {
	return &Server{
		server: new(http.Server),
		config: cfg,
		svc:    svc,
	}
}

// Start starts the HTTP server and listens for incoming requests on the configured port.
func (s *Server) Start(ctx context.Context) error {
	goHdl := gohdl.New(s.svc)
	hlz := healthzhdl.New()
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", hlz.Healthz)
	mux.HandleFunc("/", goHdl.Handle)

	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.Port),
		Handler:      mux,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		IdleTimeout:  s.config.IdleTimeout,
	}

	slog.InfoContext(ctx, "Starting HTTP server", slog.Int("port", s.config.Port))
	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.ErrorContext(ctx, "failed to start HTTP server", slog.String("error", err.Error()))
		return err
	}

	slog.InfoContext(ctx, "HTTP server closed")
	return nil
}

// Stop gracefully shuts down the HTTP server with the provided context.
func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
