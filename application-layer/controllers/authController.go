package controllers

import (
	"application-layer/wallet"
	"application-layer/models"
	"encoding/json"
	"net/http"
)

type AuthController struct {
	UserService   *wallet.UserService
	WalletService *wallet.WalletService
}

// HandleSignUp handles the signup process
func (ac *AuthController) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	// Retrieve the passphrase for wallet unlocking
	passphrase := "CSE416" // Hardcoded for testing; replace with environment variable for production

	// Perform signup
	user, privateKey, err := ac.UserService.SignUp(*ac.WalletService, passphrase)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the user data and private key
	response := struct {
		User       *models.User `json:"user"`
		PrivateKey string       `json:"private_key"`
	}{
		User:       user,
		PrivateKey: privateKey,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}