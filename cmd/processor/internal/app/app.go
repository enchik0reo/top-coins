package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/top-coins/cmd/processor/internal/config"
	"github.com/top-coins/cmd/processor/internal/kafka"
	"github.com/top-coins/cmd/processor/internal/processor"
	grpcsrv "github.com/top-coins/cmd/processor/internal/server/grpc"
	"github.com/top-coins/cmd/processor/internal/storage"
	"github.com/top-coins/cmd/processor/internal/storage/cache"
	"github.com/top-coins/cmd/processor/internal/storage/db"
	"github.com/top-coins/pkg/logger"
)

type App struct {
	cfg *config.ProcessorConfig
	log *logger.CustomLogger

	db         *db.PSQLDatabase
	cache      *cache.RedisCache
	storageSvc *storage.StorageService
	kafkaSvc   *kafka.KafkaService
	gRPCSrv    *grpcsrv.GRPCServer
	processor  *processor.ProcessorService
}

func New() *App {
	const fnName = "new processor app"

	app := &App{}

	cfg, err := config.NewConfigFromEnv()
	if err != nil {
		log.Fatalf("Can't create config from env, func: %s\n", fnName)
	}

	app.cfg = cfg
	app.log = logger.NewSlogLogger(cfg.Env)
	app.db = db.NewDatabase(&cfg.Postgres, app.log)
	app.cache = cache.NewCache(cfg.Redis, app.log)
	app.storageSvc = storage.NewStorageService(app.db, app.cache, app.log)
	app.kafkaSvc = kafka.NewKafkaService(&cfg.Kafka, app.log)
	app.processor = processor.NewProcessorService(app.db, app.db, app.kafkaSvc.Producer, app.kafkaSvc.Consumer, app.log)

	gRPCsrv, err := grpcsrv.NewServer(cfg, app.storageSvc)
	if err != nil {
		app.log.Error("Can't create new grpc server", app.log.Attr("error", err), app.log.Attr("func", fnName))
	}

	app.gRPCSrv = gRPCsrv

	return app
}

func (a *App) MustRun() {
	const fnName = "processor app: must run"

	a.log.Info("Processor app is running")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := a.gRPCSrv.Run(); err != nil {
		a.log.Error("Can't run grpc server", a.log.Attr("error", err), a.log.Attr("func", fnName))
		os.Exit(1)
	}

	if err := a.db.CheckConnection(); err != nil {
		a.log.Error("Can't check db connection", a.log.Attr("error", err), a.log.Attr("func", fnName))
		os.Exit(1)
	}

	if err := a.cache.CheckConnection(ctx); err != nil {
		a.log.Error("Can't check cache connection", a.log.Attr("error", err), a.log.Attr("func", fnName))
		os.Exit(1)
	}

	if err := a.processor.Run(ctx); err != nil {
		a.log.Error("Can't run processor", a.log.Attr("error", err), a.log.Attr("func", fnName))
		os.Exit(1)
	}

	a.log.Info("Processor app runs")

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGTERM, syscall.SIGINT)
	<-stopCh

	a.mustStop()
}

func (a *App) mustStop() {
	const fnName = "must stop"

	a.log.Info("Processor app is stopping")

	a.processor.Stop()

	if err := a.gRPCSrv.Stop(); err != nil {
		a.log.Error("Can't stop grpc server", a.log.Attr("error", err), a.log.Attr("func", fnName))
	}

	if err := a.db.Close(); err != nil {
		a.log.Error("Can't close db", a.log.Attr("error", err), a.log.Attr("func", fnName))
	}

	if err := a.cache.Close(); err != nil {
		a.log.Error("Can't close cache", a.log.Attr("error", err), a.log.Attr("func", fnName))
	}

	if err := a.kafkaSvc.Consumer.Close(); err != nil {
		a.log.Error("Can't close kafka consumer", a.log.Attr("error", err), a.log.Attr("func", fnName))
	}

	if err := a.kafkaSvc.Producer.Close(); err != nil {
		a.log.Error("Can't close kafka producer", a.log.Attr("error", err), a.log.Attr("func", fnName))
	}

	a.log.Info("Processor app stopped gracefully")
}
