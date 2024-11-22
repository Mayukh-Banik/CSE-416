package files

import (
	dht_kad "application-layer/dht"
	"application-layer/models"
	"application-layer/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
)

var (
	dirPath            = filepath.Join("..", "utils")
	UploadedFilePath   = filepath.Join(dirPath, "files.json")
	DownloadedFilePath = filepath.Join(dirPath, "downloadedFiles.json")
)

// fetch all uploaded files from JSON file
func getUploadedFiles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("tring to fetch user's uploaded files")
	file, err := os.ReadFile(UploadedFilePath)

	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	var files []models.FileMetadata
	if err := json.Unmarshal(file, &files); err != nil {
		http.Error(w, "Failed to parse files data", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(files); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestBody models.FileMetadata
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	action, err := utils.SaveOrUpdateFile(requestBody, dirPath, UploadedFilePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if action == "added" {
		PublishFile(requestBody)
		// fmt.Printf("new file %v | path: %v\n", requestBody.Hash, requestBody.Path)

		curDirectory, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting cur directory %v\n", err)
			return
		}

		newPath := filepath.Join(curDirectory, "../squidcoinFiles", requestBody.Name)
		dht_kad.FileMapMutex.Lock()
		dht_kad.FileHashToPath[requestBody.Hash] = newPath
		dht_kad.FileMapMutex.Unlock()
		// ///dht_kad.FileHashToPath[requestBody.Hash] = requestBody.Path // fix getting file path

	}

	responseMsg := fmt.Sprintf("File %s successfully: %s", action, requestBody.Name)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseMsg))
	fmt.Println(responseMsg)
}

func PublishFile(requestBody models.FileMetadata) {
	fmt.Println("publishing new file")

	// Retrieve the current metadata for the file, if it exists
	var currentMetadata models.DHTMetadata
	existingData, err := dht_kad.DHT.GetValue(dht_kad.GlobalCtx, "/orcanet/"+requestBody.Hash)
	if err == nil { // If data exists, unmarshal it
		err = json.Unmarshal(existingData, &currentMetadata)
		if err != nil {
			log.Fatal("Failed to unmarshal existing DHTMetadata:", err)
		}
	} else {
		// If no existing metadata, initialize a new DHTMetadata
		currentMetadata = models.DHTMetadata{
			Name:        requestBody.Name,
			Type:        requestBody.Type,
			Size:        requestBody.Size,
			Description: requestBody.Description,
			CreatedAt:   requestBody.CreatedAt,
			Reputation:  requestBody.Reputation,
		}
	}

	// Add the new provider to the list of current providers
	provider := models.Provider{
		PeerID:   dht_kad.PeerID,
		PeerAddr: dht_kad.DHT.Host().Addrs()[0].String(),
		IsActive: true,
		Fee:      requestBody.Fee,
	}

	currentMetadata.Providers = append(currentMetadata.Providers, provider)

	// Marshal the updated metadata
	dhtMetadataBytes, err := json.Marshal(currentMetadata)
	if err != nil {
		log.Fatal("Failed to marshal updated DHTMetadata:", err)
	}

	// Store the updated metadata in the DHT
	err = dht_kad.DHT.PutValue(dht_kad.GlobalCtx, "/orcanet/"+requestBody.Hash, dhtMetadataBytes)
	if err != nil {
		log.Fatal("failed to register updated file to dht")
	}
	fmt.Println("successfully updated file to dht with new provider", requestBody.Hash)

	// Begin providing ourselves as a provider for that file
	dht_kad.ProvideKey(dht_kad.GlobalCtx, dht_kad.DHT, requestBody.Hash)
}

// bug
func handleGetFileByHash(w http.ResponseWriter, r *http.Request) {
	// Get file hash from the query parameters (instead of the body)
	fmt.Println("getting file by hash")
	fileHash := r.URL.Query().Get("val")
	fmt.Println("filehash:", fileHash)

	if fileHash == "" {
		http.Error(w, "File hash not provided", http.StatusBadRequest)
		return
	}
	fmt.Println("before getting file data from dht:")

	// Retrieve the file data from the DHT using the file hash
	data, err := dht_kad.DHT.GetValue(dht_kad.GlobalCtx, "/orcanet/"+fileHash)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving file data: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Println("file data from dht:", data)

	// Create an instance of FileMetadata to hold the decoded data
	var metadata models.DHTMetadata

	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(data, &metadata)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding file metadata: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Println("file requested metadata: ", metadata)

	// Send the entire metadata (including providers) as JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(metadata); err != nil {
		http.Error(w, "Failed to encode file metadata", http.StatusInternalServerError)
	}
}

