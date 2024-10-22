package controllers

import (
	"encoding/json"
	"go-server/models"
	"go-server/services"
	"log"
	"net/http"
)

// Signup handles the signup request
func Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Println("Signup request received")

	// Call the wallet service to generate a key pair
	walletResponse, err := services.GenerateWallet()
	if err != nil {
		http.Error(w, "Error generating wallet", http.StatusInternalServerError)
		return
	}

	// Log success message
	log.Printf("User %s created successfully", walletResponse.PublicKey)

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
	if err := json.NewDecoder(r.Body).Decode(&challengeReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Find user by public key
	user, err := services.FindUserByPublicKey(challengeReq.PublicKey)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Get existing challenge or generate a new one
    challengeData, exists := services.GetChallenge(user.PublicKey)
    if exists {
        log.Printf("Returning existing challenge for publicKey: %s", user.PublicKey)
        response := map[string]string{"challenge": challengeData.Challenge}
        if err := json.NewEncoder(w).Encode(response); err != nil {
            http.Error(w, "Failed to encode response", http.StatusInternalServerError)
            return
        }
        return
    }

	// Generate a new challenge
	challenge, err := services.GenerateChallenge()
	if err != nil {
		http.Error(w, "Error generating challenge", http.StatusInternalServerError)
		return
	}

	// Store the challenge in memory
	if err := services.StoreChallenge(user.PublicKey, challenge); err != nil {
		http.Error(w, "Error storing challenge", http.StatusInternalServerError)
		return
	}

	// Return the challenge to the client
	response := map[string]string{"challenge": challenge}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// VerifyLogin handles the login process by verifying the signature
func VerifyChallenge(w http.ResponseWriter, r *http.Request) {
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

// AuthStatusResponse는 상태 확인 응답 구조체입니다.
type AuthStatusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// // AuthStatus는 JWT 토큰 상태를 검증하는 핸들러입니다.
// func AuthStatus(w http.ResponseWriter, r *http.Request) {
// 	// Authorization 헤더에서 토큰 추출
// 	token := r.Header.Get("Authorization")
// 	if token == "" {
// 		http.Error(w, "Missing token", http.StatusUnauthorized)
// 		return
// 	}

// 	// 토큰 검증
// 	valid, err := services.VerifyToken(token)
// 	if err != nil || !valid {
// 		http.Error(w, "Invalid token", http.StatusUnauthorized)
// 		return
// 	}

// 	// 토큰이 유효한 경우
// 	response := AuthStatusResponse{
// 		Status:  "valid",
// 		Message: "Token is valid",
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }
