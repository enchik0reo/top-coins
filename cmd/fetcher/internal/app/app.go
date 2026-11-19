package app

import (
	"context"
	"encoding/json"
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
	const fnName = "new fethcer app"

	cfg, err := config.NewConfigFromEnv()
	if err != nil {
		log.Fatalf("can't create config from env, func: %s\n", fnName)
	}

	log := logger.NewSlogLogger(cfg.Env)

	app := &App{
		cfg: cfg,
		log: log,
	}

	app.mustPrintCfg()

	fetchers := client.NewCoinsClient(app.cfg.Client.Timeout, app.log).BringFetchers()

	app.infoCollector = service.NewCoinInfoCollector(app.cfg.FetchTicker, fetchers, app.log)

	infoSender, err := sender.NewInfoSender(app.cfg, app.log)
	if err != nil {
		app.log.Error("Can't create info sender", app.log.Attr("error", err), app.log.Attr("func", fnName))
		os.Exit(1)
	}

	app.infoSender = infoSender

	return app
}

func (a *App) MustRun() {
	const fnName = "fetcher app: must run"

	a.log.Info("Fetcher service is running", "env", a.cfg.Env)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	infoCh := a.infoCollector.StartCollect(ctx)
	a.log.Debug("Coin info collector started successfully", a.log.Attr("func", fnName))

	go func() {
		for coinsInfo := range infoCh {
			a.infoSender.Publish(ctx, coinsInfo)
		}
	}()
	a.log.Debug("Coin info sender started successfully", a.log.Attr("func", fnName))

	a.log.Info("Fetcher service started successfully")

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)
	<-stopCh

	a.mustStop()
}

func (a *App) mustStop() {
	const fnName = "fetcher app: must stop"

	a.log.Info("Fetcher service is finishing")

	if err := a.infoSender.Close(); err == nil {
		a.log.Debug("Info sender closed", a.log.Attr("func", fnName))
	}

	a.log.Info("Fetcher stopped gracefully")
}

func (a *App) mustPrintCfg() {
	cfgBytes, err := json.MarshalIndent(a.cfg, "", " ")
	if err != nil {
		a.log.Error("Can't marshal config", a.log.Attr("error", err))
		os.Exit(1)
	}
	a.log.Debug(string(cfgBytes))
}
