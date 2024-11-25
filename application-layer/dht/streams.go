package dht_kad

import (
	"application-layer/models"
	"application-layer/utils"
	"application-layer/websocket"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
)

var (
	PendingRequests = make(map[string]models.Transaction) // all requests made by host node
	FileHashToPath  = make(map[string]string)             // file paths of files uploaded by host node
	Mutex           = &sync.Mutex{}
	FileMapMutex    = &sync.Mutex{}
	dir             = filepath.Join("..", "squidcoinFiles")
	RefreshResponse []models.FileMetadata

	dirPath            = filepath.Join("..", "utils")
	UploadedFilePath   = filepath.Join(dirPath, "files.json")
	DownloadedFilePath = filepath.Join(dirPath, "downloadedFiles.json")
)

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
	declineData, err := json.Marshal(declineMessage)
	if err != nil {
		log.Printf("Error marshaling decline message: %v", err)
		return
	}

	// Send decline to the target peer
	requestStream, err := CreateNewStream(DHT.Host(), targetID, "/requestResponse/p2p")
	if err != nil {
		log.Printf("Error creating stream to target peer %s: %v", targetID, err)
		return
	}
	defer requestStream.Close()

	// Write data with a newline as a delimiter
	declineData = append(declineData, '\n')
	_, err = requestStream.Write(declineData)
	if err != nil {
		log.Printf("Error writing to stream: %v", err)
		return
	}

	log.Printf("Decline message sent to peer %s for file hash %s", targetID, fileHash)
}

// // send metadata before sending file content
// func sendMetadata(stream network.Stream, fileHash string) error {
// 	// Retrieve the file data from the DHT using the file hash
// 	data, err := DHT.GetValue(GlobalCtx, "/orcanet/"+fileHash)
// 	if err != nil {
// 		return fmt.Errorf("sendMetadata: file hash not found in DHT: %w", err)
// 	}

// 	var metadata models.DHTMetadata

// 	// unmarshal JSON data into the struct
// 	err = json.Unmarshal(data, &metadata)
// 	if err != nil {
// 		return fmt.Errorf("sendMetadata: error decoding file metadata: %w", err)
// 	}

// 	metadataJSON, err := json.Marshal(metadata)
// 	if err != nil {
// 		return fmt.Errorf("sendMetadata: error encoding metadata to data: %w", err)
// 	}

// 	// newline signifies end of metadata
// 	metadataJSON = append(metadataJSON, '\n')

// 	// Write metadata to stream
// 	_, err = stream.Write(metadataJSON)
// 	if err != nil {
// 		log.Fatalf("sendMetadata: failed to write metadata to stream: %s", err)
// 	}

// 	fmt.Println("metadata sent successfully :)")
// 	return nil
// }

// send metadata before sending file content
func sendMetadata(stream network.Stream, fileHash string) error {
	// Retrieve the file data from the DHT using the file hash
	data, err := DHT.GetValue(GlobalCtx, "/orcanet/"+fileHash)
	if err != nil {
		return fmt.Errorf("sendMetadata: file hash not found in DHT: %w", err)
	}

	var metadata models.DHTMetadata

	// unmarshal JSON data into the struct
	err = json.Unmarshal(data, &metadata)
	if err != nil {
		return fmt.Errorf("sendMetadata: error decoding file metadata: %w", err)
	}

	var fileMetadata = models.FileMetadata{
		Name:              metadata.Name,
		NameWithExtension: metadata.NameWithExtension,
		Type:              metadata.Type,
		Size:              metadata.Size,
		Description:       metadata.Description,
		Hash:              fileHash,
		OriginalUploader:  false,
		IsPublished:       true, //automatically become provider when you download file
	}

	fileMetadataJSON, err := json.Marshal(fileMetadata)
	if err != nil {
		return fmt.Errorf("sendMetadata: error encoding metadata to data: %w", err)
	}

	// newline signifies end of metadata
	fileMetadataJSON = append(fileMetadataJSON, '\n')

	// Write metadata to stream
	_, err = stream.Write(fileMetadataJSON)
	if err != nil {
		log.Fatalf("sendMetadata: failed to write metadata to stream: %s", err)
	}

	fmt.Println("sendMetadata: metadata sent successfully: ", fileMetadataJSON)
	return nil
}

func sendFile(host host.Host, targetID string, fileHash string, requesterID string, fileName string) {
	fmt.Printf("Sending file %s to requester %s...\n", fileHash, targetID)

	filePath := FileHashToPath[fileHash]
	fmt.Println("sendFile: filePath: ", filePath)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("error: file %s not found in %s\n", fileHash, filePath)
		return
	}

	fmt.Printf("sending file %s to requester %s \n", fileName, requesterID)

	// create stream to send the file
	fileStream, err := CreateNewStream(host, targetID, "/sendFile/p2p")
	if err != nil {
		fmt.Println("Error creating file stream:", err)
		return
	}
	defer fileStream.Close()

	// send metadata first
	err = sendMetadata(fileStream, fileHash)
	if err != nil {
		fmt.Println("error sending file metadata")
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("error opening file %s: %v", filePath, err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := make([]byte, 4096)

	for {
		n, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("all of file sent")
				break
			}
			fmt.Printf("error reading file %s: %v\n", filePath, err)
			return
		}

		_, err = fileStream.Write(buffer[:n])
		if err != nil {
			fmt.Printf("error sending chunk to requester %s: %v\n", requesterID, err)
			return
		}
		fmt.Printf("sent %d bytes to requester %s\n", n, requesterID)
	}
}

