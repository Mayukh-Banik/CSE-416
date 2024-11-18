package files

import (
	dht_kad "application-layer/dht"
	"application-layer/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	dirPath  = filepath.Join("..", "utils")
	filePath = filepath.Join(dirPath, "files.json")
)

// fetch all uploaded files from JSON file
func getUploadedFiles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("tring to fetch user's uploaded files")
	file, err := os.ReadFile(filePath)

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

func saveOrUpdateFile(newFileData models.FileMetadata) (string, error) {
	// check if directory and file exist
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.Mkdir(dirPath, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create utils directory: %v", err)
		}
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := os.WriteFile(filePath, []byte("[]"), 0644); err != nil {
			return "", fmt.Errorf("failed to create files.json: %v", err)
		}
	}

	fmt.Println("file to be added/updated: ", newFileData)

	// read in JSON file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read files.json: %v", err)
	}

	var files []models.FileMetadata
	if err := json.Unmarshal(data, &files); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %v", err)
	}

	// update if file is already in JSON file
	isUpdated := false
	for i := range files {
		if files[i].Hash == newFileData.Hash {
			files[i] = newFileData
			isUpdated = true
			break
		}
	}

	// add file if not already in JSON file
	if !isUpdated {
		files = append([]models.FileMetadata{newFileData}, files...)
	}

	// convert updated list of files back to JSON
	updatedData, err := json.MarshalIndent(files, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %v", err)
	}

	if err := os.WriteFile(filePath, updatedData, 0644); err != nil {
		return "", fmt.Errorf("failed to write updated data to files.json: %v", err)
	}

	action := "updated"
	if !isUpdated {
		action = "added"
	}

	return action, nil
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

	action, err := saveOrUpdateFile(requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if action == "added" {
		publishFile(requestBody)
		// fmt.Printf("new file %v | path: %v\n", requestBody.Hash, requestBody.Path)
		dht_kad.FileHashToPath[requestBody.Hash] = requestBody.Path // fix getting file path
	}

	responseMsg := fmt.Sprintf("File %s successfully: %s", action, requestBody.Name)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseMsg))
	fmt.Println(responseMsg)
}

func publishFile(requestBody models.FileMetadata) {
	fmt.Println("publishing new file")

	// only one provider (uploader) for now bc it was just uploaded
	provider := []models.Provider{
		{PeerID: dht_kad.PeerID, PeerAddr: dht_kad.DHT.Host().Addrs()[0].String(), IsActive: true, Fee: requestBody.Fee},
	}

	dhtMetadata := models.DHTMetadata{
		Name:        requestBody.Name,
		Type:        requestBody.Type,
		Size:        requestBody.Size,
		Description: requestBody.Description,
		CreatedAt:   requestBody.CreatedAt,
		Reputation:  requestBody.Reputation,
		Providers:   provider,
	}

	dhtMetadataBytes, err := json.Marshal(dhtMetadata)
	if err != nil {
		log.Fatal("Failed to marshal DHTMetadata:", err)
	}

	err = dht_kad.DHT.PutValue(dht_kad.GlobalCtx, "/orcanet/"+requestBody.Hash, dhtMetadataBytes)
	if err != nil {
		log.Fatal("failed to register file to dht")
	}
	fmt.Println("successfully registered file to dht", requestBody.Hash)

	// Begin providing ourselves as a provider for that file
	dht_kad.ProvideKey(dht_kad.GlobalCtx, dht_kad.DHT, requestBody.Hash)
}

func handleGetFileByHash(w http.ResponseWriter, r *http.Request) {
	// Get file hash from the query parameters (instead of the body)
	fileHash := r.URL.Query().Get("val")
	fmt.Println("filehash:", fileHash)

	if fileHash == "" {
		http.Error(w, "File hash not provided", http.StatusBadRequest)
		return
	}

	// Retrieve the file data from the DHT using the file hash
	data, err := dht_kad.DHT.GetValue(dht_kad.GlobalCtx, "/orcanet/"+fileHash)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving file data: %v", err), http.StatusInternalServerError)
		return
	}

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

	action, err := deleteFileFromJSON(hash)
	if err != nil {
		http.Error(w, fmt.Sprint("failed to delete file from file", err), http.StatusInternalServerError)
		return
	}

	delete(dht_kad.FileHashToPath, hash) // delete from map of file hash to file path

	response := map[string]string{
		"status":  action,
		"message": "File deletion successful",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// removeFileFromJSON removes the file entry with the given hash from the JSON file
func deleteFileFromJSON(fileHash string) (string, error) {
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
