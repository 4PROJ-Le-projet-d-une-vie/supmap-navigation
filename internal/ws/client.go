package ws

import (
	"context"
	"encoding/json"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"supmap-navigation/internal/navigation"
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
		c.Manager.logger.Warn("failed to close connection", "error", err)
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
			c.Manager.logger.Warn("failed to read message", "clientID", c.ID, "error", err)
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
				c.Manager.logger.Warn("failed to write message", "clientID", c.ID, "error", err)
				return
			}
			c.Manager.logger.Debug("message sent", "clientID", c.ID, "type", msg.Type)
		case <-ticker.C:
			if err := c.Conn.Ping(c.ctx); err != nil {
				c.Manager.logger.Debug("failed to ping client", "clientID", c.ID, "error", err)
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Client) handleMessage(msg Message) {
	switch msg.Type {
	case "init":
		c.Manager.logger.Debug("received init message", "clientID", c.ID, "data", msg.Data)

		var session navigation.Session
		if err := json.Unmarshal(msg.Data, &session); err != nil {
			c.Manager.logger.Warn("failed to unmarshal init message", "clientID", c.ID, "error", err)
			return
		}

		if session.ID != c.ID {
			c.Manager.logger.Warn("Session ID mismatch", "clientID", c.ID, "session", session.ID)
			return
		}

		if err := c.Manager.sessionCache.SetSession(c.ctx, &session); err != nil {
			c.Manager.logger.Warn("failed to cache session", "clientID", c.ID, "error", err)
		}
	case "position":
		c.Manager.logger.Debug("received position message", "clientID", c.ID, "data", msg.Data)

		var pos navigation.Position
		if err := json.Unmarshal(msg.Data, &pos); err != nil {
			c.Manager.logger.Warn("failed to unmarshal position", "clientID", c.ID, "error", err)
			return
		}

		session, err := c.Manager.sessionCache.GetSession(c.ctx, c.ID)
		if err != nil {
			c.Manager.logger.Warn("failed to get session for position update", "clientID", c.ID, "error", err)
			return
		}

		session.LastPosition = pos
		session.UpdatedAt = time.Now()

		if err := c.Manager.sessionCache.SetSession(c.ctx, session); err != nil {
			c.Manager.logger.Warn("failed to update session with new position", "clientID", c.ID, "error", err)
		}
	case "route":
		c.Manager.logger.Debug("received route message", "clientID", c.ID, "data", msg.Data)
	default:
		c.Manager.logger.Debug("received unknown type message", "clientID", c.ID, "type", msg.Type)
	}
}