func SendRefreshFilesRequest(nodeID string, wg *sync.WaitGroup) error {
	defer wg.Done()

	fmt.Println("Requesting all files data from ", nodeID)
	// refreshResponse = []model.FileMetadata{}
	requestStream, err := CreateNewStream(DHT.Host(), nodeID, "/sendRefreshRequest/p2p")
	if err != nil {
		return fmt.Errorf("error sending refresh request")
	}
	defer requestStream.Close()

	request := models.RefreshRequest{
		Message:     "gimme all your files",
		RequesterID: PeerID,
		TargetID:    nodeID,
	}

	// Marshal the request metadata to JSON
	requestData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error marshaling refresh request data: %v", err)
	}

	// send JSON data over the stream
	_, err = requestStream.Write(requestData)
	if err != nil {
		return fmt.Errorf("error sending refresh request data: %v", err)
	}

	fmt.Printf("Sent refresh request to target peer %s\n", nodeID)
	return nil
}

func sendRefreshResponse(node host.Host, targetID string) error {
	fmt.Printf("sendRefreshResponse: sending all files to %v", targetID)
	dirPath := filepath.Join("..", "utils")
	filePath := filepath.Join(dirPath, "files.json")

	fileData, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to read files.json: %v", err)
	}

	// var files []models.FileMetadata
	// if err := json.Unmarshal(fileData, &files); err != nil {
	// 	return fmt.Errorf("failed to parse JSON: %v", err)

	// }

	refreshRequestStream, err := CreateNewStream(node, targetID, "/sendRefreshResponse/p2p")
	if err != nil {
		fmt.Println("Error creating file stream:", err)
	}
	defer refreshRequestStream.Close()

	reader := bufio.NewReader(fileData)
	buffer := make([]byte, 4096)

	for {
		n, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("All JSON data sent")
				break
			}
			return fmt.Errorf("error reading JSON data: %w", err)
		}

		_, err = refreshRequestStream.Write(buffer[:n])
		if err != nil {
			return fmt.Errorf("error sending byte data for JSON")
		}
		fmt.Printf("Sent %d bytes\n", n)
	}
	return nil
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
		log.Println("receieveDownloadRequest: files in FileHashToPath: ", FileHashToPath)
		fmt.Print("FILEHASHTOPATH", FileHashToPath)
		fmt.Print("Request.FileHash", request.FileHash)
		// send file to requester if it exists
		if FileHashToPath[request.FileHash] != "" {
			fmt.Println("receivedownloadrequest: sending file")
			sendFile(node, request.RequesterID, request.FileHash, PeerID, request.FileName)
		} else {
			fmt.Println("receivedownloadrequest: decline")
			sendDecline(request.RequesterID, request.FileHash)
		}

	})
}

func receiveDecline(node host.Host) {
	node.SetStreamHandler("/requestResponse/p2p", func(s network.Stream) {
		defer s.Close()
		buf := bufio.NewReader(s)

		// Read data until a newline character
		data, err := buf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				log.Printf("Stream closed by peer: %s", s.Conn().RemotePeer())
			} else {
				log.Printf("Error reading from stream: %v", err)
			}
			return
		}

		log.Printf("Raw data received: %s", data)

		// Unmarshal the JSON data
		var declineMessage map[string]string
		err = json.Unmarshal(data, &declineMessage)
		if err != nil {
			log.Printf("Error unmarshalling data: %v", err)
			return
		}

		// Process the decline message
		status, statusOK := declineMessage["status"]
		fileHash, fileHashOK := declineMessage["fileHash"]

		if statusOK && fileHashOK && status == "declined" {
			log.Printf("Received decline message for file with hash: %s", fileHash)
			websocket.NotifyFrontend(declineMessage)
			// Notify user on the frontend of the decline
			// Update transaction details to DECLINED
		} else {
			log.Println("Received invalid decline message")
		}
	})
}

