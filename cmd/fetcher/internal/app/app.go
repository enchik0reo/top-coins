package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/top-coins/cmd/fetcher/internal/client"
	"github.com/top-coins/cmd/fetcher/internal/config"
	"github.com/top-coins/cmd/fetcher/internal/sender"
	"github.com/top-coins/cmd/fetcher/internal/service"
	"github.com/top-coins/pkg/logger"
	"github.com/top-coins/pkg/models"
)

type coinInfoCollector interface {
	StartCollect(ctx context.Context) chan []models.CoinInfo
}

type coinInfoSender interface {
	Publish(ctx context.Context, coinsInfo []models.CoinInfo) error
	Close() error
}

type App struct {
	cfg *config.FetcherConfig

	infoCollector coinInfoCollector
	infoSender    coinInfoSender

	log *logger.CustomLogger
}

func New() *App {
	const fnName = "new app"

	cfg, err := config.NewConfigFromEnv()
	if err != nil {
		log.Fatalf("can't create config from env, func: %s\n", fnName)
	}

	log := logger.NewSlogLogger(cfg.Env)

	app := &App{
		cfg: cfg,
		log: log,
	}

	fetchers := client.NewCoinsClient(app.cfg.Client.Timeout, app.log).BringFetchers()

	app.infoCollector = service.NewCoinInfoCollector(app.cfg.FetchTicker, fetchers, app.log)

	infoSender, err := sender.NewInfoSender(app.cfg)
	if err != nil {
		app.log.Error("Can't create info sender", app.log.Attr("error", err), app.log.Attr("func", fnName))
		os.Exit(1)
	}

	app.infoSender = infoSender

	return app
}

func (a *App) MustRun() {
	a.log.Info("Fetcher service is running", "env", a.cfg.Env)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	infoCh := a.infoCollector.StartCollect(ctx)

	go func() {
		for coinsInfo := range infoCh {
			a.infoSender.Publish(ctx, coinsInfo)
		}
	}()

	a.log.Info("Fetcher service runs")

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)
	<-stopCh

	a.mustStop()
}

func (a *App) mustStop() {
	a.log.Info("Fetcher service is finishing")

	a.infoSender.Close()

	a.log.Info("Fetcher stopped gracefully")
}
