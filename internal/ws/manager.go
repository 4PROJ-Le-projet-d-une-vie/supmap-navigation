package ws

import (
	"context"
	"github.com/coder/websocket"
	"log/slog"
	"supmap-navigation/internal/navigation"
	"sync"
)

type Manager struct {
	clients      map[string]*Client
	register     chan *Client
	unregister   chan *Client
	broadcast    chan Message
	mu           sync.RWMutex
	ctx          context.Context
	cancel       context.CancelFunc
	logger       *slog.Logger
	sessionCache navigation.SessionCache
}

func NewManager(ctx context.Context, logger *slog.Logger, cache navigation.SessionCache) *Manager {
	ctx, cancel := context.WithCancel(ctx)
	return &Manager{
		clients:      make(map[string]*Client),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		broadcast:    make(chan Message),
		ctx:          ctx,
		cancel:       cancel,
		logger:       logger,
		sessionCache: cache,
	}
}

func (m *Manager) Start() {
	defer m.Shutdown()
	m.logger.Info("Websocket manager is running")
	for {
		select {
		case client := <-m.register:
			m.mu.Lock()
			m.clients[client.ID] = client
			m.mu.Unlock()
			m.logger.Debug("client connected", "clientID", client.ID)
		case client := <-m.unregister:
			m.mu.Lock()
			if _, ok := m.clients[client.ID]; ok {
				delete(m.clients, client.ID)
				close(client.send)
				m.logger.Debug("client disconnected", "clientID", client.ID)
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

func (m *Manager) ClientsUnsafe() map[string]*Client {
	return m.clients
}

func (m *Manager) RLock()   { m.mu.RLock() }
func (m *Manager) RUnlock() { m.mu.RUnlock() }

// HandleNewConnection creates a new client from an accepted connection.
// Can be used in an HTTP handler.
func (m *Manager) HandleNewConnection(userID string, conn *websocket.Conn) {
	client := NewClient(userID, conn, m)
	client.Start()
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
	m.logger.Info("shutting down Websocket manager")
}
