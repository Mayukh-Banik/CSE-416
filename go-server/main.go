package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// current mining status
type MiningStatus struct {
	MinedBlocks    int    `json:"minedBlocks"`
	LastMinedBlock string `json:"lastMinedBlock"`
	IsMining       bool   `json:"isMining"`
}

var (
	miningStatus MiningStatus
	client       *rpcclient.Client
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

		// implement our own block generation
		blockHashes, err := client.Generate(1)
		if err != nil {
			log.Printf("Failed to generate block: %v", err)
			miningStatus.IsMining = false
			mu.Unlock()
			return
		}

		miningStatus.MinedBlocks++
		miningStatus.LastMinedBlock = blockHashes[0].String()
		mu.Unlock()

		// mining delay
		time.Sleep(10 * time.Second)
	}
}

func main() {
	// set up btcd connection
	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:8334", // btcd RPC server
		User:         "user",
		Pass:         "password",
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	var err error
	client, err = rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatalf("Error creating btcd client: %v", err)
	}
	defer client.Shutdown()

	// default mining status
	miningStatus = MiningStatus{
		MinedBlocks:    0,
		LastMinedBlock: "",
		IsMining:       false,
	}

	router := mux.NewRouter()

	// API routes
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

	// CORS handler
	handler := cors.New(corsOptions).Handler(router)

	// server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
