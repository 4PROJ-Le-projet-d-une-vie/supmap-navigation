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
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			return handler.NewErrWithStatus(http.StatusBadRequest, errors.New("missing user_id"))
		}

		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return handler.NewErrWithStatus(http.StatusInternalServerError, fmt.Errorf("websocket accept: %w", err))
		}

		s.WebsocketManager.HandleNewConnection(userID, conn)
		return nil
	})
}
