package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// MiningStatus holds the current status of mining
type MiningStatus struct {
	MinedBlocks    int    `json:"minedBlocks"`
	LastMinedBlock string `json:"lastMinedBlock"`
	IsMining       bool   `json:"isMining"`
}

var (
	miningStatus MiningStatus
	mu           sync.Mutex
)

// fetch current mining status
func getMiningStatus(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(miningStatus)
}

// begins mining process
func startMining(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	if miningStatus.IsMining {
		http.Error(w, "Mining is already in progress", http.StatusBadRequest)
		return
	}

	miningStatus.IsMining = true
	miningStatus.MinedBlocks = 0
	miningStatus.LastMinedBlock = ""
	go mine()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(miningStatus)
}

// stopMining handles POST requests to stop the mining process
func stopMining(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	if !miningStatus.IsMining {
		http.Error(w, "Mining is not currently active", http.StatusBadRequest)
		return
	}

	miningStatus.IsMining = false

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(miningStatus)
}

// mine simulates the mining process by incrementing mined blocks periodically
func mine() {
	for {
		mu.Lock()
		if !miningStatus.IsMining {
			mu.Unlock()
			return
		}
		miningStatus.MinedBlocks++
		miningStatus.LastMinedBlock = time.Now().Format(time.RFC3339)
		mu.Unlock()

		// Simulate mining delay
		time.Sleep(5 * time.Second)
	}
}

func main() {
	// Initialize default mining status
	miningStatus = MiningStatus{
		MinedBlocks:    0,
		LastMinedBlock: "",
		IsMining:       false,
	}

	// Create a new router
	router := mux.NewRouter()

	// Define API routes
	router.HandleFunc("/api/mining-status", getMiningStatus).Methods("GET")
	router.HandleFunc("/api/start-mining", startMining).Methods("POST")
	router.HandleFunc("/api/stop-mining", stopMining).Methods("POST")

	// Configure CORS
	corsOptions := cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Frontend origin
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}

	// Create a CORS handler
	handler := cors.New(corsOptions).Handler(router)

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
