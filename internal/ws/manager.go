package ws

import (
	"context"
	"log"
	"sync"
)

type Manager struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan Message
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewManager(ctx context.Context) *Manager {
	ctx, cancel := context.WithCancel(ctx)
	return &Manager{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Message),
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (m *Manager) Start() {
	for {
		select {
		case client := <-m.register:
			m.mu.Lock()
			m.clients[client.ID] = client
			m.mu.Unlock()
			log.Printf("client %q connected", client.ID)
		case client := <-m.unregister:
			m.mu.Lock()
			if _, ok := m.clients[client.ID]; ok {
				delete(m.clients, client.ID)
				close(client.send)
				log.Printf("client %q disconnected", client.ID)
			}
			m.mu.Unlock()
		case message := <-m.broadcast:
			m.mu.RLock()
			for _, client := range m.clients {
				select {
				case client.send <- message:
				default:
					go m.forceDisconnect(client)
				}
			}
			m.mu.RUnlock()
		case <-m.ctx.Done():
			return
		}
	}
}

func (m *Manager) Broadcast(message Message) {
	m.broadcast <- message
}

func (m *Manager) forceDisconnect(c *Client) {
	c.Close()
}

func (m *Manager) Shutdown() {
	m.cancel()
	m.mu.Lock()
	for _, client := range m.clients {
		client.Close()
	}
	m.mu.Unlock()
}