// deleting a file requires removing it from the json and "removing" the node as a provider of the file
func deleteFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request methods", http.StatusMethodNotAllowed)
		return
	}

	hash := r.URL.Query().Get("hash")
	if hash == "" {
		http.Error(w, "file hash not provided", http.StatusBadRequest)
		return
	}

	isDownloaded := r.URL.Query().Get("isDownloaded")
	if isDownloaded == "" {
		http.Error(w, "json file not provided", http.StatusBadRequest)
		return
	}

	var filePath string
	if isDownloaded == "true" {
		filePath = DownloadedFilePath
	} else {
		filePath = UploadedFilePath
	}

	// update so it works for both uplaoded and downloaded files
	action, err := deleteFileFromJSON(hash, filePath)
	if err != nil {
		http.Error(w, fmt.Sprint("failed to delete file from file", err), http.StatusInternalServerError)
		return
	}

	dht_kad.FileMapMutex.Lock()
	delete(dht_kad.FileHashToPath, hash) // delete from map of file hash to file path
	dht_kad.FileMapMutex.Unlock()

	response := map[string]string{
		"status":  action,
		"message": "File deletion successful",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// removeFileFromJSON removes the file entry with the given hash from the JSON file
func deleteFileFromJSON(fileHash string, filePath string) (string, error) {
	// Read the JSON file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read files.json: %v", err)
	}

	var files []models.FileMetadata
	if err := json.Unmarshal(data, &files); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Find and remove the file by hash
	updatedFiles := []models.FileMetadata{}
	fileFound := false
	for _, file := range files {
		if file.Hash != fileHash {
			updatedFiles = append(updatedFiles, file)
		} else {
			fileFound = true
		}
	}

	if !fileFound {
		return "not found", fmt.Errorf("file not found in files.json")
	}

	// Write the updated files list back to files.json
	updatedData, err := json.MarshalIndent(updatedFiles, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %v", err)
	}

	if err := os.WriteFile(filePath, updatedData, 0644); err != nil {
		return "", fmt.Errorf("failed to write updated data to files.json: %v", err)
	}

	return "deleted", nil
}

// functions below are used in marketplace to get all dht files
func getAdjacentNodeFilesMetadata(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Trying to get adjacent node files in backend")
	dht_kad.RefreshResponse = dht_kad.RefreshResponse[:0]

	relayNode := "12D3KooWDpJ7As7BWAwRMfu1VU2WCqNjvq387JEYKDBj4kx6nXTN"
	bootstrapNode := "12D3KooWQd1K1k8XA9xVEzSAu7HUCodC7LJB6uW5Kw4VwkRdstPE"

	// Retrieve connected peers
	adjacentNodes := dht_kad.Host.Peerstore().Peers()
	fmt.Println("Connected peers:", adjacentNodes)

	var sendWG sync.WaitGroup
	var responseWG sync.WaitGroup

	// Convert PeerID to strings
	// var peers []string
	for _, peer := range adjacentNodes {
		peerID := peer.String()
		if peerID != relayNode && peerID != bootstrapNode && peerID != dht_kad.PeerID && nodeSupportRefreshStreams(peer) {
			sendWG.Add(1)
			responseWG.Add(1)
			go func(peerID string) {
				defer responseWG.Done()
				go dht_kad.SendRefreshFilesRequest(peerID, &sendWG)
			}(peerID)
			// peers = append(peers, peer.String())
		}
	}

	sendWG.Wait()
	responseWG.Wait()

	// // create stream to every adjacent node
	// for _, peer := range peers {
	// 	if peer != "12D3KooWDpJ7As7BWAwRMfu1VU2WCqNjvq387JEYKDBj4kx6nXTN" &&
	// 		peer != "12D3KooWQd1K1k8XA9xVEzSAu7HUCodC7LJB6uW5Kw4VwkRdstPE" &&
	// 		peer != dht_kad.PeerID { // cannot connect to relay node, bootstrap node, or self
	// 		dht_kad.SendRefreshFilesRequest(peer)
	// 	}
	// }

	// dht_kad.SendRefreshFilesRequest("12D3KooWFZ8nwUD3cxtqLHvord4cXU1M7vcoUoEwrouADQskxsVJ")
	<-time.After(3 * time.Second)

	fmt.Println("getAdjacentNodeFilesMetadata: received everyone's uploaded files: ", dht_kad.RefreshResponse)
	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Make sure the status is OK
	fmt.Println("getAdjacentNodeFilesMetadata: back to frontend...")

	// Encode response
	/**
	if err := json.NewEncoder(w).Encode(peers); err != nil {
		http.Error(w, "Failed to encode adjacent nodes", http.StatusInternalServerError)
	}
	*/

	if err := json.NewEncoder(w).Encode(dht_kad.RefreshResponse); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
	fmt.Println("getAdjacentNodeFilesMetadata: back to frontend...")

}

func nodeSupportRefreshStreams(peerID peer.ID) bool {
	supportSendRefreshRequest := false
	supportSendRefreshResponse := false

	protocols, _ := dht_kad.Host.Peerstore().GetProtocols(peerID)
	fmt.Printf("protocols supported by peer %v: %v\n", peerID, protocols)

	for _, protocol := range protocols {
		if protocol == "/sendRefreshRequest/p2p" {
			supportSendRefreshRequest = true
		} else if protocol == "/sendRefreshResponse/p2p" {
			supportSendRefreshResponse = true
		}
	}
	return supportSendRefreshRequest && supportSendRefreshResponse
}

// republish files in the dht incase the TTL expired - called upon successful login
func republishFiles(filePath string) {
	// Open the JSON file
	file, err := os.Open(filePath) // Replace "files.json" with your file name
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No files to republish, node has not uploaded any files.")
			return
		}
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Read file content
	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Parse JSON into structs
	var files []models.FileMetadata
	err = json.Unmarshal(byteValue, &files)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Print parsed structs
	for _, file := range files {
		fmt.Printf("Republishing File: %+v\n", file)
		PublishFile(file)
	}

}
