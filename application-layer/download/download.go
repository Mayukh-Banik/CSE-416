package download

import (
	dht_kad "application-layer/dht"
	"application-layer/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// api calls
func handleDownloadRequest(w http.ResponseWriter, r *http.Request) {
	var request models.Transaction

	// Decode the incoming request data into the transaction struct
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	request.RequesterID = dht_kad.DHT.Host().ID().String()

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

	// Send acknowledgment back to the requester
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "request sent"})
}

// prob gonna delete
func handleGetPendingRequests(w http.ResponseWriter, r *http.Request) {
	dht_kad.Mutex.Lock()
	defer dht_kad.Mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")

	// Convert the map to a slice of transactions
	var transactions []models.Transaction
	for _, transaction := range dht_kad.PendingRequests {
		transactions = append(transactions, transaction)
	}

	err := json.NewEncoder(w).Encode(transactions)
	if err != nil {
		http.Error(w, "Failed to encode pending requests", http.StatusInternalServerError)
		return
	}
}

// additional downloading stuff

func updateDownloadJSON(metadata models.FileMetadata) {
	panic("unimplemented")
}
