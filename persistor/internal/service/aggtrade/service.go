package aggtrade

import (
	"context"
	"fmt"
	aggregatorpb "github.com/majidmvulle/binance-trading-chart-service/persistor/internal/clients/aggregator"
	"github.com/majidmvulle/binance-trading-chart-service/persistor/internal/models"
	"google.golang.org/grpc"
	"log"
)

type aggTradeRepo interface {
	SaveTick(ctx context.Context, tick models.AggTradeTick) error
}

type service struct {
	aggTradeRepo aggTradeRepo
}

func NewService(repo aggTradeRepo) *service {
	return &service{
		aggTradeRepo: repo,
	}
}

func (s *service) HandleStream(ctx context.Context,
	stream grpc.ServerStreamingClient[aggregatorpb.StreamResponse]) error {
	for {
		resp, err := stream.Recv()
		if err != nil {
			return fmt.Errorf("error receiving from stream: %v", err)
		}

		if err := s.aggTradeRepo.SaveTick(ctx, models.AggTradeTick{
			Symbol:    resp.Symbol,
			Open:      resp.Open,
			High:      resp.High,
			Low:       resp.Low,
			Close:     resp.Close,
			Volume:    resp.Volume,
			Timestamp: resp.Timestamp.AsTime(),
		}); err != nil {
			log.Printf("error saving tick: %v", err)
		}
	}
}
