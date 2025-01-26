package main

import (
	"context"
	"github.com/majidmvulle/binance-trading-chart-service/persistor/config"
	"github.com/majidmvulle/binance-trading-chart-service/persistor/internal/db"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func main() {
	cfg := config.Config()

	dbInstance, err := db.New(context.Background(),
		db.WithReadDSN(cfg.Database.ReadDSN),
		db.WithWriteDSN(cfg.Database.WriteDSN),
	)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	conn, err := grpc.NewClient(cfg.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	errGrp := errgroup.Group{}

	errGrp.Go(func() error {
		return RegisterAggregatorClient(ctx, conn, dbInstance.DB())
	})

	if err := errGrp.Wait(); err != nil {
		log.Fatalf("clients failing: %v", err)
	}
}
