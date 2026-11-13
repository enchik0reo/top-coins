package grpcsrv

import (
	"context"

	"github.com/cockroachdb/errors"
	protov1 "github.com/top-coins/pkg/api/proto/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *GRPCServer) GetCoinInfo(ctx context.Context, req *protov1.GetCoinLastInfoRequest) (*protov1.GetCoinLastInfoResponse, error) {
	const fnName = "grpc server: get coin info"

	coin, err := s.storage.GetCoinLastInfo(ctx, req.GetName())
	if err != nil {
		return nil, errors.Wrap(err, fnName)
	}

	pCoin := protov1.Coin{
		Id:          coin.ID,
		RequestId:   coin.RequestID,
		Name:        string(coin.Name),
		Price:       coin.Price,
		Currency:    string(coin.Currency),
		Source:      string(coin.Source),
		RequestTime: timestamppb.New(coin.RequestTime),
	}

	return &protov1.GetCoinLastInfoResponse{
		Coin: &pCoin,
	}, nil
}

func (s *GRPCServer) GetCoinHistory(ctx context.Context, req *protov1.GetCoinHistoryRequest) (*protov1.GetCoinHistoryResponse, error) {
	const fnName = "grpc server: get coin history"

	history, err := s.storage.GetCoinHistory(ctx, req.GetName(), req.GetSinceTime().AsTime())
	if err != nil {
		return nil, errors.Wrap(err, fnName)
	}

	pHistory := make([]*protov1.Coin, 0, len(history))

	for _, coin := range history {
		pHistory = append(pHistory, &protov1.Coin{
			Id:          coin.ID,
			RequestId:   coin.RequestID,
			Name:        string(coin.Name),
			Price:       coin.Price,
			Currency:    string(coin.Currency),
			Source:      string(coin.Source),
			RequestTime: timestamppb.New(coin.RequestTime),
		})
	}

	return &protov1.GetCoinHistoryResponse{
		History: pHistory,
	}, nil
}
