package binance

import (
	"bytes"
	"encoding/json"
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

		event := &WsPartialBookDepthEvent{}
		message = bytes.Replace(message, []byte(",[]]"), []byte("   ]"), -1)
		err := json.Unmarshal(message, event)
		if err != nil {
			// TODO: callback if there is an error
			fmt.Println(string(message))
			fmt.Println(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler)
}

// WsPartialBookDepthEvent define websocket depth event
type WsPartialBookDepthEvent struct {
	LastUpdateId int64                `json:"lastUpdateId"`
	Bids         [][2]decimal.Decimal `json:"bids"`
	Asks         [][2]decimal.Decimal `json:"asks"`
}
