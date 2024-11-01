package controllers

import (
	"encoding/json"
	"go-server/services"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// JWT generation function
func generateJWT(publicKey string) (string, error) {
	claims := jwt.MapClaims{
		"publicKey": publicKey,
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // Expiration time set to 24 hours later
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

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
	log.Printf("RequestChallenge@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	// Generate a new challenge
	challenge, err := services.GenerateChallenge()
	if err != nil {
		http.Error(w, "Failed to generate challenge", http.StatusInternalServerError)
		return
	}

	log.Printf("[challenge]: %s", challenge)

	// The challenge can be stored in a session or cookie. Here, it's simply stored in memory.
	services.StoreChallenge(challenge)

	// Send the challenge to the client
	json.NewEncoder(w).Encode(map[string]string{"challenge": challenge})
}

// VerifyChallenge handles the login process by verifying the signature
func VerifyChallenge(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PublicKey string `json:"public_key"`
		Signature string `json:"signature"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Printf("[Received PublicKey]: \n %s", req.PublicKey)
	log.Printf("[Received Signature]: \n %s", req.Signature)

	// Signature verification
	verified, err := services.VerifySignature(req.PublicKey, req.Signature)
	if err != nil || !verified {
		log.Printf("Error verifying signature: %v", err)
		http.Error(w, "Invalid signature. Login failed.", http.StatusUnauthorized)
		return
	}

	// On successful signature verification, generate a JWT token
	token, err := generateJWT(req.PublicKey)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Set the token as an HttpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	})

	// Success response
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

// AuthStatusResponse is the structure for the status check response
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
