package main

import (
	"context"
	"fmt"
	"github.com/majidmvulle/binance-trading-chart-service/persistor/config"
	"log"
	"time"

	aggregatorpb "github.com/majidmvulle/binance-trading-chart-service/persistor/internal/clients/aggregator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.Config()

	conn, err := grpc.NewClient(cfg.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	client := aggregatorpb.NewAggregatorServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.StreamCandlesticks(ctx, &aggregatorpb.StreamRequest{})
	if err != nil {
		log.Fatalf("could not stream candlesticks: %v", err)
	}

	log.Println("connected to gRPC server, listening for candlestick stream...")
	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Printf("error receiving from stream: %v", err)
			return
		}

		candleTime := resp.Timestamp.AsTime()

		fmt.Printf("received Candlestick: Symbol=%s, Timestamp=%s, Open=%.2f, High=%.2f, Low=%.2f, Close=%.2f, Volume=%.2f\n",
			resp.Symbol, candleTime.Format(time.RFC3339), resp.Open, resp.High, resp.Low, resp.Close, resp.Volume)
	}
}
