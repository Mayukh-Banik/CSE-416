package download

import (
	dht_kad "application-layer/dht"
	"application-layer/models"
	"application-layer/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// api calls
func handleDownloadRequest(w http.ResponseWriter, r *http.Request) {
	var request models.Transaction

	// Decode the incoming request data into the transaction struct
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}
	fmt.Println("handling download request for", request.FileHash)
	request.RequesterID = dht_kad.DHT.Host().ID().String()
	request.TransactionID = uuid.New().String()
	// Log the requester and provider IDs
	fmt.Printf("Requesting file download: Requester: %v | Provider: %v\n", request.RequesterID, request.TargetID)

	// Update the request status to "pending"
	request.Status = "pending"
	request.CreatedAt = time.Now().Format(time.RFC3339)

	// Connect to the target peer and send the download request via P2P
	if err := dht_kad.ConnectToPeerUsingRelay(dht_kad.DHT.Host(), request.TargetID); err != nil {
		http.Error(w, "Failed to connect to target peer", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	// fmt.Println("Connected peers:", dht_kad.Host.Peerstore().Peers())

	// just testing if nodes are connected
	// dht_kad.SendDataToPeer(dht_kad.DHT.Host(), request.TargetID)

	// actually send the download request
	if err := dht_kad.SendDownloadRequest(request); err != nil {
		http.Error(w, "Failed to send download request", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Store request in pending requests map (thread-safe)
	dht_kad.Mutex.Lock()
	dht_kad.PendingRequests[request.FileHash] = request
	dht_kad.Mutex.Unlock()

	utils.AddOrUpdateTransaction(request)

	// Send acknowledgment back to the requester
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "request sent"})
}
