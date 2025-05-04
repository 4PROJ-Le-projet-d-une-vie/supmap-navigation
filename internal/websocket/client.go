package websocket

import (
	"context"
	"github.com/coder/websocket"
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type Client struct {
	ID      string
	Conn    *websocket.Conn
	Manager *Manager
	Send    chan Message
	cancel  context.CancelFunc
}
