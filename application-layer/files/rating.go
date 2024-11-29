package files

import (
	dht_kad "application-layer/dht"
	"application-layer/models"
	"encoding/json"
	"fmt"
	"net/http"
)

/*
	voting system is similar to that of stackoverflow and reddit
	upvote = +1, downvote = -1
*/

// Handle voting for both upvotes and downvotes
func handleVote(w http.ResponseWriter, r *http.Request) {
	fileHash := r.URL.Query().Get("fileHash")
	voteType := r.URL.Query().Get("voteType")

	if fileHash == "" || voteType == "" {
		http.Error(w, "File hash or vote type not provided", http.StatusBadRequest)
		return
	}

	if err := votingHelper(fileHash, voteType); err != nil {
		http.Error(w, fmt.Sprintf("Voting failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Vote '%s' recorded for file %s", voteType, fileHash)
}

func votingHelper(fileHash, voteType string) error {
	// Retrieve metadata from DHT
	data, err := dht_kad.DHT.GetValue(dht_kad.GlobalCtx, "/orcanet/"+fileHash)
	if err != nil {
		return fmt.Errorf("failed to retrieve file data: %v", err)
	}

	var metadata models.DHTMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return fmt.Errorf("failed to decode metadata: %v", err)
	}

	// Validate vote type
	if voteType != "upvote" && voteType != "downvote" {
		return fmt.Errorf("invalid vote type: %s", voteType)
	}

	// Check if the user is a provider
	provider, exists := metadata.Providers[dht_kad.PeerID]
	if !exists {
		return fmt.Errorf("user is not a provider and cannot vote")
	}

	// Ensure the user hasn't already voted
	if provider.Rating != "" {
		return fmt.Errorf("user has already voted")
	}

	// Apply vote logic
	if voteType == "upvote" {
		metadata.Upvote++
		metadata.Rating++
	} else if voteType == "downvote" {
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

	return dht_kad.DHT.PutValue(dht_kad.GlobalCtx, "/orcanet/"+fileHash, updatedData)
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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(metadata.Rating); err != nil {
		http.Error(w, "Failed to encode file rating", http.StatusInternalServerError)
	}
}
