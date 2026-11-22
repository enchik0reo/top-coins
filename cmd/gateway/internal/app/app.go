package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	grpcclient "github.com/top-coins/cmd/gateway/internal/client/grpc"
	"github.com/top-coins/cmd/gateway/internal/config"
	"github.com/top-coins/cmd/gateway/internal/server"
	"github.com/top-coins/cmd/gateway/internal/server/handler"
	"github.com/top-coins/pkg/logger"
)

type App struct {
	cfg config.GatewayConfig
	log *logger.CustomLogger

	client     *grpcclient.Client
	httpServer *server.Server
}

func New() *App {
	const fnName = "new gateway app"

	cfg, err := config.NewConfigFromEnv()
	if err != nil {
		log.Fatalf("Can't create config from env, func: %s\n", fnName)
	}

	log := logger.NewSlogLogger(cfg.Env)

	app := &App{
		cfg: cfg,
		log: log,
	}

	cfgStr, err := app.cfg.PrettyStr()
	if err != nil {
		app.log.Error("Can't create pretty config", app.log.Attr("error", err), app.log.Attr("func", fnName))
		os.Exit(1)
	}
	app.log.Info("Config loaded", app.log.Attr("func", fnName), app.log.Attr("config", cfgStr))

	client, err := grpcclient.NewClient(cfg.GRPCClient.Address, app.log)
	if err != nil {
		app.log.Error("Can't create new client", app.log.Attr("error", err))
		os.Exit(1)
	}

	app.client = client
	app.log.Debug("gRPC client created", app.log.Attr("func", fnName))

	router := handler.NewCustomRouter(client, cfg.ApiServer, app.log)
	srv := server.New(router, &cfg.ApiServer, app.log)
	app.httpServer = srv

	return app
}

func (a *App) MustRun() {
	const fnName = "app: must run"

	a.log.Info("Gateway service is running", "env", a.cfg.Env)

	go func() {
		if err := a.httpServer.Start(); err != nil {
			a.log.Error("Failed over working api sever", a.log.Attr("error", err), a.log.Attr("func", fnName))
			os.Exit(1)
		}
	}()
	a.log.Debug("HTTP server started successfully", a.log.Attr("func", fnName))

	a.log.Info("Gateway service started successfully")

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)
	<-stopCh

	a.mustStop()
}

func (a *App) mustStop() {
	const fnName = "app: must stop"

	a.log.Info("Gateway service is finishing")

	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.CtxTimeout)
	defer cancel()

	if err := a.httpServer.Stop(ctx); err != nil {
		a.log.Error("Failed to close server session", a.log.Attr("error", err), a.log.Attr("func", fnName))
	} else {
		a.log.Debug("HTTP server stopped", a.log.Attr("func", fnName))
	}

	if err := a.client.Close(); err != nil {
		a.log.Error("Failed to close client session", a.log.Attr("error", err), a.log.Attr("func", fnName))
	} else {
		a.log.Debug("gRPC client closed", a.log.Attr("func", fnName))
	}

	a.log.Info("Gateway stopped gracefully")
}
