package aggregator_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/majidmvulle/binance-trading-chart-service/ingestor/internal/clients/binance"
	aggregatorsvc "github.com/majidmvulle/binance-trading-chart-service/ingestor/internal/services/aggregator"
)

func TestAggregator_AggregateTrade_NewCandlestick(t *testing.T) {
	agg := aggregatorsvc.NewAggregator()
	tradeTime := time.Now().UTC().Truncate(time.Minute)

	tradeData := binance.TradeData{
		Symbol:    "BTCUSDT",
		Price:     "100.0",
		Quantity:  "1.0",
		TradeTime: tradeTime.UnixMilli(),
	}

	candle, err := agg.AggregateTrade(tradeData)
	if err != nil {
		t.Fatalf("aggregateTrade failed: %v", err)
	}

	expectedCandle := &aggregatorsvc.Candlestick{
		Symbol:    "BTCUSDT",
		Open:      100.0,
		High:      100.0,
		Low:       100.0,
		Close:     100.0,
		Volume:    1.0,
		Timestamp: tradeTime,
	}

	if !reflect.DeepEqual(candle, expectedCandle) {
		t.Errorf("Aggregated candlestick is incorrect. \ngot: %#v \nwant: %#v", candle, expectedCandle)
	}
}

func TestAggregator_AggregateTrade_UpdateCandlestick(t *testing.T) {
	agg := aggregatorsvc.NewAggregator()
	tradeTime := time.Now().UTC().Truncate(time.Minute)

	tradeData1 := binance.TradeData{
		Symbol:    "BTCUSDT",
		Price:     "100.0",
		Quantity:  "1.0",
		TradeTime: tradeTime.UnixMilli(),
	}
	_, _ = agg.AggregateTrade(tradeData1) // First trade - creates candle

	tradeData2 := binance.TradeData{
		Symbol:    "BTCUSDT",
		Price:     "102.5",
		Quantity:  "0.5",
		TradeTime: tradeTime.UnixMilli(), // Same minute - update candle
	}
	updatedCandle, err := agg.AggregateTrade(tradeData2)
	if err != nil {
		t.Fatalf("AggregateTrade failed: %v", err)
	}

	expectedCandle := &aggregatorsvc.Candlestick{
		Symbol:    "BTCUSDT",
		Open:      100.0,
		High:      102.5,
		Low:       100.0,
		Close:     102.5,
		Volume:    1.5,
		Timestamp: tradeTime,
	}

	if !reflect.DeepEqual(updatedCandle, expectedCandle) {
		t.Errorf("Updated candlestick is incorrect. \ngot: %#v \nwant: %#v", updatedCandle, expectedCandle)
	}
}

func TestAggregator_AggregateTrade_MultipleTradesInMinute(t *testing.T) {
	agg := aggregatorsvc.NewAggregator()
	tradeTime := time.Now().UTC().Truncate(time.Minute)

	trades := []binance.TradeData{
		{Symbol: "BTCUSDT", Price: "100.0", Quantity: "1.0", TradeTime: tradeTime.UnixMilli()},
		{Symbol: "BTCUSDT", Price: "99.5", Quantity: "0.8", TradeTime: tradeTime.UnixMilli()},
		{Symbol: "BTCUSDT", Price: "101.0", Quantity: "1.2", TradeTime: tradeTime.UnixMilli()},
		{Symbol: "BTCUSDT", Price: "100.5", Quantity: "0.5", TradeTime: tradeTime.UnixMilli()},
	}

	var lastCandle *aggregatorsvc.Candlestick
	for _, trade := range trades {
		candle, err := agg.AggregateTrade(trade)
		if err != nil {
			t.Fatalf("aggregateTrade failed for trade %+v: %v", trade, err)
		}
		lastCandle = candle
	}

	expectedCandle := &aggregatorsvc.Candlestick{
		Symbol:    "BTCUSDT",
		Open:      100.0,
		High:      101.0,
		Low:       99.5,
		Close:     100.5,
		Volume:    0.0,
		Timestamp: tradeTime,
	}

	if lastCandle != nil && (lastCandle.Symbol != expectedCandle.Symbol ||
		lastCandle.Open != expectedCandle.Open ||
		lastCandle.High != expectedCandle.High ||
		lastCandle.Low != expectedCandle.Low ||
		lastCandle.Close != expectedCandle.Close ||
		!lastCandle.Timestamp.Equal(expectedCandle.Timestamp)) {
		t.Errorf("Candlestick after multiple trades is incorrect (excluding volume). \ngot: %#v \nwant: %#v",
			lastCandle, expectedCandle)
	}
}

