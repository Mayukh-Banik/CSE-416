package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Adjust for production security
	},
}

// Thread-safe client management
var (
	clients = make(map[*websocket.Conn]bool) // Tracks active connections
	mu      sync.Mutex                      // Protects the clients map
)

// Handles new WebSocket connections
func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket Upgrade error:", err)
		return
	}
	log.Println("New WebSocket connection established")

	// Add the client to the clients map
	mu.Lock()
	clients[conn] = true
	mu.Unlock()

	// Continuously read messages (if needed)
	go func(conn *websocket.Conn) {
		defer func() {
			// Remove the client on disconnect
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			conn.Close()
			log.Println("WebSocket connection closed")
		}()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Println("WebSocket read error:", err)
				break
			}
		}
	}(conn)
}

// Sends messages to all connected clients
func NotifyFrontend(declineMessage map[string]string) {
	mu.Lock()
	defer mu.Unlock()

	for client := range clients {
		err := client.WriteJSON(declineMessage)
		if err != nil {
			log.Println("Error sending WebSocket message:", err)
			client.Close()
			delete(clients, client)
		}
	}
}
