package main

import (
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

// Message represents a JSON message between client and server.
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// handleWS handle a single Websocket connection.
func handleWS(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"}, // A changer en prod
	})
	if err != nil {
		log.Printf("WebSocket accept error: %v", err)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "internal error")

	log.Println("New client connected")

	// Read/Write loop
	for {
		var msg Message
		err := wsjson.Read(ctx, conn, &msg)
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}

		switch msg.Type {
		case "location":
			log.Printf("Received location: %v", msg.Data)
			resp := Message{
				Type: "info",
				Data: "Ok, localisation reçue",
			}
			_ = wsjson.Write(ctx, conn, resp)
		default:
			resp := Message{
				Type: "error",
				Data: "Type de message inconnu",
			}
			_ = wsjson.Write(ctx, conn, resp)
		}
	}

	if err := conn.Close(websocket.StatusNormalClosure, ""); err != nil {
		log.Printf("Close error: %v", err)
	}
	log.Println("Client disconnected")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/navigation", handleWS)

	srv := &http.Server{
		Addr:         ":8081",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Println("supmap-navigation WS server listening on :8081/navigation")
	log.Fatal(srv.ListenAndServe())
}
