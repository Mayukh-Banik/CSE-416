package controllers

import (
	"encoding/json"
	"go-server/services"
	"log"
	"net/http"
)

// Signup handles the signup request, generates keys, and returns the public and private keys
func Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Call the wallet service to generate a key pair
	walletResponse, err := services.GenerateWallet()
	if err != nil {
		http.Error(w, "Error generating wallet", http.StatusInternalServerError)
		return
	}

	// Log a success message
	log.Printf("Wallet successfully generated with public key: %s", walletResponse.PublicKey)

	// Return both public and private keys to the frontend
	json.NewEncoder(w).Encode(walletResponse)
}


// Login handles the login process (future implementation)
func Login(w http.ResponseWriter, r *http.Request) {
	// Logic for login (to be implemented)
}
