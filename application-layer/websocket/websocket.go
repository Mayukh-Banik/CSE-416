package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var wsClients = make(map[*websocket.Conn]bool) // Connected WebSocket clients
var wsBroadcast = make(chan string)            // Channel for broadcast messages

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins; customize as needed
	},
}

// WsHandler handles incoming WebSocket connections
func WsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer ws.Close()

	wsClients[ws] = true
	log.Println("New WebSocket connection established")

	// Keep the connection alive
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Printf("WebSocket connection closed: %v", err)
			delete(wsClients, ws)
			break
		}
	}
}

// BroadcastMessages continuously sends messages from wsBroadcast to all connected WebSocket clients
func BroadcastMessages() {
	for {
		msg := <-wsBroadcast
		for client := range wsClients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Printf("WebSocket broadcast error: %v", err)
				client.Close()
				delete(wsClients, client)
			}
		}
	}
}

// SendMessage sends a message to the WebSocket broadcast channel
func SendMessage(message string) {
	wsBroadcast <- message
}