func TestAggregator_AggregateTrade_DifferentSymbols(t *testing.T) {
	agg := aggregatorsvc.NewAggregator()
	tradeTime := time.Now().UTC().Truncate(time.Minute)

	tradeBTC := binance.TradeData{Symbol: "BTCUSDT", Price: "100.0", Quantity: "1.0", TradeTime: tradeTime.UnixMilli()}
	tradeETH := binance.TradeData{Symbol: "ETHUSDT", Price: "50.0", Quantity: "2.0", TradeTime: tradeTime.UnixMilli()}

	candleBTC, errBTC := agg.AggregateTrade(tradeBTC)
	if errBTC != nil {
		t.Fatalf("aggregateTrade failed for BTCUSDT: %v", errBTC)
	}
	candleETH, errETH := agg.AggregateTrade(tradeETH)
	if errETH != nil {
		t.Fatalf("aggregateTrade failed for ETHUSDT: %v", errETH)
	}

	expectedCandleBTC := &aggregatorsvc.Candlestick{Symbol: "BTCUSDT", Open: 100.0, High: 100.0, Low: 100.0,
		Close: 100.0, Volume: 1.0, Timestamp: tradeTime}
	expectedCandleETH := &aggregatorsvc.Candlestick{Symbol: "ETHUSDT", Open: 50.0, High: 50.0, Low: 50.0,
		Close: 50.0, Volume: 2.0, Timestamp: tradeTime}

	if !reflect.DeepEqual(candleBTC, expectedCandleBTC) {
		t.Errorf("aggregated candlestick for BTCUSDT is incorrect. \ngot: %#v \nwant: %#v",
			candleBTC, expectedCandleBTC)
	}
	if !reflect.DeepEqual(candleETH, expectedCandleETH) {
		t.Errorf("aggregated candlestick for ETHUSDT is incorrect. \ngot: %#v \nwant: %#v",
			candleETH, expectedCandleETH)
	}
}

func TestAggregator_AggregateTrade_ZeroQuantityTrade(t *testing.T) {
	agg := aggregatorsvc.NewAggregator()
	tradeTime := time.Now().UTC().Truncate(time.Minute)

	tradeData := binance.TradeData{
		Symbol:    "BTCUSDT",
		Price:     "105.0",
		Quantity:  "0.0", // Zero quantity trade
		TradeTime: tradeTime.UnixMilli(),
	}

	candle, err := agg.AggregateTrade(tradeData)
	if err != nil {
		t.Fatalf("aggregateTrade failed: %v", err)
	}

	expectedCandle := &aggregatorsvc.Candlestick{
		Symbol:    "BTCUSDT",
		Open:      105.0,
		High:      105.0,
		Low:       105.0,
		Close:     105.0,
		Volume:    0.0,
		Timestamp: tradeTime,
	}

	if !reflect.DeepEqual(candle, expectedCandle) {
		t.Errorf("aggregated candlestick for zero quantity trade is incorrect. \ngot: %#v \nwant: %#v",
			candle, expectedCandle)
	}
}

