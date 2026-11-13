package processor

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"

	"github.com/top-coins/cmd/processor/internal/kafka"

	"github.com/top-coins/pkg/logger"
	"github.com/top-coins/pkg/models"
)

type saver interface {
	SaveCoinInfo(ctx context.Context, coinsInfo []models.CoinInfo) error
}

type getter interface {
	GetCoinLastInfo(ctx context.Context, name string) (models.CoinInfo, error)
	GetCoinHistory(ctx context.Context, name string, sinceTime time.Time) ([]models.CoinInfo, error)
}

type publisher interface {
	Publish(ctx context.Context, coinInfo models.CoinInfo) error
}

type consumer interface {
	StartConsume(ctx context.Context) (<-chan *kafka.CoinMessage, error)
}

type ProcessorService struct {
	log    *logger.CustomLogger
	stopCh chan struct{}

	saver     saver
	getter    getter
	publisher publisher
	consumer  consumer
}

func NewProcessorService(svr saver, gtr getter, pbl publisher, cnsm consumer, log *logger.CustomLogger) *ProcessorService {
	stopCh := make(chan struct{}, 1)

	return &ProcessorService{
		log:    log,
		stopCh: stopCh,

		saver:     svr,
		getter:    gtr,
		publisher: pbl,
		consumer:  cnsm,
	}
}

func (p *ProcessorService) Run(ctx context.Context) error {
	const fnName = "processor service: Run"

	infoCh, err := p.consumer.StartConsume(ctx)
	if err != nil {
		return errors.Wrap(err, fnName)
	}

	for {
		select {
		case <-ctx.Done():
			p.log.Debug("Context canceled", p.log.Attr("func", fnName))
			break
		case <-p.stopCh:
			p.log.Debug("Service stopped gracefully", p.log.Attr("func", fnName))
			break
		case msg := <-infoCh:
			if err := p.processMessage(ctx, msg); err != nil {
				p.log.Error("Can't process message", p.log.Attr("error", err), p.log.Attr("func", fnName))
			}
		}
	}
}

func (p *ProcessorService) Stop() {
	close(p.stopCh)
}

func (p *ProcessorService) processMessage(ctx context.Context, msg *kafka.CoinMessage) error {
	const fnName = "must proceess message"

	info := msg.GetCoinInfo()

	defer func() {
		if err := msg.Commit(ctx); err != nil {
			p.log.Error("Can't commit message", p.log.Attr("error", err), p.log.Attr("func", fnName))
		}
	}()

	if err := msg.GetError(); err != nil {
		return errors.Wrap(err, fnName)
	}

	durHour := time.Now().Add(-time.Hour)
	hourHistory, err := p.getter.GetCoinHistory(ctx, string(info[0].Name), durHour)
	if err != nil {
		return errors.Wrap(err, fnName)
	}

	if err := p.saver.SaveCoinInfo(ctx, info); err != nil {
		return errors.Wrap(err, fnName)
	}

	go p.notifyIfIncrease(ctx, info, hourHistory)

	return nil
}

func (p *ProcessorService) notifyIfIncrease(ctx context.Context, currentCoinInfo, hourHistory []models.CoinInfo) {
	const fnName = "notify if increase"

	currCoinPrice := coinAvgPrice(currentCoinInfo)
	historyPrice := coinAvgPrice(hourHistory)

	if currCoinPrice > historyPrice {

		coin, err := p.getter.GetCoinLastInfo(ctx, string(currentCoinInfo[0].Name))
		if err != nil {
			p.log.Error("Can't get coin last info", p.log.Attr("error", err), p.log.Attr("func", fnName))

			return
		}

		if err := p.publisher.Publish(ctx, coin); err != nil {
			p.log.Error("Can't publish coin info", p.log.Attr("error", err), p.log.Attr("func", fnName))

			return
		}
	}
}

func coinAvgPrice(coins []models.CoinInfo) uint64 {
	var tempSum uint64
	var currCoinAvg uint64

	for _, coin := range coins {
		tempSum += coin.Price
	}

	currCoinAvg = tempSum / uint64(len(coins))

	return currCoinAvg
}
