package websocket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var WsConn *websocket.Conn                           // Declare this globally
var WsConnections = make(map[string]*websocket.Conn) // A map of userID to WebSocket connection

// WebSocket handler to establish a connection with the frontend
func handleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // You can add more checks for security here
		},
	}

	// Upgrade the HTTP connection to a WebSocket
	WsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}

	// Assume the frontend sends its user ID in the query string
	nodeID := r.URL.Query().Get("nodeID")
	WsConnections[nodeID] = WsConn

	// Optionally, listen for incoming messages (or just keep the connection open)
	go func() {
		for {
			_, _, err := WsConn.ReadMessage()
			if err != nil {
				// Handle disconnection or error
				delete(WsConnections, nodeID)
				WsConn.Close()
				break
			}
		}
	}()
}
