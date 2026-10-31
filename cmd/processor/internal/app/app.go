package app

import (
	grpcsrv "github.com/top-coins/cmd/processor/internal/server/grpc"
)

type App struct {
	// log
	// cfg

	server *grpcsrv.GRPCServer
	// kafka
	// db
	// redis

	// business logic
}

func New() *App {

	// log

	// cfg, err := config.NewConfigFromEnv()

	// db

	// srv, err := grpcsrv.NewServer(cfg)
	// if err != nil {
	// 	// return nil, err
	// }

	return &App{}
}

func (a *App) MustRun() {

}

func (a *App) mustStop() {

}
