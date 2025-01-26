package aggregator

import (
	"github.com/majidmvulle/binance-trading-chart-service/ingestor/internal/services/aggregator"
	aggregatorpb "github.com/majidmvulle/binance-trading-chart-service/ingestor/pkg/api/aggregator"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
)

type Server struct {
	aggregatorpb.UnimplementedAggregatorServiceServer
	candlestickChan chan *aggregator.Candlestick
}

func NewServer(candlestickChan chan *aggregator.Candlestick) *Server {
	return &Server{
		candlestickChan: candlestickChan,
	}
}

func (s *Server) StreamCandlesticks(_ *aggregatorpb.StreamRequest,
	stream aggregatorpb.AggregatorService_StreamCandlesticksServer) error {
	log.Println("client connected for candlestick stream")
	for candle := range s.candlestickChan {
		resp := &aggregatorpb.StreamResponse{
			Symbol:    candle.Symbol,
			Open:      candle.Open,
			High:      candle.High,
			Low:       candle.Low,
			Close:     candle.Close,
			Volume:    candle.Volume,
			Timestamp: timestamppb.New(candle.Timestamp),
		}

		if err := stream.Send(resp); err != nil {
			log.Printf("error sending candlestick to stream: %v", err)

			return err
		}

		log.Printf("Sent candlestick to stream: Symbol=%s, Timestamp=%s, Close=%.2f", candle.Symbol,
			candle.Timestamp.Format(time.RFC3339), candle.Close)
	}
	log.Println("candlestick stream channel closed, ending gRPC stream")
	return nil
}
