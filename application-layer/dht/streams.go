package dht_kad

import (
	"application-layer/models"
	"application-layer/websocket"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
)

var PendingRequests = make(map[string]models.Transaction) // all requests made by host node
var FilePath = make(map[string]string)                    // file paths of files uploaded by host node
var Mutex = &sync.Mutex{}

// SENDING FUNCTIONS

func SendDownloadRequest(requestMetadata models.Transaction) error {
	// create stream to send the download request
	fmt.Println("Sending download request via stream /sendRequest/p2p")
	requestStream, err := CreateNewStream(DHT.Host(), requestMetadata.TargetID, "/sendRequest/p2p")
	if err != nil {
		return fmt.Errorf("error sending download request: %v", err)
	}
	defer requestStream.Close()

	// Marshal the request metadata to JSON
	requestData, err := json.Marshal(requestMetadata)
	if err != nil {
		return fmt.Errorf("error marshaling download request data: %v", err)
	}

	// send JSON data over the stream
	_, err = requestStream.Write(requestData)
	if err != nil {
		return fmt.Errorf("error sending download request data: %v", err)
	}

	fmt.Printf("Sent download request for file hash %s to target peer %s\n", requestMetadata.FileHash, requestMetadata.TargetID)
	return nil
}

func sendDecline(targetID string, fileHash string) {
	declineMessage := map[string]string{
		"status":   "declined",
		"fileHash": fileHash,
	}
	declineData, _ := json.Marshal(declineMessage)

	// send decline to the target peer
	requestStream, err := CreateNewStream(DHT.Host(), targetID, "/requestResponse/p2p")
	if err == nil {
		requestStream.Write(declineData)
		requestStream.Close()
	}
}

func sendFile(host host.Host, targetID string, fileHash string, requesterID string) {
	fmt.Printf("Sending file %s to requester %s...\n", fileHash, targetID)

	// create stream to send the file
	fileStream, err := CreateNewStream(DHT.Host(), targetID, "/sendFile/p2p")
	if err != nil {
		fmt.Println("Error creating file stream:", err)
		return
	}
	defer fileStream.Close()

	// path := FilePath[fileHash]
	// fileType := filepath.Ext(path)

	// fileContent, err := os.Open(path)
	// if err != nil {
	// 	fmt.Println("error sending file to requester: %v", err)
	// }
	// defer fileContent.Close()

	// // Simulate file sending (this should be replaced with actual file retrieval and transfer logic)
	// fileData := []byte("This is the file data for " + fileHash)
	// _, err = fileStream.Write(fileData)
	// if err != nil {
	// 	fmt.Println("Error sending file:", err)
	//

	// for testing below - must implement actual file sharing
	_, err = fileStream.Write([]byte("sending file to peer\n"))
	if err != nil {
		log.Fatalf("Failed to write to stream: %s", err)
	}
}

// RECEIVING FUNCTIONS

func receieveDownloadRequest(node host.Host) {
	fmt.Println("listening for download requests")
	// listen for streams on "/sendRequest/p2p"
	node.SetStreamHandler("/sendRequest/p2p", func(s network.Stream) {
		defer s.Close()
		// Create a buffered reader to read data from the stream
		buf := bufio.NewReader(s)
		// Read data from the stream
		data, err := io.ReadAll(buf) // everything - should just be a transaction struct

		if err != nil {
			if err == io.EOF {
				log.Printf("Stream closed by peer: %s", s.Conn().RemotePeer())
			} else {
				log.Printf("Error reading from stream: %v", err)
			}
			return
		}

		var request models.Transaction
		err = json.Unmarshal(data, &request)
		if err != nil {
			fmt.Printf("error unmarshalling file request: %v", err)
			return
		}
		log.Printf("Received data: %s", data)

		// send file to requester if it exists
		if FilePath[request.FileHash] != "" {
			sendFile(DHT.Host(), request.RequesterID, request.FileHash, PeerID)
		} else {
			sendDecline(request.RequesterID, request.FileHash)
		}

	})
}

func receieveFile(node host.Host) {
	fmt.Println("listening for file data")
	// listen for streams on "/sendFile/p2p"
	node.SetStreamHandler("/sendFile/p2p", func(s network.Stream) {
		defer s.Close()
		buf := bufio.NewReader(s)
		data, err := buf.ReadBytes('\n') // Reads until a newline character

		// chunck, err := io.ReadAll(buf) // receieve chunk by chunck and put back together

		if err != nil {
			if err == io.EOF {
				log.Printf("Stream closed by peer: %s", s.Conn().RemotePeer())
			} else {
				log.Printf("Error reading from stream: %v", err)
			}
			return
		}
		log.Printf("Received data: %s", data)
	})
}

func receiveDecline(node host.Host) {
	node.SetStreamHandler("/requestResponse/p2p", func(s network.Stream) {
		defer s.Close()
		buf := bufio.NewReader(s)
		data, err := buf.ReadBytes('\n') // Reads until a newline character

		if err != nil {
			if err == io.EOF {
				log.Printf("Stream closed by peer: %s", s.Conn().RemotePeer())
			} else {
				log.Printf("Error reading from stream: %v", err)
			}
			return
		}
		log.Printf("Received data: %s", data)
		// Unmarshal the JSON data
		var declineMessage map[string]string
		err = json.Unmarshal(data, &declineMessage)
		if err != nil {
			log.Println("Error unmarshalling data:", err)
			return
		}

		// Process the decline message
		if status, ok := declineMessage["status"]; ok && status == "declined" {
			fileHash := declineMessage["fileHash"]
			log.Printf("Received decline message for file with hash: %s", fileHash)
			// notify user on the front end of decline
			// update transaction detail to DECLINED
		} else {
			log.Println("Received invalid decline message")
		}
		log.Printf("Received data: %s", data)
	})
}

// OTHER - IGNORE WILL PROB DELETE
func handleDownloadRequestOrResponse(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// Check if the request exists in pendingRequests
	Mutex.Lock()
	existingTransaction, exists := PendingRequests[transaction.FileHash]
	Mutex.Unlock()

	if !exists {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	// Handle based on the transaction status
	switch transaction.Status {
	case "accepted":
		existingTransaction.Status = "accepted"
		// Send file to requester
		sendFile(DHT.Host(), existingTransaction.TargetID, existingTransaction.FileHash, existingTransaction.RequesterID)
	case "declined":
		existingTransaction.Status = "declined"
		// Notify decline
		sendDecline(existingTransaction.TargetID, existingTransaction.FileHash)
	}

	// Update the transaction status in pendingRequests
	Mutex.Lock()
	PendingRequests[transaction.FileHash] = existingTransaction
	Mutex.Unlock()
}

func NotifyFrontendOfPendingRequest(request models.Transaction) {
	// Prepare acknowledgment message
	acknowledgment := map[string]string{
		"status":    request.Status,
		"fileHash":  request.FileHash,
		"requester": request.RequesterID,
	}
	acknowledgmentData, _ := json.Marshal(acknowledgment)

	// Retrieve the WebSocket connection for the specific user
	if wsConn, exists := websocket.WsConnections[request.TargetID]; exists {
		// Send the notification over the WebSocket connection
		if err := wsConn.WriteJSON(acknowledgmentData); err != nil {
			fmt.Println("Error sending notification to frontend:", err)
		}
	} else {
		fmt.Println("WebSocket connection not found for node:", request.TargetID)
	}
}

// set up stream handlers/listeners?

func setupStreams(node host.Host) {
	receieveDownloadRequest(node)
	receiveDecline(node)
	receieveFile(node)
}
