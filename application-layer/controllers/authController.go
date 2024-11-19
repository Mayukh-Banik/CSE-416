package controllers

import (
	"application-layer/wallet"
	"encoding/json"
	"net/http"
)

type AuthController struct {
	UserService   *wallet.UserService
	WalletService *wallet.WalletService
}

// HandleSignUp handles the signup process
func (ac *AuthController) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	// Perform signup
	user, err := ac.UserService.SignUp(*ac.WalletService)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the user data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
