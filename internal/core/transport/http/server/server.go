package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	"go.uber.org/zap"
)

type Server struct {
	config Config
	log    *core_logger.Logger
	server *http.Server
}

func NewServer(
	config Config,
	log *core_logger.Logger,
	handler http.Handler,
) *Server {
	return &Server{
		config: config,
		log:    log,
		server: &http.Server{
			Addr:              config.Addr,
			Handler:           handler,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	baseCtx := core_logger.ToContext(context.WithoutCancel(ctx), s.log)
	s.server.BaseContext = func(net.Listener) context.Context { return baseCtx }

	errChan := make(chan error, 1)
	go func() {
		s.log.Warn("start HTTP server", zap.String("addr", s.config.Addr))
		errChan <- s.server.ListenAndServe()
	}()

	select {
	case err := <-errChan:
		return fmt.Errorf("listen and serve: %w", err)

	case <-ctx.Done():
		s.log.Warn("shutdown HTTP server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()

		if err := s.server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		if err := <-errChan; err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("listen and serve: %w", err)
		}

		s.log.Warn("HTTP server stopped")

		return nil
	}
}
