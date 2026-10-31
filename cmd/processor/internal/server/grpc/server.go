package grpcsrv

import (
	"context"
	"net"

	"github.com/cockroachdb/errors"
	"github.com/top-coins/cmd/processor/internal/config"
	protov1 "github.com/top-coins/pkg/api/proto/v1"

	"google.golang.org/grpc"
)

type storage interface {
	LastCoin()    // TODO right signature
	CoinHistory() // TODO right signature
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
	const fnName = "server: run"

	if err := s.server.Serve(s.listener); err != nil {
		return errors.Wrap(err, fnName)
	}

	return nil
}

func (s *GRPCServer) Stop() {
	if s.server != nil {
		s.server.GracefulStop()
	}
}

func (s *GRPCServer) GetLastCoin(context.Context, *protov1.GetLastCoinRequest) (*protov1.GetLastCoinResponse, error) {
	// s.fetcher.LastCoin()
	panic("implement me")
}

func (s *GRPCServer) GetCoinHistory(context.Context, *protov1.GetCoinHistoryRequest) (*protov1.GetCoinHistoryResponse, error) {
	// s.fetcher.CoinHistory()
	panic("implement me")
}
