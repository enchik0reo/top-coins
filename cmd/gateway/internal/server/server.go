package server

import (
	"context"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/top-coins/cmd/gateway/internal/config"
	"github.com/top-coins/pkg/logger"
)

type Server struct {
	log     *logger.CustomLogger
	httpSrv *http.Server
}

func New(handler http.Handler, cfg *config.ApiServer, l *logger.CustomLogger) *Server {
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      handler,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &Server{
		log:     l,
		httpSrv: srv,
	}
}

func (s *Server) Start() error {
	const fnName = "server: start"

	if err := s.httpSrv.ListenAndServe(); err != nil {
		return errors.Wrap(err, fnName)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	const fnName = "server: stop"

	if err := s.httpSrv.Shutdown(ctx); err != nil {
		return errors.Wrap(err, fnName)
	}

	return nil
}
