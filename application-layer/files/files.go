package files

import (
	dht_kad "application-layer/dht"
	"application-layer/models"
	"application-layer/utils"
	"encoding/json"
	"errors"
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
	dirPath             = filepath.Join("..", "utils")
	UploadedFilePath    = filepath.Join(dirPath, "files.json")
	DownloadedFilePath  = filepath.Join(dirPath, "downloadedFiles.json")
	transactionFilePath = filepath.Join(dirPath, "transactionFiles.json")
	FileCopyPath        = filepath.Join("..", "squidcoinFiles")
	republishMutex      sync.Mutex
	republished         = false
)

// fetch all uploaded files from JSON file
func getFiles(w http.ResponseWriter, r *http.Request) {
	fileType := r.URL.Query().Get("file")
	fmt.Printf("trying to fetch user's %v files \n", fileType)

	var filePath string
	if fileType == "uploaded" {
		fmt.Println("getting uploaded files")
		filePath = UploadedFilePath
	} else {
		fmt.Println("getting downloaded files")
		filePath = DownloadedFilePath
	}

	file, err := os.ReadFile(filePath)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("No %s files found for user\n", fileType)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]models.FileMetadata{}) // Return empty array
			return
		}
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	var files []models.FileMetadata
	if err := json.Unmarshal(file, &files); err != nil {
		http.Error(w, "Failed to parse files data", http.StatusInternalServerError)
		return
	}

	republishMutex.Lock()
	defer republishMutex.Unlock()
	fmt.Println("republished bool: ", republished)
	if !republished {
		if fileType == "uploaded" {
			republishFiles(UploadedFilePath)
		} else {
			republishFiles(DownloadedFilePath)
		}
		republished = true
	}

	fmt.Println("fileHashToPath:", dht_kad.FileHashToPath)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(files); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// add new file or update existing file
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

	var filePath string
	fmt.Println("original uploader: ", requestBody.OriginalUploader)
	if requestBody.OriginalUploader {
		fmt.Println("updating uploaded file path")
		filePath = UploadedFilePath
	} else {
		fmt.Println("updating downloaded file path")
		filePath = DownloadedFilePath
	}
	fmt.Println("filePath: ", filePath)
	action, err := utils.SaveOrUpdateFile(requestBody, dirPath, filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// update publishing info - reflects if user is currently providing or not
	PublishFile(requestBody)

	responseMsg := fmt.Sprintf("File %s successfully: %s", action, requestBody.Name)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseMsg))
	fmt.Println(responseMsg)
}

func PublishFile(requestBody models.FileMetadata) {
	fmt.Println("publishing new file")

	err := dht_kad.UpdateFileInDHT(requestBody)
	if err != nil {
		fmt.Printf("unable to update file in the dht %w\n", err)
		return
	}

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("error getting currentDir %v\n", err)
		return
	}

	newPath := filepath.Join(currentDir, "../squidcoinFiles", requestBody.NameWithExtension)
	dht_kad.FileMapMutex.Lock()
	dht_kad.FileHashToPath[requestBody.Hash] = newPath
	fmt.Println("PublishFile: fileHashToPath: ", dht_kad.FileHashToPath)
	dht_kad.FileMapMutex.Unlock()
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

	originalUploader := r.URL.Query().Get("originalUploader")
	if originalUploader == "" {
		http.Error(w, "did not specify is user uplaoded the file or downloaded it", http.StatusBadRequest)
		return
	}
	fmt.Println("user is original uploader?", originalUploader)

	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "file name not provided", http.StatusBadRequest)
		return
	}
	fmt.Println("trying to delete file", name)

	var filePath string
	if originalUploader == "true" {
		filePath = UploadedFilePath
	} else {
		filePath = DownloadedFilePath
	}

	// VERY BUGGY
	err := deleteFileContent(name) // currently using file name but we should switch to hash
	if err != nil {
		http.Error(w, fmt.Sprint("failed to delete file from squidcoinFiles", err), http.StatusInternalServerError)
		return
	}

	dht_kad.FileMapMutex.Lock()
	delete(dht_kad.FileHashToPath, hash) // delete from map of file hash to file path
	dht_kad.FileMapMutex.Unlock()

	// update so it works for both uploaded and downloaded files
	action, err := deleteFileFromJSON(hash, filePath)
	if err != nil { // will still show up but will not be able to provide
		http.Error(w, fmt.Sprint("failed to delete file json file", err), http.StatusInternalServerError)
		return
	}

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
	fmt.Println("deleting file from JSON:", filePath)
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

// currently using file name but user can download files of the same name
// from different providers so we have to switch to file hash
func deleteFileContent(hash string) error {
	filePath := filepath.Join(FileCopyPath, hash)

	// Attempt to delete the file
	err := os.Remove(filePath)
	if err != nil {
		fmt.Printf("Failed to delete file: %v\n", err)
		return err
	}

	fmt.Println("File deleted successfully")
	return nil
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

	<-time.After(10 * time.Second)

	fmt.Println("getAdjacentNodeFilesMetadata: received everyone's uploaded files: ", dht_kad.RefreshResponse)
	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Make sure the status is OK
	fmt.Println("getAdjacentNodeFilesMetadata: back to frontend...")

	// Encode response
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
	fmt.Println("republishing files in ", filePath)
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

// transactions page
func getTransactions(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getting transaction history")
	transactionFile, err := os.ReadFile(transactionFilePath)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("transactionFiles.json could not be found")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]models.Transaction{}) // Return empty array
			return
		}
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	var transactions []models.Transaction
	if err := json.Unmarshal(transactionFile, &transactions); err != nil {
		http.Error(w, "Failed to parse transaction files data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
