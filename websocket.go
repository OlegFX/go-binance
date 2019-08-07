package binance

import (
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

func wsServeMax(cfg *wsConfig, handler WsHandler, done chan interface{}) (lsDone chan interface{}, err error) {
	c, _, err := websocket.DefaultDialer.Dial(cfg.endpoint, nil)
	if err != nil {
		return
	}

	if done == nil {
		done = make(chan interface{}, 2)
	}
	lsDone = make(chan interface{}, 2)

	go func() {
		v := <-done
		lsDone <- v
		if v != nil {
			c.Close()
		}
	}()

	go func() {
		defer func() {
			_ = c.Close()

			if pan := recover(); pan != nil {
				go wsServeMax(cfg, handler, done)
				return
			}

			close(lsDone)
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
