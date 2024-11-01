package main

import (
	"context"
<<<<<<< HEAD
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
=======
	"log"
	"net/http"
	"os"
>>>>>>> dianne
	"time"

	"go-server/routes"
	"go-server/utils"

	"go-server/services"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

<<<<<<< HEAD
// MiningStatus represents the current mining state
type MiningStatus struct {
	MinedBlocks    int    `json:"minedBlocks"`
	LastMinedBlock string `json:"lastMinedBlock"`
	IsMining       bool   `json:"isMining"`
}
=======
// // current mining status
// type MiningStatus struct {
// 	MinedBlocks    int    `json:"minedBlocks"`
// 	LastMinedBlock string `json:"lastMinedBlock"`
// 	IsMining       bool   `json:"isMining"`
// }
>>>>>>> dianne

// var (
// 	miningStatus MiningStatus
// 	client       *rpcclient.Client
// 	mu           sync.Mutex
// )

// // fetch current mining status
// func getMiningStatus(w http.ResponseWriter, r *http.Request) {
// 	mu.Lock()
// 	defer mu.Unlock()

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(miningStatus)
// }

// // begins mining process
// func startMining(w http.ResponseWriter, r *http.Request) {
// 	mu.Lock()
// 	defer mu.Unlock()

// 	if miningStatus.IsMining {
// 		http.Error(w, "Mining is already in progress", http.StatusBadRequest)
// 		return
// 	}

// 	miningStatus.IsMining = true
// 	miningStatus.MinedBlocks = 0
// 	miningStatus.LastMinedBlock = ""
// 	go mine()

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(miningStatus)
// }

<<<<<<< HEAD
// stops the mining process
func stopMining(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
=======
// // stopMining handles POST requests to stop the mining process
// func stopMining(w http.ResponseWriter, r *http.Request) {
// 	mu.Lock()
// 	defer mu.Unlock()
>>>>>>> dianne

// 	if !miningStatus.IsMining {
// 		http.Error(w, "Mining is not currently active", http.StatusBadRequest)
// 		return
// 	}

<<<<<<< HEAD
	miningStatus.IsMining = false
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(miningStatus)
}
=======
// 	miningStatus.IsMining = false

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(miningStatus)
// }
>>>>>>> dianne

// mine simulates the mining process by incrementing mined blocks periodically
// func mine() {
// 	for {
// 		mu.Lock()
// 		if !miningStatus.IsMining {
// 			mu.Unlock()
// 			return
// 		}

<<<<<<< HEAD
		// Generate a new block
		blockHashes, err := client.Generate(1)
		if err != nil {
			log.Printf("Failed to generate block: %v", err)
			miningStatus.IsMining = false
			mu.Unlock()
			return
		}
=======
// 		// implement our own block generation
// 		blockHashes, err := client.Generate(1)
// 		if err != nil {
// 			log.Printf("Failed to generate block: %v", err)
// 			miningStatus.IsMining = false
// 			mu.Unlock()
// 			return
// 		}
>>>>>>> dianne

// 		miningStatus.MinedBlocks++
// 		miningStatus.LastMinedBlock = blockHashes[0].String()
// 		mu.Unlock()

<<<<<<< HEAD
		// Delay between block generations
		time.Sleep(10 * time.Second)
	}
}
=======
// 		// mining delay
// 		time.Sleep(10 * time.Second)
// 	}
// }
>>>>>>> dianne

// main initializes the server, routing, and CORS
func main() {
	// set up btcd connection
<<<<<<< HEAD
	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:8334",
		User:         "user",
		Pass:         "password",
		HTTPPostMode: true,
		DisableTLS:   true,
	}
=======
	// connCfg := &rpcclient.ConnConfig{
	// 	Host:         "localhost:8334", // btcd RPC server
	// 	User:         "user",
	// 	Pass:         "password",
	// 	HTTPPostMode: true,
	// 	DisableTLS:   true,
	// }
>>>>>>> dianne

	// defer func() {
	//     if err := utils.MongoClient.Disconnect(context.TODO()); err != nil {
	//         log.Fatalf("Error disconnecting from MongoDB: %v", err)
	//     }
	// }()

	// Load the .env file to get environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// default mining status
	// miningStatus = MiningStatus{
	// 	MinedBlocks:    0,
	// 	LastMinedBlock: "",
	// 	IsMining:       false,
	// }

	// Initialize MongoDB connection using the environment variable
	MONGODB_URI := os.Getenv("MONGODB_URI")
	if MONGODB_URI == "" {
		log.Fatalf("MONGODB_URI environment variable is not set")
	}

	utils.ConnectMongo(MONGODB_URI)

	defer func() {
		if err := utils.MongoClient.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Error disconnecting from MongoDB: %v", err)
		}
	}()

<<<<<<< HEAD
	// Configure CORS
	corsOptions := cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Adjust as needed
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
=======
	// Periodically clean up expired Challenges
	go func() {
		for {
			time.Sleep(time.Minute)
			services.CleanupExpiredChallenges(time.Minute * 5)
		}
	}()

	// Initialize routes
	router := routes.InitRoutes()

	// Add CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Frontend origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
>>>>>>> dianne
		AllowCredentials: true,
	})

<<<<<<< HEAD
	// CORS handler
	handler := cors.New(corsOptions).Handler(router)

	// Create server with context for graceful shutdown
	server := &http.Server{
		Addr:    ":8081",
		Handler: handler,
=======
	// Use the CORS handler directly in the `ListenAndServe` function
	log.Println("Server starting on port 8080...")
	err = http.ListenAndServe(":8080", c.Handler(router)) // Pass the handler here
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
>>>>>>> dianne
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
