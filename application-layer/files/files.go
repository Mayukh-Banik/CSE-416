package files

import (
	dht_kad "application-layer/dht"
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

	var files []FileMetadata
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

func saveOrUpdateFile(newFileData FileMetadata) (string, error) {
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

	// read in JSON file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read files.json: %v", err)
	}

	var files []FileMetadata
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
		files = append([]FileMetadata{newFileData}, files...)
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

	var requestBody FileMetadata
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
	}

	responseMsg := fmt.Sprintf("File %s successfully: %s", action, requestBody.Name)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseMsg))
	fmt.Println(responseMsg)
}

func publishFile(requestBody FileMetadata) {
	fmt.Println("publishing new file")

	// only one provider (uploader) for now bc it was just uploaded
	provider := []Provider{
		{PeerID: dht_kad.PeerID, IsActive: true, Fee: requestBody.Fee},
	}

	dhtMetadata := DHTMetadata{
		Name:        requestBody.Name,
		Type:        requestBody.Type,
		Size:        requestBody.Size,
		Description: requestBody.Description,
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

// get list of providers for a file hash
// func handleGetProvidersByFileHash(w http.ResponseWriter, r *http.Request) {
// 	// Get file hash from the query parameters (instead of the body)
// 	fileHash := r.URL.Query().Get("val")
// 	fmt.Println("filehash:", fileHash)

// 	if fileHash == "" {
// 		http.Error(w, "File hash not provided", http.StatusBadRequest)
// 		return
// 	}

// 	// data := []byte(fileHash)
// 	// hash := sha256.Sum256(data)
// 	// mh, err := multihash.EncodeName(hash[:], "sha2-256")

// 	// if err != nil {
// 	// 	fmt.Printf("Error encoding multihash: %v\n", err)
// 	// 	http.Error(w, "Error encoding hash", http.StatusInternalServerError)
// 	// 	return
// 	// }
// 	data, err := dht_kad.DHT.GetValue(dht_kad.GlobalCtx, "/orcanet/"+fileHash)
// 	if err != nil {
// 		return
// 	}

// 	// Create an instance of FileMetadata to hold the decoded data
// 	var metadata DHTMetadata

// 	// Unmarshal the JSON data into the struct
// 	err = json.Unmarshal(data, &metadata)
// 	if err != nil {
// 		return
// 	}
// 	fmt.Println(metadata)

// 	providers := metadata.Providers
// 	fmt.Println(providers)

// 	// providers, err := dht_kad.ProviderStore.GetProviders(dht_kad.GlobalCtx, mh)
// 	// fmt.Println(err)

// 	// var providerList []string
// 	// for _, provider := range providers {
// 	// 	fmt.Println("Provider PeerID:", provider.ID)
// 	// 	providerList = append(providerList, provider.ID.String())
// 	// 	fmt.Println("Provider PeerID:", provider.ID.String())
// 	// }

// 	// Send the list of providers as JSON response
// 	w.Header().Set("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(providers); err != nil {
// 		http.Error(w, "Failed to encode providers", http.StatusInternalServerError)
// 	}
// }

func handleGetProvidersByFileHash(w http.ResponseWriter, r *http.Request) {
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
	var metadata DHTMetadata

	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(data, &metadata)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding file metadata: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Println(metadata)

	// Send the entire metadata (including providers) as JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(metadata); err != nil {
		http.Error(w, "Failed to encode file metadata", http.StatusInternalServerError)
	}
}