func TestAggregator_AggregateTrade_ZeroPriceTrade(t *testing.T) {
	agg := aggregatorsvc.NewAggregator()
	tradeTime := time.Now().UTC().Truncate(time.Minute)

	tradeData := binance.TradeData{
		Symbol:    "BTCUSDT",
		Price:     "0.0", // Zero price trade
		Quantity:  "1.0",
		TradeTime: tradeTime.UnixMilli(),
	}

	candle, err := agg.AggregateTrade(tradeData)
	if err != nil {
		t.Fatalf("AggregateTrade failed: %v", err)
	}

	expectedCandle := &aggregatorsvc.Candlestick{
		Symbol:    "BTCUSDT",
		Open:      0.0, // Open, High, Low, Close can be 0.0 if first trade is 0 price
		High:      0.0,
		Low:       0.0,
		Close:     0.0,
		Volume:    1.0, // Volume should still accumulate even with 0 price
		Timestamp: tradeTime,
	}

	if !reflect.DeepEqual(candle, expectedCandle) {
		t.Errorf("aggregated candlestick for zero price trade is incorrect. \ngot: %#v \nwant: %#v",
			candle, expectedCandle)
	}
}

func TestAggregator_AggregateTrade_VolumeAccumulation(t *testing.T) {
	agg := aggregatorsvc.NewAggregator()
	tradeTime := time.Date(2025, time.January, 27, 10, 30, 0, 0, time.UTC).
		Truncate(time.Minute)

	trades := []binance.TradeData{
		{Symbol: "BTCUSDT", Price: "100.0", Quantity: "1.0", TradeTime: tradeTime.UnixMilli()}, // Trade 1
		{Symbol: "BTCUSDT", Price: "99.5", Quantity: "0.8", TradeTime: tradeTime.UnixMilli()},  // Trade 2
		{Symbol: "BTCUSDT", Price: "101.0", Quantity: "1.2", TradeTime: tradeTime.UnixMilli()}, // Trade 3
		{Symbol: "BTCUSDT", Price: "100.5", Quantity: "0.5", TradeTime: tradeTime.UnixMilli()}, // Trade 4
	}

	var lastCandle *aggregatorsvc.Candlestick
	for _, trade := range trades {
		candle, err := agg.AggregateTrade(trade)
		if err != nil {
			t.Fatalf("AggregateTrade failed for trade %+v: %v", trade, err)
		}
		lastCandle = candle
	}

	expectedCandle := &aggregatorsvc.Candlestick{
		Symbol:    "BTCUSDT",
		Open:      100.0,     // Price of the first trade
		High:      101.0,     // Highest price among all trades
		Low:       99.5,      // Lowest price among all trades
		Close:     100.5,     // Price of the last trade
		Volume:    3.5,       // Sum of quantities (1.0 + 0.8 + 1.2 + 0.5)
		Timestamp: tradeTime, // Start of the minute
	}

	if !reflect.DeepEqual(lastCandle, expectedCandle) {
		t.Errorf("candlestick after multiple trades is incorrect. \ngot: %#v \nwant: %#v",
			lastCandle, expectedCandle)
	}

	if lastCandle == nil {
		t.Errorf("expected Candlestick to be non-nil")

		return
	}

	if lastCandle.Symbol != expectedCandle.Symbol {
		t.Errorf("symbol mismatch: got %s, want %s", lastCandle.Symbol, expectedCandle.Symbol)
	}
	if lastCandle.Open != expectedCandle.Open {
		t.Errorf("open price mismatch: got %.2f, want %.2f", lastCandle.Open, expectedCandle.Open)
	}
	if lastCandle.High != expectedCandle.High {
		t.Errorf("high price mismatch: got %.2f, want %.2f", lastCandle.High, expectedCandle.High)
	}
	if lastCandle.Low != expectedCandle.Low {
		t.Errorf("low price mismatch: got %.2f, want %.2f", lastCandle.Low, expectedCandle.Low)
	}
	if lastCandle.Close != expectedCandle.Close {
		t.Errorf("close price mismatch: got %.2f, want %.2f", lastCandle.Close, expectedCandle.Close)
	}
	if lastCandle.Volume != expectedCandle.Volume {
		t.Errorf("volume mismatch: got %.2f, want %.2f", lastCandle.Volume, expectedCandle.Volume)
	}
	if !lastCandle.Timestamp.Equal(expectedCandle.Timestamp) {
		t.Errorf("Timestamp mismatch: got %v, want %v", lastCandle.Timestamp, expectedCandle.Timestamp)
	}
}
