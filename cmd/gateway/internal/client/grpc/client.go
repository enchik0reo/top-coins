package grpcclient

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	protov1 "github.com/top-coins/pkg/api/proto/v1"
	"github.com/top-coins/pkg/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var ErrEmptyResponce = fmt.Errorf("responce is empty")

type Client struct {
	conn       *grpc.ClientConn
	grpcClient protov1.CoinInfoServiceClient
}

func NewClient(addr string) (*Client, error) {
	const fnName = "new client"

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrap(err, fnName)
	}

	return &Client{
		conn:       conn,
		grpcClient: protov1.NewCoinInfoServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	const fnName = "client: close"

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return errors.Wrap(err, fnName)
		}
	}
	return nil
}

func (c *Client) CoinLastInfo(ctx context.Context, name string) (models.CoinInfo, error) {
	const fnName = "client: get last coin"

	req := &protov1.GetCoinLastInfoRequest{
		Name: name,
	}

	resp, err := c.grpcClient.GetCoinLastInfo(ctx, req)
	if err != nil {
		return models.CoinInfo{}, errors.Wrap(err, fnName)
	}

	pbCoin := resp.GetCoin()
	if pbCoin == nil {
		return models.CoinInfo{}, errors.Wrap(ErrEmptyResponce, fnName)
	}

	return models.CoinInfo{
		ID:          pbCoin.GetId(),
		Name:        models.CoinName(pbCoin.GetName()),
		Price:       pbCoin.GetPrice(),
		Currency:    models.CoinCurrency(pbCoin.GetCurrency()),
		Source:      models.CoinSource(pbCoin.GetSource()),
		RequestTime: pbCoin.GetRequestTime().AsTime(),
	}, nil
}

func (c *Client) CoinHistory(ctx context.Context, name string, since time.Time) ([]models.CoinInfo, error) {
	const fnName = "client: coin history"

	req := &protov1.GetCoinHistoryRequest{
		Name:      name,
		SinceTime: timestamppb.New(since),
	}

	resp, err := c.grpcClient.GetCoinHistory(ctx, req)
	if err != nil {
		return nil, err
	}

	pbCoinHistory := resp.GetHistory()
	if pbCoinHistory == nil {
		return nil, errors.Wrap(ErrEmptyResponce, fnName)
	}

	res := make([]models.CoinInfo, 0, len(pbCoinHistory))

	for _, pbCoin := range pbCoinHistory {
		res = append(res, models.CoinInfo{
			ID:          pbCoin.GetId(),
			RequestID:   pbCoin.GetRequestId(),
			Name:        models.CoinName(pbCoin.GetName()),
			Price:       pbCoin.GetPrice(),
			Currency:    models.CoinCurrency(pbCoin.GetCurrency()),
			Source:      models.CoinSource(pbCoin.GetSource()),
			RequestTime: pbCoin.GetRequestTime().AsTime(),
		})
	}

	return res, nil
}
