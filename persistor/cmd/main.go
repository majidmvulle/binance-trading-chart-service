package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/majidmvulle/binance-trading-chart-service/ingestor/internal/grpc" // Replace with your module path
)

func main() {
	serverAddr := "localhost:50051" // gRPC server address

	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := grpc.NewCandlestickServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Set timeout for client
	defer cancel()

	stream, err := client.StreamCandlesticks(ctx, &grpc.StreamCandlestickRequest{})
	if err != nil {
		log.Fatalf("Could not stream candlesticks: %v", err)
	}

	log.Println("Connected to gRPC server, listening for candlestick stream...")
	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Printf("Error receiving from stream: %v", err)
			return // Exit client on stream error
		}

		candleTime := resp.Timestamp.AsTime() // Convert Timestamp to time.Time

		fmt.Printf("Received Candlestick: Symbol=%s, Timestamp=%s, Open=%.2f, High=%.2f, Low=%.2f, Close=%.2f, Volume=%.2f\n",
			resp.Symbol, candleTime.Format(time.RFC3339), resp.Open, resp.High, resp.Low, resp.Close, resp.Volume)
	}
}
