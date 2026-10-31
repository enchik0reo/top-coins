package grpcclient

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	protov1 "github.com/top-coins/pkg/api/proto/v1"
	"github.com/top-coins/pkg/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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

// Close закрывает соединение (удобно вызывать через defer)
func (c *Client) Close() error {
	const fnName = "client: close"

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return errors.Wrap(err, fnName)
		}
	}
	return nil
}

// GetLastCoin скрывает детали proto и возвращает доменную модель
func (c *Client) GetLastCoin(ctx context.Context, id string) (*models.CoinInfo, error) {
	panic("implement me!")
	// req := &protov1.GetLastCoinRequest{
	// 	Id: id,
	// }

	// resp, err := c.grpcClient.GetLastCoin(ctx, req)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "grpc call GetLastCoin failed")
	// }

	// pbCoin := resp.GetCoin()
	// if pbCoin == nil {
	// 	return nil, fmt.Errorf("empty coin received")
	// }

	// // Маппинг из proto в доменную модель
	// return &models.CoinInfo{
	// 	ID:        pbCoin.GetId(),
	// 	Name:      models.CoinName(pbCoin.GetName()),
	// 	Price:     pbCoin.GetPrice(),
	// 	Currency:  models.CoinCurrency(pbCoin.GetCurrency()),
	// 	Source:    models.CoinSource(pbCoin.GetSource()),
	// 	CreatedAt: pbCoin.GetCreatedAt().AsTime(),
	// }, nil
}

// GetCoinHistory аналогично инкапсулирует историю
func (c *Client) GetCoinHistory(ctx context.Context, id string, since time.Time) ([]*models.CoinInfo, error) {
	panic("implement me!")
	// 	req := &protov1.GetCoinHistoryRequest{
	// 		Id:        id,
	// 		SinceTime: timestamppb.New(since),
	// 	}

	// return nil, nil
}
