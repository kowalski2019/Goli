package handler

import (
	ws "goli/websocket"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}

// ServeWebSocket handles websocket requests from clients
func ServeWebSocket(hub *ws.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &ws.Client{
		Hub:  hub,
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	client.Hub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in new goroutines
	go client.WritePump()
	go client.ReadPump()
}
