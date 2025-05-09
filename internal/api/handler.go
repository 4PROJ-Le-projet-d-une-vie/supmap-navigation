package api

import (
	"errors"
	"fmt"
	"github.com/coder/websocket"
	"github.com/matheodrd/httphelper/handler"
	"net/http"
)

func (s *Server) wsHandler() http.HandlerFunc {
	return handler.Handler(func(w http.ResponseWriter, r *http.Request) error {
		sessionID := r.URL.Query().Get("session_id")
		if sessionID == "" {
			return handler.NewErrWithStatus(http.StatusBadRequest, errors.New("missing session_id"))
		}

		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return handler.NewErrWithStatus(http.StatusInternalServerError, fmt.Errorf("websocket accept: %w", err))
		}

		s.WebsocketManager.HandleNewConnection(sessionID, conn)
		return nil
	})
}
