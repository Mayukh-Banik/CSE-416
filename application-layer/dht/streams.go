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
	"os"
	"path/filepath"
	"sync"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
)

var (
	PendingRequests        = make(map[string]models.Transaction) // all requests made by host node
	FileHashToPath         = make(map[string]string)             // file paths of files uploaded by host node
	Mutex                  = &sync.Mutex{}
	FileMapMutex           = &sync.Mutex{}
	dir                    = filepath.Join("..", "squidcoinFiles")
	MarketplaceFiles       []models.FileMetadata
	MarketplaceFilesSignal = make(chan struct{})
	dirPath                = filepath.Join("..", "utils")
	UploadedFilePath       = filepath.Join(dirPath, "files.json")
	DownloadedFilePath     = filepath.Join(dirPath, "downloadedFiles.json")
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

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("sendMetadata: error encoding metadata to data: %w", err)
	}

	// newline signifies end of metadata
	metadataJSON = append(metadataJSON, '\n')

	// Write metadata to stream
	_, err = stream.Write(metadataJSON)
	if err != nil {
		log.Fatalf("sendMetadata: failed to write metadata to stream: %s", err)
	}

	fmt.Println("metadata sent successfully :)")
	return nil
}

func sendFile(host host.Host, targetID string, fileHash string, requesterID string, fileName string) {
	fmt.Printf("Sending file %s to requester %s...\n", fileHash, targetID)

	filePath := FileHashToPath[fileHash]

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

func SendMarketFilesRequest(nodeID string) error {
	fmt.Println("Requesting marketplace files data from cloud node ", nodeID)

	// refreshResponse = []model.FileMetadata{}
	requestStream, err := CreateNewStream(DHT.Host(), nodeID, "/sendRefreshRequest/p2p")
	if err != nil {
		return fmt.Errorf("error sending refresh request")
	}
	defer requestStream.Close()

	requestData := []byte(PeerID)

	// send JSON data over the stream
	_, err = requestStream.Write(requestData)
	if err != nil {
		return fmt.Errorf("error sending refresh request: %v", err)
	}

	fmt.Printf("Sent refresh request to cloud node %s\n", nodeID)
	return nil
}

func SendCloudNodeFiles(nodeID string, fileMetadata models.FileMetadata) error {
	stream, err := CreateNewStream(DHT.Host(), nodeID, "/sendRefreshRequest/p2p")
	if err != nil {
		return fmt.Errorf("error sending file to cloud node")
	}
	defer stream.Close()

	fileData, err := json.Marshal(fileMetadata)
	if err != nil {
		return fmt.Errorf("sendCloudNodeFiles: failed to marshal file metadata: %v", err)
	}

	_, err = stream.Write(fileData)
	if err != nil {
		return fmt.Errorf("sendCloudNodeFiles: failed to send file metadata to cloud node")
	}

	fmt.Printf("Sent file metadata to cloud node %s\n", nodeID)
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
		log.Println("files in FileHashToPath: ", FileHashToPath)

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
		outputPath := filepath.Join(dir, metadata.Name)
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

		FileHashToPath[metadata.Hash] = filepath.Join(dir, metadata.Name) // add file and its path to the map
		FileMapMutex.Unlock()

		ProvideKey(GlobalCtx, DHT, metadata.Hash) // must be published - update dht with new provider
	})
}

func receiveMarketplaceFiles(node host.Host) {
	fmt.Println("Listening for refresh response")
	node.SetStreamHandler("/marketplaceFiles/p2p", func(s network.Stream) {
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

		Mutex.Lock()
		MarketplaceFiles = append(MarketplaceFiles[:0], fileData...)
		Mutex.Unlock()

		MarketplaceFilesSignal <- struct{}{}
		fmt.Println("Received files for marketplace:", MarketplaceFiles)
	})
}

// listen on streams
func setupStreams(node host.Host) {
	receieveDownloadRequest(node)
	receiveDecline(node)
	receiveFile(node)
	receiveMarketplaceFiles(node)
}
