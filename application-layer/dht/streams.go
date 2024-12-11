package dht_kad

import (
	"application-layer/models"
	"application-layer/utils"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

var (
	PendingRequests = make(map[string]models.Transaction) // all requests made by host node
	FileHashToPath  = make(map[string]string)             // file paths of files uploaded by host node
	Mutex           = &sync.Mutex{}
	FileMapMutex    = &sync.Mutex{}

	dir             = filepath.Join("..", "squidcoinFiles")
	RefreshResponse []models.FileMetadata
	ProxyResponse   []models.Proxy

	MarketplaceFiles []models.FileMetadata
	proxyHistory     []models.ProxyHistoryEntry

	MarketplaceFilesSignal = make(chan struct{})
	ProxiesSignal          = make(chan struct{}, 1)
	Proxies                []models.Proxy

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

func sendDecline(transaction models.Transaction) {
	transaction.Status = "declined"

	transactionJson, err := json.Marshal(transaction)
	if err != nil {
		log.Printf("sendDecline: error encoding metadata to data: %v", err)
		return
	}

	// Send decline to the target peer
	requestStream, err := CreateNewStream(DHT.Host(), transaction.RequesterID, "/requestResponse/p2p")
	if err != nil {
		log.Printf("Error creating stream to target peer %s: %v", transaction.RequesterID, err)
		return
	}
	defer requestStream.Close()

	// Write data with a newline as a delimiter
	transactionJson = append(transactionJson, '\n')
	_, err = requestStream.Write(transactionJson)
	if err != nil {
		log.Printf("Error writing to stream: %v", err)
		return
	}

	log.Printf("Decline message sent to peer %s for file hash %s", transaction.RequesterID, transaction.FileHash)
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

	fmt.Println("sending metadata for file: ", metadata.NameWithExtension)
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

func sendFile(host host.Host, request models.Transaction) {
	requesterID, fileHash := request.RequesterID, request.FileHash

	fmt.Printf("Sending file %s to requester %s...\n", fileHash, requesterID)

	filePath := FileHashToPath[request.FileHash]
	fmt.Println("sendFile: filePath: ", filePath)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("error: file %s not found in %s\n", fileHash, filePath)
		return
	}

	fmt.Printf("sending file %s to requester %s \n", request.FileName, requesterID)

	// create stream to send the file
	fileStream, err := CreateNewStream(host, request.RequesterID, "/sendFile/p2p")
	if err != nil {
		fmt.Println("Error creating file stream:", err)
		return
	}
	defer fileStream.Close()

	// sending transaction details before file metadata and content
	transactionData, err := json.Marshal(request)
	if err != nil {
		fmt.Printf("error marshaling transaction data: %v\n", err)
		return
	}

	transactionData = append(transactionData, '\n')

	_, err = fileStream.Write(transactionData)
	if err != nil {
		fmt.Printf("error sending transaction data: %v\n", err)
		return
	}
	fmt.Printf("Sent transaction data sent for %s\n", request.TransactionID)

	// send metadata next
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

func SendCloudNodeFiles(fileMetadata models.FileMetadata) error {
	stream, err := CreateNewStream(DHT.Host(), Cloud_node_id, "/nodeFiles/p2p")
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

	fmt.Printf("Sent file metadata to cloud node %s\n", Cloud_node_id)
	return nil
}

func sendSuccessConfirmation(transaction models.Transaction) {
	// Send decline to the target peer
	confirmationStream, err := CreateNewStream(DHT.Host(), transaction.TargetID, "/requestResponse/p2p")
	if err != nil {
		log.Printf("Error creating stream to target peer %s: %v", transaction.TargetID, err)
		return
	}
	defer confirmationStream.Close()

	transactionData, err := json.Marshal(transaction)
	if err != nil {
		fmt.Printf("error marshaling download request data: %v\n", err)
		return
	}

	_, err = confirmationStream.Write(transactionData)
	if err != nil {
		log.Printf("Error writing to stream: %v", err)
		return
	}
}

func SendProxyRequest(peerID string, wg *sync.WaitGroup) {
	defer wg.Done()

	// // Convert string peerID to libp2p PeerID
	// peerIDObj, err := peer.Decode(peerID)
	// if err != nil {
	// 	fmt.Printf("SendProxyRequest: Failed to decode PeerID %s: %v\n", peerID, err)
	// 	return
	// }

	// Open a new stream to the peer
	stream, err := CreateNewStream(DHT.Host(), peerID, "/proxy/metadata/1.0.0")
	if err != nil {
		fmt.Printf("SendProxyRequest: Failed to create stream to peer %s: %v\n", peerID, err)
		return
	}
	defer stream.Close()

	// Write a request message to the stream (if needed)
	request := models.RefreshRequest{
		Message:     "Requesting proxy metadata",
		RequesterID: Host.ID().String(),
		TargetID:    peerID,
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		fmt.Printf("SendProxyRequest: Failed to marshal request for peer %s: %v\n", peerID, err)
		return
	}

	_, err = stream.Write(requestBytes)
	if err != nil {
		fmt.Printf("SendProxyRequest: Failed to send request to peer %s: %v\n", peerID, err)
		return
	}

	// Read the response
	responseBytes := make([]byte, 4096)
	stream.SetReadDeadline(time.Now().Add(10 * time.Second))
	n, err := stream.Read(responseBytes)
	if err != nil {
		fmt.Printf("SendProxyRequest: Failed to read response from peer %s: %v\n", peerID, err)
		return
	}

	// Unmarshal the response into a Proxy object
	var proxy models.Proxy
	err = json.Unmarshal(responseBytes[:n], &proxy)
	if err != nil {
		fmt.Printf("SendProxyRequest: Failed to unmarshal response from peer %s: %v\n", peerID, err)
		return
	}

	// Append the proxy metadata to ProxyResponse
	ProxyResponse = append(ProxyResponse, proxy)
	fmt.Printf("SendProxyRequest: Successfully received proxy metadata from peer %s\n", peerID)
}

func handleHistoryStream(stream network.Stream) {
	defer stream.Close()

	var history []models.ProxyHistoryEntry
	decoder := json.NewDecoder(stream)
	if err := decoder.Decode(&history); err != nil {
		log.Printf("Error decoding history: %v", err)
		return
	}

	// Process the received history
	fmt.Printf("Received proxy history: %v\n", history)
	// Add your logic to store or process the history
}

func SendHistoryToHost(hostPeerIDStr string) error {

	// Decode the peer ID
	hostPeerID, err := peer.Decode(hostPeerIDStr)
	if err != nil {
		return fmt.Errorf("failed to decode peer ID: %v", err)
	}
	fmt.Printf("DECODED PEER ID: %s\n", hostPeerID)

	// Attempt to create a stream
	fmt.Printf("Attempting to create a stream to peer %s using protocol %s\n", hostPeerID, "/orcanet/p2p")
	stream, err := CreateNewStream(DHT.Host(), hostPeerIDStr, "/orcanet/p2p")
	if err != nil {
		fmt.Printf("ERROR: Failed to create stream to peer %s: %v\n", hostPeerID, err)
		return fmt.Errorf("failed to create stream: %v", err)
	}
	defer stream.Close()
	fmt.Printf("Stream successfully opened to peer %s\n", hostPeerID)

	// Send the proxy history
	fmt.Printf("Sending proxy history to peer %s\n", hostPeerID)
	encoder := json.NewEncoder(stream)
	fmt.Printf("Sending data: %+v\n", proxyHistory)

	if err := encoder.Encode(proxyHistory); err != nil {
		fmt.Printf("ERROR: Failed to encode and send history: %v\n", err)
		return fmt.Errorf("failed to encode and send history: %v", err)
	}
	fmt.Println("Proxy history sent successfully")

	return nil
}

// RECEIVING FUNCTIONS
func receivedHistory(node host.Host) {
	node.SetStreamHandler("/orcanet/p2p", func(stream network.Stream) {
		fmt.Println("Received a stream for /orcanet/p2p")
		defer stream.Close()

		decoder := json.NewDecoder(stream)
		var receivedHistory []models.ProxyHistoryEntry // Update YourHistoryType accordingly
		if err := decoder.Decode(&receivedHistory); err != nil {
			fmt.Printf("Error decoding received history: %v\n", err)
			return
		}
		fmt.Println("Received history:", receivedHistory)
	})
}
func receiveDownloadRequest(node host.Host) {
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

		utils.AddOrUpdateTransaction(request)

		// send file to requester if it exists
		if FileHashToPath[request.FileHash] != "" {
			fmt.Println("receivedownloadrequest: sending file")
			sendFile(node, request)
		} else {
			fmt.Println("receivedownloadrequest: decline")
			sendDecline(request)
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

		var declineMessage models.Transaction
		err = json.Unmarshal(data, &declineMessage)
		if err != nil {
			fmt.Printf("error unmarshalling file request: %v", err)
			return
		}
		declineMessage.Status = "declined"
		utils.AddOrUpdateTransaction(declineMessage)
	})
}

func receiveFile(node host.Host) error {
	fmt.Println("listening for file data")
	// listen for streams on "/sendFile/p2p"
	node.SetStreamHandler("/sendFile/p2p", func(s network.Stream) {
		defer s.Close()

		// read metadata - use metadata struct later
		buf := bufio.NewReader(s)

		// read in transaction details first
		transactionJSON, err := buf.ReadBytes('\n') // Read until newline
		if err != nil {
			log.Fatalf("Failed to read metadata: %v", err)
		}

		// Parse JSON metadata
		var transaction models.Transaction

		err = json.Unmarshal(transactionJSON, &transaction)
		if err != nil {
			log.Fatalf("Failed to unmarshal metadata: %v", err)
		}

		fmt.Printf("Received metadata: transactionID=%s\n", transaction.TransactionID)

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

		fmt.Printf("Received metadata: FileName=%s\n", metadata.NameWithExtension)

		// check if squidcoinFiles directory exists
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %v", dir, err)
		}

		// open file for writing
		outputPath := filepath.Join(dir, metadata.NameWithExtension)
		fmt.Println("receiveFile: outputPath", outputPath)
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

		transaction.Status = "complete"
		fmt.Println("receiveFile: transaction", transaction)
		utils.AddOrUpdateTransaction(transaction)

		// ProvideKey(GlobalCtx, DHT, metadata.Hash) // must be published - update dht with new provider
		err = UpdateFileInDHT(metadata)
		if err != nil {
			// is it a failure if the user receives the file but cannot be added to the dht?
			fmt.Errorf("failed to update dht metadata")
		}

		sendSuccessConfirmation(transaction)
	})
	return nil
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

func receiveSuccessConfirmation(node host.Host) {
	node.SetStreamHandler("/requestResponse/p2p", func(s network.Stream) {
		defer s.Close()
		buf := bufio.NewReader(s)

		// Read data until a newline character
		data, err := io.ReadAll(buf)
		if err != nil {
			if err == io.EOF {
				log.Printf("Stream closed by peer: %s", s.Conn().RemotePeer())
			} else {
				log.Printf("Error reading from stream: %v", err)
			}
			return
		}
		log.Printf("Raw data received: %s", data)

		var successMessage models.Transaction
		err = json.Unmarshal(data, &successMessage)
		if err != nil {
			fmt.Printf("error unmarshalling file request: %v", err)
			return
		}
		utils.AddOrUpdateTransaction(successMessage)
	})
}
func receiveProxies(node host.Host) {
	fmt.Println("Listening for proxy data")

	node.SetStreamHandler("/proxies/p2p", func(s network.Stream) {
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

		// Assuming proxies are serialized as JSON, we will parse them
		var proxies []models.Proxy // Replace ProxyData with the actual struct type for proxies
		err := json.Unmarshal(receivedData.Bytes(), &proxies)
		fmt.Println("Proxy data received:", proxies)
		if err != nil {
			log.Fatalf("Error unmarshaling received proxy data: %v", err)
		}

		// Assuming you have a global variable to store proxies
		Mutex.Lock()
		Proxies = append(Proxies[:0], proxies...) // Clear and update proxies
		Mutex.Unlock()

		// Notify other parts of your system that proxies have been received
		ProxiesSignal <- struct{}{}
		fmt.Println("Received proxies:", Proxies)
	})
}

// listen on streams
func setupStreams(node host.Host) {
	receiveDownloadRequest(node)
	receiveDecline(node)
	receiveFile(node)
	receivedHistory(node)
	receiveMarketplaceFiles(node)
	receiveSuccessConfirmation(node)
}
