package download

import (
	dht_kad "application-layer/dht"
	"application-layer/models"
	wsConn "application-layer/websocket"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
)

var PendingRequests = make(map[string]models.Transaction)
var mutex = &sync.Mutex{}

func handleDownloadRequest(w http.ResponseWriter, r *http.Request) {
	// handle when user tries to download file they already have (same hash)
	var request models.Transaction

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}
	fmt.Println("Handling download request: ", request.FileHash)

	// Set the requester's ID (assumed to be from the node's local ID)
	request.RequesterID = dht_kad.DHT.Host().ID().String()

	if request.RequesterID == request.TargetID {
		http.Error(w, "Cannot request self as a provider", http.StatusBadRequest)
		return
	}
	fmt.Printf("requester:%v | provider:%v \n ", request.RequesterID, request.TargetID)

	request.Status = "pending"
	request.CreatedAt = time.Now().Format(time.RFC3339)

	// Connect to the target peer and send the request over P2P stream
	if err := dht_kad.ConnectToPeerUsingRelay(dht_kad.DHT.Host(), request.TargetID); err != nil {
		http.Error(w, "Failed to connect to target peer", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if err := sendDownloadRequest(request); err != nil {
		http.Error(w, "Failed to send download request", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Store request in pending requests
	PendingRequests[request.FileHash] = request

	// Send acknowledgment back to the requester
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "request sent"})
}

func sendDownloadRequest(requestMetadata models.Transaction) error {
	// Create a new P2P stream to send the download request
	requestStream, err := dht_kad.CreateNewStream(dht_kad.DHT.Host(), requestMetadata.TargetID, "/sendRequest/p2p")
	if err != nil {
		return fmt.Errorf("error sending download request: %v", err)
	}
	defer requestStream.Close()

	// Marshal the request metadata to JSON
	requestData, err := json.Marshal(requestMetadata)
	if err != nil {
		return fmt.Errorf("error marshaling download request data: %v", err)
	}

	// Send the JSON data over the stream
	_, err = requestStream.Write(requestData)
	if err != nil {
		return fmt.Errorf("error sending download request data: %v", err)
	}

	fmt.Printf("Sent download request for file hash %s to target peer %s\n", requestMetadata.FileHash, requestMetadata.TargetID)
	return nil
}

func NotifyFrontendOfPendingRequest(request models.Transaction) {
	// Prepare the acknowledgment message
	acknowledgment := map[string]string{
		"status":    request.Status,
		"fileHash":  request.FileHash,
		"requester": request.RequesterID,
	}
	acknowledgmentData, _ := json.Marshal(acknowledgment)

	// Retrieve the WebSocket connection for the specific user
	if wsConn, exists := wsConn.WsConnections[request.TargetID]; exists {
		// Send the notification over the WebSocket connection
		if err := wsConn.WriteJSON(acknowledgmentData); err != nil {
			fmt.Println("Error sending notification to frontend:", err)
		}
	} else {
		fmt.Println("WebSocket connection not found for node:", request.TargetID)
	}
}

func handleDownloadRequestOrResponse(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// Check if the request exists in pendingRequests
	existingTransaction, exists := PendingRequests[transaction.FileHash]
	if !exists {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	// Handle based on the transaction status
	switch transaction.Status {
	case "accepted":
		existingTransaction.Status = "accepted"
		// Send file to requester
		sendFile(existingTransaction.TargetID, existingTransaction.FileHash)
	case "declined":
		existingTransaction.Status = "declined"
		// Notify decline
		notifyDecline(existingTransaction.TargetID, existingTransaction.FileHash)
	}

	// Update the transaction status in pendingRequests
	PendingRequests[transaction.FileHash] = existingTransaction
}

func notifyDecline(targetID string, fileHash string) {
	declineMessage := map[string]string{
		"status":   "declined",
		"fileHash": fileHash,
	}
	declineData, _ := json.Marshal(declineMessage)

	// Send the decline response back to the target peer
	requestStream, err := dht_kad.CreateNewStream(dht_kad.DHT.Host(), targetID, "/requestResponse/p2p")
	if err == nil {
		requestStream.Write(declineData)
		requestStream.Close()
	}
}

func sendFile(targetID string, fileHash string) {
	// In P2P systems, file transfer will be done through streams as well.
	// You need to implement logic for transferring the file to the requester.
	// This is just a placeholder to show where file transfer would occur.
	fmt.Printf("Sending file %s to requester %s...\n", fileHash, targetID)

	// Create a new stream to send the file
	fileStream, err := dht_kad.CreateNewStream(dht_kad.DHT.Host(), targetID, "/sendFile/p2p")
	if err != nil {
		fmt.Println("Error creating file stream:", err)
		return
	}
	defer fileStream.Close()

	// Logic to retrieve and send the file in chunks goes here
	// For now, we'll just simulate sending a file
	fileData := []byte("This is the file data for " + fileHash)
	_, err = fileStream.Write(fileData)
	if err != nil {
		fmt.Println("Error sending file:", err)
	}
}

func HandleDownloadRequestStream(stream network.Stream) {
	defer stream.Close()

	// Decode the request metadata
	var requestMetadata models.Transaction
	if err := json.NewDecoder(stream).Decode(&requestMetadata); err != nil {
		fmt.Println("Error decoding download request:", err)
		return
	}

	// Store the request in pending requests for the target node to approve
	PendingRequests[requestMetadata.FileHash] = requestMetadata

	// Notify the target frontend of a pending download request
	NotifyFrontendOfPendingRequest(requestMetadata)
}

func handleGetPendingRequests(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")

	// Convert the map to a slice of transactions
	var transactions []models.Transaction
	for _, transaction := range PendingRequests {
		transactions = append(transactions, transaction)
	}

	err := json.NewEncoder(w).Encode(transactions)
	if err != nil {
		http.Error(w, "Failed to encode pending requests", http.StatusInternalServerError)
		return
	}
}
