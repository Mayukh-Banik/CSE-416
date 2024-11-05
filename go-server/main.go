package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"go-server/routes"
	"go-server/utils"

	"go-server/services"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

// // current mining status
// type MiningStatus struct {
// 	MinedBlocks    int    `json:"minedBlocks"`
// 	LastMinedBlock string `json:"lastMinedBlock"`
// 	IsMining       bool   `json:"isMining"`
// }

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

// // stopMining handles POST requests to stop the mining process
// func stopMining(w http.ResponseWriter, r *http.Request) {
// 	mu.Lock()
// 	defer mu.Unlock()

// 	if !miningStatus.IsMining {
// 		http.Error(w, "Mining is not currently active", http.StatusBadRequest)
// 		return
// 	}

// 	miningStatus.IsMining = false

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(miningStatus)
// }

// mine simulates the mining process by incrementing mined blocks periodically
// func mine() {
// 	for {
// 		mu.Lock()
// 		if !miningStatus.IsMining {
// 			mu.Unlock()
// 			return
// 		}

// 		// implement our own block generation
// 		blockHashes, err := client.Generate(1)
// 		if err != nil {
// 			log.Printf("Failed to generate block: %v", err)
// 			miningStatus.IsMining = false
// 			mu.Unlock()
// 			return
// 		}

// 		miningStatus.MinedBlocks++
// 		miningStatus.LastMinedBlock = blockHashes[0].String()
// 		mu.Unlock()

// 		// mining delay
// 		time.Sleep(10 * time.Second)
// 	}
// }

func main() {
	// set up btcd connection
	// connCfg := &rpcclient.ConnConfig{
	// 	Host:         "localhost:8334", // btcd RPC server
	// 	User:         "user",
	// 	Pass:         "password",
	// 	HTTPPostMode: true,
	// 	DisableTLS:   true,
	// }

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
		AllowCredentials: true,
	})

	// Use the CORS handler directly in the `ListenAndServe` function
	log.Println("Server starting on port 8080...")
	err = http.ListenAndServe(":8080", c.Handler(router)) // Pass the handler here
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
	


}