package binance

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

// WsPartialBookDepthHandler handle websocket depth event
type WsPartialBookDepthHandler func(event *WsPartialBookDepthEvent)

// WsPartialBookDepthServe serve websocket depth handler with a symbol.
// Valid <levels> are 5, 10, or 20.
func WsPartialBookDepthServe(symbol string, handler WsPartialBookDepthHandler, levels int) (chan struct{}, error) {
	endpoint := fmt.Sprintf("%s/%s@depth%d", baseURL, strings.ToLower(symbol), levels)
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		handler(parseWsPartialBookDepthEvent(message, levels))
	}
	return wsServe(cfg, wsHandler)
}

// WsPartialBookDepthEvent define websocket depth event
type WsPartialBookDepthEvent struct {
	LastUpdateId int64                `json:"lastUpdateId"`
	Bids         [][2]decimal.Decimal `json:"bids"`
	Asks         [][2]decimal.Decimal `json:"asks"`
}
