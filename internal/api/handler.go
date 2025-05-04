package api

import (
	"github.com/coder/websocket"
	"github.com/matheodrd/httphelper/handler"
	"net/http"
	"supmap-navigation/internal/ws"
)

func (s *Server) wsHandler() http.HandlerFunc {
	return handler.Handler(func(w http.ResponseWriter, r *http.Request) error {
		id := r.URL.Query().Get("user_id")
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return handler.NewErrWithStatus(http.StatusInternalServerError, err)
		}
		client := ws.NewClient(id, conn, s.WebsocketManager)
		client.Start()
		return nil
	})
}
