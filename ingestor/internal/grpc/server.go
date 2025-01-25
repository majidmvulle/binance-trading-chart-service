package grpc

import (
	"fmt"
	aggregatorgrpc "github.com/majidmvulle/binance-trading-chart-service/ingestor/internal/grpc/aggregator"
	"log"
	"net"

	"google.golang.org/grpc"

	aggregatorproto "github.com/majidmvulle/binance-trading-chart-service/ingestor/pkg/api/aggregator"
)
