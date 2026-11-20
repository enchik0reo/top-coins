package grpcsrv

import (
	"context"
	"net"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/top-coins/cmd/processor/internal/config"
	protov1 "github.com/top-coins/pkg/api/proto/v1"
	"github.com/top-coins/pkg/logger"
	"github.com/top-coins/pkg/models"

	"google.golang.org/grpc"
)

var _ protov1.CoinInfoServiceServer = (*GRPCServer)(nil)

type storage interface {
	GetCoinLastInfo(ctx context.Context, name string) (models.CoinInfo, error)
	GetCoinHistory(ctx context.Context, name string, sinceTime time.Time) ([]models.CoinInfo, error)
}

type GRPCServer struct {
	protov1.UnimplementedCoinInfoServiceServer

	log      *logger.CustomLogger
	storage  storage
	server   *grpc.Server
	listener net.Listener
}

func NewServer(cfg *config.ProcessorConfig, stor storage, log *logger.CustomLogger) (*GRPCServer, error) {
	const fnName = "new server"

	lis, err := net.Listen("tcp", cfg.Server.Address)
	if err != nil {
		return nil, errors.Wrap(err, fnName)
	}

	srv := grpc.NewServer()

	coinSrv := &GRPCServer{
		log:      log,
		storage:  stor,
		server:   srv,
		listener: lis,
	}

	protov1.RegisterCoinInfoServiceServer(srv, coinSrv)

	return coinSrv, nil
}

func (s *GRPCServer) Run() error {
	const fnName = "grpc server: run"

	err := s.server.Serve(s.listener)
	if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return errors.Wrap(err, fnName)
	}

	return nil
}

func (s *GRPCServer) Stop(ctx context.Context) error {
	const fnName = "grpc server: stop"

	if s.server == nil && s.listener == nil {
		return nil
	}

	var err error

	if s.server != nil {
		doneCh := make(chan struct{})

		go func() {
			s.server.GracefulStop()
			close(doneCh)
		}()

		select {
		case <-ctx.Done():
			err = errors.Join(err, ctx.Err())
			s.server.Stop()
		case <-doneCh:
		}
	}

	if s.listener != nil {
		errLis := s.listener.Close()
		if errLis != nil {
			err = errors.Join(err, errLis)
		}
	}

	return errors.Wrap(err, fnName)
}
