package main

import (
	"fmt"
	"github.com/majidmvulle/binance-trading-chart-service/ingestor/internal/grpc/aggregator"
	aggregatorsvc "github.com/majidmvulle/binance-trading-chart-service/ingestor/internal/services/aggregator"
	aggregatorpb "github.com/majidmvulle/binance-trading-chart-service/ingestor/pkg/api/aggregator"
	"google.golang.org/grpc"
	"log"
	"net"
)

type options struct {
	candlestickChan chan *aggregatorsvc.Candlestick
}

type Option func(o *options)

type ServerWrapper struct {
	grpcServer *grpc.Server
	options    *options
}

func NewGrpcServer(opts ...Option) *ServerWrapper {
	opt := options{}

	for _, o := range opts {
		o(&opt)
	}

	return &ServerWrapper{ // Return ServerWrapper
		grpcServer: grpc.NewServer(),
		options:    &opt,
	}
}

func WithCandlestickChan(candlestickChan chan *aggregatorsvc.Candlestick) Option {
	return func(o *options) {
		o.candlestickChan = candlestickChan
	}
}

func (s *ServerWrapper) StartGRPCServer(port uint16) error { // Return *grpc.Server and error
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	if s.options.candlestickChan != nil {
		aggregatorpb.RegisterAggregatorServiceServer(s.grpcServer, aggregator.NewServer(s.options.candlestickChan))
	}

	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve gRPC: %v", err)
	}

	log.Printf("gRPC serving on port :%d", port)

	return nil
}

func (s *ServerWrapper) GracefulStop() {
	log.Println("stopping gRPC server gracefully...")
	s.grpcServer.GracefulStop()
	log.Println("gRPC server stopped.")
}
