package main

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
)

type wsConnection struct {
	wsConn *websocket.Conn
	mu     *sync.Mutex
	logger *logging.Logger
}

func (c *wsConnection) SendTextMessage(payload []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.wsConn.WriteMessage(websocket.TextMessage, payload)
	if err != nil {
		c.logger.Warningf("WS Write error: %v", err)
	}
}

// Close closes the underlying network connection without sending or waiting for a close frame
func (c *wsConnection) Close() {
	c.mu.Lock()
	c.wsConn.Close()
	c.mu.Unlock()
}
