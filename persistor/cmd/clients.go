package main

import (
	"context"
	"fmt"
	"log"

	aggregatorpb "github.com/majidmvulle/binance-trading-chart-service/persistor/internal/clients/aggregator"
	aggtraderepo "github.com/majidmvulle/binance-trading-chart-service/persistor/internal/repository/aggtrade"
	aggtradesvc "github.com/majidmvulle/binance-trading-chart-service/persistor/internal/service/aggtrade"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func RegisterAggregatorClient(ctx context.Context, conn *grpc.ClientConn, db *gorm.DB) error {
	svc := aggtradesvc.NewService(aggtraderepo.NewRepository(db))
	client := aggregatorpb.NewAggregatorServiceClient(conn)

	stream, err := client.StreamCandlesticks(ctx, &aggregatorpb.StreamRequest{})
	if err != nil {
		log.Fatalf("could not stream candlesticks from aggregator service: %v", err)
	}

	errGrp := errgroup.Group{}

	errGrp.Go(func() error {
		return svc.HandleStream(ctx, stream)
	})

	log.Println("connected to gRPC server, listening from aggregator service...")

	if err := errGrp.Wait(); err != nil {
		return fmt.Errorf("aggregator client is failing: %w", err)
	}

	return nil
}
