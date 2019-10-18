package binance

import (
	"time"
	"strings"

	"github.com/gorilla/websocket"
)

// WsHandler handle raw websocket message
type WsHandler func(message []byte)

type wsConfig struct {
	endpoint string
}

func newWsConfig(endpoint string) *wsConfig {
	return &wsConfig{
		endpoint: endpoint,
	}
}

type wsServeFunc func(*wsConfig, WsHandler) (chan struct{}, error)

func wsServe(cfg *wsConfig, handler WsHandler) (done chan struct{}, err error) {
	c, _, err := websocket.DefaultDialer.Dial(cfg.endpoint, nil)
	if err != nil {
		return
	}
	done = make(chan struct{}, 2)
	go func() {
		defer func() {
			_ = c.Close()
			close(done)
		}()

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				return
			}
			go handler(message)
		}
	}()

	return
}

func wsServeMax(cfg *wsConfig, handler WsHandler, done chan interface{}, closed *bool) (err error) {
	c, _, err := websocket.DefaultDialer.Dial(cfg.endpoint, nil)
	if err != nil {
		return
	}

	if done == nil {
		done = make(chan interface{}, 2)
	}

	go func() {
		defer func() {
			_ = c.Close()

			if pan := recover(); pan != nil {
				go wsServeMax(cfg, handler, done, closed)
				return
			}

			select {
			case _, ok := <-done:
				if !ok {
					break
				}
			default:
				close(done)
			}
		}()

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					time.Sleep(time.Second * 3)
					go wsServeMax(cfg, handler, done, closed)
				} else {
					time.Sleep(time.Second * 3)
					go wsServeMax(cfg, handler, done, closed)
				}
				return
			}

			if *closed {
				return
			}

			handler(message)
		}
	}()

	return
}