func receiveFile(node host.Host) {
	fmt.Println("listening for file data")
	// listen for streams on "/sendFile/p2p"
	node.SetStreamHandler("/sendFile/p2p", func(s network.Stream) {
		defer s.Close()

		// read metadata - use metadata struct later
		buf := bufio.NewReader(s)
		// Read metadata
		metadataJSON, err := buf.ReadBytes('\n') // Read until newline
		if err != nil {
			log.Fatalf("Failed to read metadata: %v", err)
		}

		// Parse JSON metadata
		var metadata models.FileMetadata

		err = json.Unmarshal(metadataJSON, &metadata)
		if err != nil {
			log.Fatalf("Failed to unmarshal metadata: %v", err)
		}

		fmt.Printf("Received metadata: FileName=%s\n", metadata.Name)

		// open file for writing
		outputPath := filepath.Join(dir, metadata.NameWithExtension)
		file, err := os.Create(outputPath)
		if err != nil {
			log.Printf("error creating file %s: %v\n", outputPath, err)
		}
		defer file.Close()

		// read and write chunks of data
		buffer := make([]byte, 4086)
		for {
			n, err := buf.Read(buffer)
			if err != nil {
				if err == io.EOF {
					log.Printf("file %s received and saved to %s\n", metadata.Name, outputPath)
					break
				}
				log.Printf("Ererrorror reading file chunk: %v\n", err)
				return
			}

			_, writeErr := file.Write(buffer[:n])
			if writeErr != nil {
				log.Printf("error writing to file %s: %v\n", outputPath, writeErr)
				return
			}

			log.Printf("receieved and wrote %d bytes of file %s\n", n, metadata.Name)
		}

		// after successfully downloading file, the user is now a provider of the file
		utils.SaveOrUpdateFile(metadata, dirPath, DownloadedFilePath)

		FileMapMutex.Lock()

		fmt.Println("file name with extension: ", metadata.NameWithExtension)
		FileHashToPath[metadata.Hash] = filepath.Join(dir, metadata.NameWithExtension) // add file and its path to the map
		FileMapMutex.Unlock()

		ProvideKey(GlobalCtx, DHT, metadata.Hash) // must be published - update dht with new provider
	})
}

func receiveRefreshRequest(node host.Host) error {
	fmt.Println("listening for refresh requests")
	node.SetStreamHandler("/sendRefreshRequest/p2p", func(s network.Stream) {
		defer s.Close()
		buf := bufio.NewReader(s)

		data, err := io.ReadAll(buf) //read in request

		if err != nil {
			if err == io.EOF {
				log.Printf("Stream closed by peer :%s", s.Conn().RemotePeer())
			} else {
				log.Printf("Error receiving refresh request %v", err)
			}
			return
		}

		var refreshReq models.RefreshRequest
		err = json.Unmarshal(data, &refreshReq)

		sendRefreshResponse(node, refreshReq.RequesterID)
	})
	return nil
}

func receiveRefreshResponse(node host.Host) {
	fmt.Println("Listening for refresh response")
	node.SetStreamHandler("/sendRefreshResponse/p2p", func(s network.Stream) {
		defer s.Close()

		// Use a buffer to read the incoming data
		var receivedData bytes.Buffer
		buf := make([]byte, 4096) // Chunk size should match sender's buffer

		for {
			n, err := s.Read(buf)
			if err != nil {
				if err == io.EOF {
					fmt.Println("All data received from sender")
					break
				}
				log.Fatalf("Failed to read data: %v", err)
			}
			receivedData.Write(buf[:n])
		}

		// Parse the accumulated JSON data
		var fileData []models.FileMetadata
		err := json.Unmarshal(receivedData.Bytes(), &fileData)
		fmt.Println("file data received from refresh", fileData)
		if err != nil {
			log.Fatalf("Error unmarshaling received data: %v", err)
		}

		// for _,file := range fileData {
		// 	RefreshResponse = append(RefreshResponse, file)

		// }
		Mutex.Lock()
		RefreshResponse = append(RefreshResponse, fileData...)
		Mutex.Unlock()
		fmt.Println("Received files for marketplace:", RefreshResponse)
	})
}

// listen on streams
func setupStreams(node host.Host) {
	receieveDownloadRequest(node)
	receiveDecline(node)
	receiveFile(node)
	receiveRefreshRequest(node)
	receiveRefreshResponse(node)
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
		sendFile(DHT.Host(), existingTransaction.TargetID, existingTransaction.FileHash, existingTransaction.RequesterID, existingTransaction.FileName)
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

// func NotifyFrontendOfPendingRequest(request models.Transaction) {
// 	// Prepare acknowledgment message
// 	acknowledgment := map[string]string{
// 		"status":    request.Status,
// 		"fileHash":  request.FileHash,
// 		"requester": request.RequesterID,
// 	}
// 	acknowledgmentData, _ := json.Marshal(acknowledgment)

// 	// Retrieve the WebSocket connection for the specific user
// 	if wsConn, exists := websocket.WsConnections[request.TargetID]; exists {
// 		// Send the notification over the WebSocket connection
// 		if err := wsConn.WriteJSON(acknowledgmentData); err != nil {
// 			fmt.Println("Error sending notification to frontend:", err)
// 		}
// 	} else {
// 		fmt.Println("WebSocket connection not found for node:", request.TargetID)
// 	}
// }
