package main

import (
	"fmt"
	aggregatorgrpc "github.com/majidmvulle/binance-trading-chart-service/ingestor/internal/grpc/aggregator"
	"github.com/majidmvulle/binance-trading-chart-service/ingestor/internal/services/aggregator"
	aggregatorproto "github.com/majidmvulle/binance-trading-chart-service/ingestor/pkg/api/aggregator"
	"google.golang.org/grpc"
	"log"
	"net"
)

type options struct {
	candlestickChan chan *aggregator.Candlestick
}

type Option func(o *options)

type Server struct {
	grpcServer *grpc.Server
	options    *options
}

func NewGrpcServer(opts ...Option) *Server {
	opt := options{}

	for _, o := range opts {
		o(&opt)
	}

	return &Server{
		grpcServer: grpc.NewServer(),
	}
}

func WithCandlestickChan(candlestickChan chan *aggregator.Candlestick) Option {
	return func(o *options) {
		o.candlestickChan = candlestickChan
	}
}

func (s *Server) StartGRPCServer(port uint16) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()

	if s.options.candlestickChan != nil {
		aggregatorproto.RegisterAggregatorServiceServer(grpcServer, aggregatorgrpc.NewServer(s.options.candlestickChan))
	}

	log.Printf("gRPC server listening on port :%d", port)

	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	return nil
}

func (s *Server) GracefulStop() {
	log.Println("stopping gRPC server gracefully...")
	s.grpcServer.GracefulStop()
	log.Println("gRPC server stopped.")
}
