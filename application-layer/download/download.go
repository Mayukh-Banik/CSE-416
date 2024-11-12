package download

import (
	dht_kad "application-layer/dht"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
)

type DownloadRequest struct {
	TargetID    string `json:"TargetID`
	FileHash    string `json:"FileHash`
	RequesterID string `json:"RequesterID`
	Status      string `json:"Status`
	CreatedAt   string `json:"CreatedAt`
}

type DownloadResponse struct {
	FileHash string
	TargetID string
	Accepted bool
}

func handleDownloadRequest(w http.ResponseWriter, r *http.Request) {
	var request DownloadRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}
	fmt.Println("Handling download request: ", request.FileHash)

	// Set the requester's ID (assumed to be from the node's local ID)
	request.RequesterID = dht_kad.DHT.Host().ID().String()
	fmt.Printf("provider ID: %s\n", request.TargetID)
	fmt.Printf("requester ID: %s\n", request.RequesterID)
	// Decode the provider's ID and find the provider's address
	targetID, err := peer.Decode(request.TargetID)
	if err != nil {
		http.Error(w, "Invalid provider ID", http.StatusBadRequest)
		return
	}

	// is this needed??
	dht_kad.FindProviders(request.FileHash)
	provider, _ := dht_kad.FindSpecificProvider(request.FileHash, targetID)
	fmt.Printf("Found provider addresses: %s\n", provider.Addrs)

	dht_kad.ConnectToPeerUsingRelay(dht_kad.DHT.Host(), request.TargetID)
	sendDownloadRequest(request)

	fmt.Printf("Received download request: %+v\n", request)

	// Notify the provider (this function is not defined in the provided code, so assuming it's part of your logic)
	// notifyProvider(request)

	// Send a response back to the requester
	w.Header().Set("Content-Type", "application/json")
}

func sendDownloadRequest(requestMetadata DownloadRequest) error {
	requestMetadata.CreatedAt = time.Now().Local().String()
	requestStream, err := dht_kad.CreateNewStream(dht_kad.Host, requestMetadata.TargetID, "/sendRequest/p2p")
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

func notifyProvider(request DownloadRequest) {

}

func respondToDownloadRequest(w http.ResponseWriter, r *http.Request) {
	var response DownloadResponse
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if response.Accepted {
		startFileTransfer(response.TargetID, response.FileHash)
	} else {
		notifyDecline(response.TargetID, response.FileHash)
	}
}

func startFileTransfer(RequesterID string, FileHash string) {

}

func notifyDecline(RequesterID string, FileHash string) {

}
