package rest

import (
	"context"
	"errors"
	"github.com/folivorra/get_order/internal/config"
	"log/slog"
	"net/http"
)

type Server struct {
	server *http.Server
	cfg    config.Config
	logger *slog.Logger
}

func NewServer(server *http.Server, cfg config.Config, logger *slog.Logger) *Server {
	return &Server{
		server: server,
		cfg:    cfg,
		logger: logger,
	}
}

func (s *Server) Run() error {
	s.logger.Info("server started",
		slog.String("addr", s.server.Addr),
	)
	
	err := s.server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		s.logger.Warn("server panic",
			slog.String("addr", s.server.Addr),
			slog.String("error", err.Error()),
		)
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, s.cfg.ServerHTTPShutdownTimeout)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		s.logger.Warn("failed to shutdown server",
			slog.String("addr", s.server.Addr),
			slog.String("error", err.Error()),
		)
	}

	s.logger.Info("server has been stopped",
		slog.String("addr", s.server.Addr),
	)
}
