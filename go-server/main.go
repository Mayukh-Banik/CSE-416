package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// MiningStatus represents the current mining state
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

// stops the mining process
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

		// Generate a new block
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

		// Delay between block generations
		time.Sleep(10 * time.Second)
	}
}

// main initializes the server, routing, and CORS
func main() {
	// set up btcd connection
	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:8334",
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
		AllowedOrigins:   []string{"http://localhost:3000"}, // Adjust as needed
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}

	// CORS handler
	handler := cors.New(corsOptions).Handler(router)

	// Create server with context for graceful shutdown
	server := &http.Server{
		Addr:    ":8081",
		Handler: handler,
	}

	// Set up signal handling for graceful shutdown
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Starting server on :8081")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %s\n", err.Error())
		}
	}()

	<-shutdownChan
	log.Println("Shutting down server...")

	// Gracefully shut down server with a timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
