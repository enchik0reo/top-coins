package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/top-coins/cmd/gateway/internal/config"
	"github.com/top-coins/pkg/logger"
	"github.com/top-coins/pkg/models"
)

type coinClient interface {
	LastCoinInfo(ctx context.Context, id string) (*models.CoinInfo, error)
	CoinHistory(ctx context.Context, id string, since time.Time) ([]*models.CoinInfo, error)
}

type CustomRouter struct {
	*chi.Mux

	coinClient coinClient

	log     *logger.CustomLogger
	timeout time.Duration
}

func NewCustomRouter(c coinClient, cfg config.ApiServer, lg *logger.CustomLogger) http.Handler {
	r := &CustomRouter{
		Mux:        chi.NewRouter(),
		coinClient: c,
		log:        lg,
		timeout:    cfg.Timeout,
	}

	r.Mux.Use(middleware.RequestID)
	r.Mux.Use(middleware.Recoverer)
	r.Mux.Use(middleware.URLFormat)
	r.Mux.Use(middleware.RealIP)
	r.Mux.Use(loggerMw(lg))
	r.Mux.Use(corsSettings(cfg.AllowDomains))

	r.Get("/api/v1/newest_info", r.getLastCoinInfo())
	r.Get("/api/v1/history", r.getCoinHistory())

	return r
}
