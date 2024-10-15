package controllers

import (
	"encoding/json"
	"go-server/services"
	"log"
	"net/http"
	"go-server/models"
)

// Signup handles the signup request
func Signup(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    // Call the wallet service to generate a key pair
    walletResponse, err := services.GenerateWallet()
    if err != nil {
        http.Error(w, "Error generating wallet", http.StatusInternalServerError)
        return
    }

    // Log success message
    log.Printf("User %s created successfully", walletResponse.UserID)

    // Return the wallet response to the frontend
    json.NewEncoder(w).Encode(walletResponse)
}


// Login handles the login process (to be implemented)
func Login(w http.ResponseWriter, r *http.Request) {
	// Logic for login (to be implemented)
}

// RequestChallenge handles the request for a login challenge
func RequestChallenge(w http.ResponseWriter, r *http.Request) {
    var challengeReq models.ChallengeRequest
    json.NewDecoder(r.Body).Decode(&challengeReq)

    challenge, err := services.GenerateChallenge(challengeReq.PublicKey)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"challenge": challenge})
}

// VerifyLogin handles the login process by verifying the signature
func VerifyLogin(w http.ResponseWriter, r *http.Request) {
    var challengeResponse models.ChallengeResponse
    json.NewDecoder(r.Body).Decode(&challengeResponse)

    verified, err := services.VerifySignature(challengeResponse.UserID, challengeResponse.Signature)
    if err != nil || !verified {
        http.Error(w, `{"error": "Invalid signature. Login failed."}`, http.StatusUnauthorized)
        return
    }

    // Successfully authenticated
    json.NewEncoder(w).Encode(map[string]string{"status": "login successful"})
}

// LoginWithWalletID handles the login using only the walletId (public key)
func LoginWithWalletID(w http.ResponseWriter, r *http.Request) {
	var req struct {
		WalletID string `json:"walletId"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Call the service to login with walletId (public key)
	success, err := services.LoginWithWalletID(req.WalletID)
	if err != nil {
		http.Error(w, "Login failed: "+err.Error(), http.StatusUnauthorized)
		return
	}

	if success {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
	} else {
		http.Error(w, "Login failed", http.StatusUnauthorized)
	}
}