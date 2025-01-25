package main

import (
	"context"
	"fmt"
	"github.com/majidmvulle/binance-trading-chart-service/ingestor/internal/clients/binance"
	"github.com/majidmvulle/binance-trading-chart-service/ingestor/internal/services/aggregator"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/majidmvulle/binance-trading-chart-service/ingestor/config"
)

func main() {
	cfg := config.Config()

	client := binance.NewClient(&binance.Config{
		WebsocketBaseURL: cfg.Binance.WebsocketBaseURL,
		Symbols:          cfg.Binance.Symbols,
	})

	aggregatorSvc := aggregator.NewAggregator()
	grpcServer := NewGrpcServer(WithCandlestickChan(aggregatorSvc.CandlestickChan))

	if err := client.Connect(); err != nil {
		log.Fatalf("failed to connect to Binance WebSocket: %v", err)
	}

	defer func(client *binance.Client) {
		err := client.Close()
		if err != nil {
			log.Fatalf("failed to close Binance WebSocket: %v", err)
		}
	}(client)

	tradeChan := make(chan binance.TradeData)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := client.ReadAggregatedTicks(ctx, tradeChan); err != nil {
			log.Printf("error reading aggregated trades: %v", err)
			cancel()
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := grpcServer.StartGRPCServer(cfg.App.GrpcPort); err != nil {
			log.Fatalf("failed to start gRPC server: %v", err)
		}
	}()
	defer func() {
		grpcServer.GracefulStop()
	}()
	log.Println("listening for aggTrades...")

	for {
		select {
		case tick := <-tradeChan:
			candle, err := aggregatorSvc.AggregateTrade(tick)
			if err != nil {
				log.Printf("error aggregating trade: %v", err)

				continue
			}

			if cfg.App.Debug {
				fmt.Printf("Candlestick updated: Symbol=%s, Timestamp=%s, Open=%.2f, High=%.2f, "+
					"Low=%.2f, Close=%.2f, Volume=%.2f\n",
					candle.Symbol, candle.Timestamp.Format(time.RFC3339), candle.Open, candle.High, candle.Low,
					candle.Close, candle.Volume)
			}

		case <-interrupt:
			log.Println("interrupt, shutting down...")
			cancel()
			time.Sleep(time.Second)
			return
		case <-ctx.Done():
			log.Println("context done, exiting...")
			return
		}
	}
}
