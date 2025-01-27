package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strings"

	"github.com/gorilla/websocket"
)

type Config struct {
	WebsocketBaseURL string
	Symbols          []string
}

type AggTrade struct {
	Stream string    `json:"stream"`
	Data   TradeData `json:"data"`
}

type TradeData struct {
	EventType     string `json:"e"`
	EventTime     int64  `json:"E"`
	Symbol        string `json:"s"`
	AggTradeID    int64  `json:"a"`
	Price         string `json:"p"`
	Quantity      string `json:"q"`
	FirstTradeID  int64  `json:"f"`
	LastTradeID   int64  `json:"l"`
	TradeTime     int64  `json:"T"`
	IsMarketMaker bool   `json:"m"`
	Ignore        bool   `json:"M"`
}

type Client struct {
	symbols      []string
	websocketURL string
	conn         *websocket.Conn
}

func NewClient(cfg *Config) *Client {
	return &Client{
		symbols:      cfg.Symbols,
		websocketURL: cfg.WebsocketBaseURL,
	}
}

func (c *Client) Connect() error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	log.Printf("connecting to base URL %s", c.websocketURL)

	streamURL, err := url.Parse(c.websocketURL)
	if err != nil {
		return fmt.Errorf("failed to parse websocket url: %w", err)
	}

	streamURL.Path = path.Join(streamURL.Path, "stream")
	symbolsLen := len(c.symbols)
	combinedStreams := ""

	for i, symbol := range c.symbols {
		stream := fmt.Sprintf("%s@aggTrade", strings.ToLower(symbol))
		combinedStreams = fmt.Sprintf("%s%s", combinedStreams, stream)

		if i < symbolsLen-1 {
			combinedStreams = fmt.Sprintf("%s/", combinedStreams)
		}
	}

	query := streamURL.Query()
	query.Set("streams", combinedStreams)
	streamURL.RawQuery = query.Encode()

	conn, _, err := websocket.DefaultDialer.Dial(streamURL.String(), nil)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}

	c.conn = conn

	return nil
}

func (c *Client) ReadAggregatedTicks(ctx context.Context, tradeChan chan<- TradeData) error {
	if c.conn == nil {
		if err := c.Connect(); err != nil {
			return fmt.Errorf("connection error: %w", err)
		}
	}

	defer close(tradeChan)

	for {
		select {
		case <-ctx.Done():
			log.Println("context cancelled, closing websocket")

			return ctx.Err()
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				log.Println("read error:", err)

				return err
			}

			var aggTrade AggTrade

			if err := json.Unmarshal(message, &aggTrade); err != nil {
				log.Printf("error unmarshalling tick data: %v, message: %s", err, string(message))

				continue
			}

			tradeChan <- aggTrade.Data
		}
	}
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}
