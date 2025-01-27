package aggregator

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/majidmvulle/binance-trading-chart-service/ingestor/internal/clients/binance"
)

// Candlestick represents a 1-minute OHLCV candlestick.
type Candlestick struct {
	Symbol    string    `json:"symbol"`
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    float64   `json:"volume"`
	Timestamp time.Time `json:"timestamp"`
}

// Aggregator manages the aggregation of trade data into candlesticks.
type Aggregator struct {
	candlesticks    map[string]map[string]*Candlestick
	CandlestickChan chan *Candlestick
	lastMinute      map[string]time.Time
}

// NewAggregator creates a new Aggregator instance.
func NewAggregator() *Aggregator {
	return &Aggregator{
		candlesticks:    make(map[string]map[string]*Candlestick),
		CandlestickChan: make(chan *Candlestick),
		lastMinute:      make(map[string]time.Time),
	}
}

// AggregateTrade processes a TradeData and updates/creates a candlestick.
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

	lastMinuteForSymbol := a.lastMinute[trade.Symbol]

	if !minuteStart.Equal(lastMinuteForSymbol) && !lastMinuteForSymbol.IsZero() {
		prevCandleKey := fmt.Sprintf("%s-%s", trade.Symbol, lastMinuteForSymbol.Format(time.RFC3339))
		symbolCandlesticks := a.candlesticks[trade.Symbol]

		completedCandle, ok := symbolCandlesticks[prevCandleKey]
		if ok {
			a.CandlestickChan <- completedCandle
			log.Printf("Completed Candlestick for %s-%s, Close=%.2f, Volume=%.2f", completedCandle.Symbol,
				completedCandle.Timestamp.Format(time.RFC3339), completedCandle.Close, completedCandle.Volume)

			delete(symbolCandlesticks, prevCandleKey)
		}
	}

	a.lastMinute[trade.Symbol] = minuteStart

	symbolCandlesticksMap := a.candlesticks[trade.Symbol]
	if symbolCandlesticksMap == nil {
		symbolCandlesticksMap = make(map[string]*Candlestick)
		a.candlesticks[trade.Symbol] = symbolCandlesticksMap
	}

	candlestickKey := fmt.Sprintf("%s-%s", trade.Symbol, minuteStart.Format(time.RFC3339))

	candle, exists := symbolCandlesticksMap[candlestickKey]
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
		symbolCandlesticksMap[candlestickKey] = candle
		candle.Volume += quantityFloat
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
