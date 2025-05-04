package websocket

import (
	"context"
	"encoding/json"
	"github.com/coder/websocket"
)

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type Client struct {
	ID      string
	Conn    *websocket.Conn
	Manager *Manager
	send    chan Message
	ctx     context.Context
	cancel  context.CancelFunc
}
