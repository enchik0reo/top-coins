package grpcsrv

import (
	"context"
	"net"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/top-coins/cmd/processor/internal/config"
	protov1 "github.com/top-coins/pkg/api/proto/v1"
	"github.com/top-coins/pkg/models"

	"google.golang.org/grpc"
)

type storage interface {
	GetCoinLastInfo(ctx context.Context, name string) (models.CoinInfo, error)
	GetCoinHistory(ctx context.Context, name string, sinceTime time.Time) ([]models.CoinInfo, error)
}

type GRPCServer struct {
	protov1.UnimplementedCoinInfoServiceServer

	storage  storage
	server   *grpc.Server
	listener net.Listener
}

func NewServer(cfg *config.ProcessorConfig, stor storage) (*GRPCServer, error) {
	const fnName = "new server"

	lis, err := net.Listen("tcp", cfg.Server.Port)
	if err != nil {
		return nil, errors.Wrap(err, fnName)
	}

	srv := grpc.NewServer()

	coinSrv := &GRPCServer{
		storage:  stor,
		server:   srv,
		listener: lis,
	}

	protov1.RegisterCoinInfoServiceServer(srv, coinSrv)

	return coinSrv, nil
}

func (s *GRPCServer) Run() error {
	const fnName = "grpc server: run"

	if err := s.server.Serve(s.listener); err != nil {
		return errors.Wrap(err, fnName)
	}

	return nil
}

func (s *GRPCServer) Stop() error {
	const fnName = "grpc server: stop"

	if s.server != nil {
		s.server.GracefulStop()
	}

	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			return errors.Wrap(err, fnName)
		}
	}

	return nil
}
