package controllers

import (
	"application-layer/services"
	"encoding/json"
	"net/http"
)

type AuthController struct {
	Service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{Service: service}
}

func (ac *AuthController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// 예제 요청 데이터 처리
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if ac.Service.Login(creds.Username, creds.Password) {
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func (ac *AuthController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// 회원가입 요청 처리
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if ac.Service.Register(creds.Username, creds.Password) {
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	} else {
		http.Error(w, "Registration failed", http.StatusInternalServerError)
	}
}

// import (
// 	"application-layer/models"
// 	"application-layer/wallet"
// 	"encoding/json"
// 	"net/http"
// 	"os"
// )

// type AuthController struct {
// 	UserService   *wallet.UserService
// 	WalletService *wallet.WalletService
// }

// // NewAuthController initializes the AuthController.
// func NewAuthController(userService *wallet.UserService, walletService *wallet.WalletService) *AuthController {
// 	return &AuthController{
// 		UserService:   userService,
// 		WalletService: walletService,
// 	}
// }

// // HandleSignUp handles the signup process
// func (ac *AuthController) HandleSignUp(w http.ResponseWriter, r *http.Request) {
// 	// Retrieve the passphrase for wallet unlocking
// 	passphrase := os.Getenv("WALLET_PASSPHRASE")

// 	// Perform signup
// 	user, privateKey, err := ac.UserService.SignUp(*ac.WalletService, passphrase)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Respond with the user data and private key
// 	response := struct {
// 		User       *models.User `json:"user"`
// 		PrivateKey string       `json:"private_key"`
// 	}{
// 		User:       user,
// 		PrivateKey: privateKey,
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }

// // HandleLoginRequest handles generating a challenge for login.
// func (ac *AuthController) HandleLoginRequest(w http.ResponseWriter, r *http.Request) {
// 	type LoginRequest struct {
// 		Address string `json:"address"`
// 	}

// 	var req LoginRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Invalid request", http.StatusBadRequest)
// 		return
// 	}

// 	challenge, err := ac.UserService.GenerateChallenge(req.Address)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(challenge)
// }

// // HandleLogin handles verifying the signed challenge for login.
// func (ac *AuthController) HandleLogin(w http.ResponseWriter, r *http.Request) {
// 	type LoginVerification struct {
// 		Address   string `json:"address"`
// 		Signature string `json:"signature"`
// 	}

// 	var req LoginVerification
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Invalid request", http.StatusBadRequest)
// 		return
// 	}

// 	// Retrieve the challenge
// 	challenge, err := ac.UserService.GetChallenge(req.Address)
// 	if err != nil {
// 		http.Error(w, "Challenge expired or not found", http.StatusUnauthorized)
// 		return
// 	}

// 	// Verify the signature using btcctl (using btcctl's command-line utility)
// 	valid, err := ac.WalletService.VerifySignature(req.Address, challenge.Challenge, req.Signature)
// 	if err != nil || !valid {
// 		http.Error(w, "Invalid signature", http.StatusUnauthorized)
// 		return
// 	}

// 	// Login successful, remove the challenge
// 	ac.UserService.RemoveChallenge(req.Address)
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("Login successful"))
// }
