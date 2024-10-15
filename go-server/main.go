package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
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

// WalletResponse represents the response format for wallet generation
type WalletResponse struct {
	Message    string `json:"message"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

var (
	miningStatus MiningStatus
	mu           sync.Mutex
)

// Get current mining status
func getMiningStatus(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(miningStatus)
}

// Start mining process
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

// Stop mining process
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

// Mine blocks (simulated)
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

		time.Sleep(5 * time.Second) // Simulate mining time
	}
}

// Wallet generation handler
func generateWalletHandler(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("/Users/dianne/go/bin/btcwallet", "--create")

	output, err := cmd.CombinedOutput() // Captures both stdout and stderr
	if err != nil {
		log.Printf("Error generating wallet: %v, output: %s", err, string(output))
		http.Error(w, "Failed to generate wallet. Please check server logs.", http.StatusInternalServerError)
		return
	}

	// Example: Parsing output (you'll need to adjust based on actual btcwallet output)
	wallet := WalletResponse{
		Message:    "Wallet generated successfully",
		PublicKey:  extractPublicKeyFromOutput(),
		PrivateKey: extractPrivateKeyFromOutput(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wallet)

	log.Printf("Wallet generated successfully at %s", time.Now().Format(time.RFC3339))
}

func extractPublicKeyFromOutput() string {
    // Logic to extract public key
    return "public_key_example"
}

func extractPrivateKeyFromOutput() string {
    // Logic to extract private key
    return "private_key_example"
}
func main() {
	// Default mining status
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
	router.HandleFunc("/api/generate-wallet", generateWalletHandler).Methods("POST") // Wallet generation route

	// Configure CORS
	corsOptions := cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Frontend origin
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}

	handler := cors.New(corsOptions).Handler(router)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
