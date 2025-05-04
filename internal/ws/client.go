package ws

import (
	"context"
	"encoding/json"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"log"
	"time"
)

const (
	// sendChannelSize controls the max number
	// of messages that can be queued for a client.
	sendChannelSize = 16
	pingPeriod      = (60 * 9 * time.Second) / 10
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

func NewClient(id string, conn *websocket.Conn, manager *Manager) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		ID:      id,
		Conn:    conn,
		Manager: manager,
		send:    make(chan Message, sendChannelSize),
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (c *Client) Start() {
	go c.readPump()
	go c.writePump()
	c.Manager.register <- c
}

func (c *Client) Close() {
	if err := c.Conn.Close(websocket.StatusNormalClosure, "bye :P"); err != nil {
		log.Println("close connection err:", err)
	}
	c.cancel()
}

func (c *Client) Send(msg Message) {
	select {
	case c.send <- msg:
	default:
		c.Manager.forceDisconnect(c)
	}
}

func (c *Client) readPump() {
	defer func() {
		c.Manager.unregister <- c
		c.Close()
	}()

	for {
		var msg Message
		if err := wsjson.Read(c.ctx, c.Conn, &msg); err != nil {
			log.Printf("client %q read message err: %v", c.ID, err)
			break
		}
		c.handleMessage(msg)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				_ = c.Conn.Close(websocket.StatusNormalClosure, "bye :P")
				return
			}
			if err := wsjson.Write(c.ctx, c.Conn, msg); err != nil {
				log.Printf("client %q write message err: %v", c.ID, err)
				return
			}
		case <-ticker.C:
			if err := c.Conn.Ping(c.ctx); err != nil {
				log.Printf("client %q ping err: %v", c.ID, err)
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Client) handleMessage(msg Message) {
	switch msg.Type {
	case "position":
		log.Printf("client %q received position msg: %v", c.ID, msg.Data)
	case "route":
		log.Printf("client %q received route msg: %v", c.ID, msg.Data)
	default:
		log.Printf("client %q received unknwon message type %q", c.ID, msg.Type)
	}
}
