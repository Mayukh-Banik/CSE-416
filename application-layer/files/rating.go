package files

import (
	dht_kad "application-layer/dht"
	"application-layer/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

/*
	voting system is similar to that of stackoverflow and reddit
	upvote = +1, downvote = -1
*/

// Handle voting for both upvotes and downvotes
func handleVote(w http.ResponseWriter, r *http.Request) {
	log.Println("in handleVote")

	fileHash := r.URL.Query().Get("fileHash")
	voteType := r.URL.Query().Get("voteType")
	log.Printf("file hash: %s | vote type: %s\n", fileHash, voteType)

	if fileHash == "" || voteType == "" {
		http.Error(w, `{"error": "File hash or vote type not provided"}`, http.StatusBadRequest)
		return
	}

	if err := votingHelper(fileHash, voteType); err != nil {
		log.Printf("handleVote error: %v", err)
		http.Error(w, fmt.Sprintf(`{"error": "Voting failed: %v"}`, err), http.StatusInternalServerError)
		return
	}

	log.Println("handleVote: successfully voted!")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "Vote '%s' recorded for file %s"}`, voteType, fileHash)
}

func votingHelper(fileHash string, voteType string) error {
	// Retrieve metadata from DHT
	fmt.Println("in voting helper")

	// check if user has already voted
	data, err := dht_kad.DHT.GetValue(dht_kad.GlobalCtx, "/orcanet/"+fileHash)
	if err != nil {
		fmt.Println("error retrieving file data from dht")
		return fmt.Errorf("failed to retrieve file data: %v", err)
	}

	var metadata models.DHTMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return fmt.Errorf("failed to decode metadata: %v", err)
	}

	fmt.Printf("votingHelper: file metadata from dht: %v\n", metadata)

	// Validate vote type
	if voteType != "upvote" && voteType != "downvote" {
		fmt.Println("votingHelper: vote type:", voteType)
		return fmt.Errorf("invalid vote type: %s", voteType)
	}

	// Check if the user is a provider
	provider, exists := metadata.Providers[dht_kad.PeerID]
	if !exists {
		fmt.Println("user is not a provider")
		return fmt.Errorf("user is not a provider and cannot vote")
	}

	// Ensure the user hasn't already voted
	if provider.Rating != "" {
		fmt.Println("user has already voted")
		return fmt.Errorf("user has already voted")
	}

	// Apply vote logic
	fmt.Println("voting helper: vote type:", voteType)
	if voteType == "upvote" {
		fmt.Println("votingHelper: upvoting")
		metadata.Upvote++
		metadata.Rating++
	} else if voteType == "downvote" {
		fmt.Println("votingHelper: downvoting")
		metadata.Rating--
		metadata.Downvote++
	}

	// Update provider's vote status
	provider.Rating = voteType
	metadata.Providers[dht_kad.PeerID] = provider
	metadata.NumRaters++ // Increment the number of raters

	// Store updated metadata in DHT
	updatedData, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to encode updated metadata: %v", err)
	}
	fmt.Println("votingHelper: metadata after updating vote: ", metadata)

	err = dht_kad.DHT.PutValue(dht_kad.GlobalCtx, "/orcanet/"+fileHash, updatedData)
	if err != nil {
		fmt.Println("votingHelper: error publishing file to DHT", err)
	}

	updateRatingLocally(fileHash, voteType)

	return nil
}

func handleGetRating(w http.ResponseWriter, r *http.Request) {
	fileHash := r.URL.Query().Get("fileHash")

	if fileHash == "" {
		http.Error(w, "File hash not provided", http.StatusBadRequest)
		return
	}

	data, err := dht_kad.DHT.GetValue(dht_kad.GlobalCtx, "/orcanet/"+fileHash)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving file data: %v", err), http.StatusInternalServerError)
		return
	}

	var metadata models.DHTMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding file metadata: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Println("rating for file hash: ", fileHash, metadata.Rating)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(metadata.Rating); err != nil {
		http.Error(w, "Failed to encode file rating", http.StatusInternalServerError)
	}
}

func updateRatingLocally(fileHash string, voteType string) error {
	// check if directory and file exist
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return fmt.Errorf("utils directory doesnt exist --> cannot vote %v", err)
	}

	if _, err := os.Stat(DownloadedFilePath); os.IsNotExist(err) {
		return fmt.Errorf("downloaded file path does not exist --> cannot vote: %v", err)
	}

	// read in JSON file
	data, err := os.ReadFile(DownloadedFilePath)
	if err != nil {
		return fmt.Errorf("failed to read downloadedFiles.json: %v", err)
	}

	var files []models.FileMetadata
	if err := json.Unmarshal(data, &files); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	for i := range files {
		if files[i].Hash == fileHash {
			files[i].HasVoted = true
			files[i].VoteType = voteType
			break
		}
	}

	// convert updated list of files back to JSON
	updatedData, err := json.MarshalIndent(files, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	if err := os.WriteFile(DownloadedFilePath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write updated data to files.json: %v", err)
	}

	return nil
}
