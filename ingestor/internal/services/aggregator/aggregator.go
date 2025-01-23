package aggregator

import (
	"fmt"
	"github.com/majidmvulle/binance-trading-chart-service/ingestor/internal/clients/binance"
	"strconv"
	"time"
)

type Candlestick struct {
	Symbol    string    `json:"symbol"`
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    float64   `json:"volume"`
	Timestamp time.Time `json:"timestamp"`
}

type Aggregator struct {
	candlesticks map[string]*Candlestick
}

func NewAggregator() *Aggregator {
	return &Aggregator{
		candlesticks: make(map[string]*Candlestick),
	}
}

func (a *Aggregator) AggregateTrade(trade binance.TradeData) (*Candlestick, error) {
	priceFloat, err := strconv.ParseFloat(trade.Price, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse price: %w", err)
	}
	quantityFloat, err := strconv.ParseFloat(trade.Quantity, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse quantity: %w", err)
	}

	tradeTime := time.UnixMilli(trade.TradeTime).UTC()
	minuteStart := tradeTime.Truncate(time.Minute)

	candlestickKey := fmt.Sprintf("%s-%s", trade.Symbol, minuteStart.Format(time.RFC3339))

	candle, exists := a.candlesticks[candlestickKey]
	if !exists {
		candle = &Candlestick{
			Symbol:    trade.Symbol,
			Open:      priceFloat,
			High:      priceFloat,
			Low:       priceFloat,
			Close:     priceFloat,
			Volume:    0.0,
			Timestamp: minuteStart,
		}

		a.candlesticks[candlestickKey] = candle
	} else {
		candle.High = maxFloat64(candle.High, priceFloat)
		candle.Low = minFloat64(candle.Low, priceFloat)
		candle.Close = priceFloat
		candle.Volume += quantityFloat
	}

	return candle, nil
}

func maxFloat64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func minFloat64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
