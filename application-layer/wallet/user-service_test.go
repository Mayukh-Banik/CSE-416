package wallet

import (
	"testing"
)

func TestSignUp(t *testing.T) {
	// Step 1: Initialize services
	walletService := NewWalletService("user", "password")
	userService := NewUserService()
	passphrase := "CSE416" // Hardcoded passphrase for testing

	// Step 2: Perform signup
	user, privateKey, err := userService.SignUp(*walletService, passphrase)
	if err != nil {
		t.Fatalf("Signup failed: %v", err)
	}

	// Step 3: Validate results
	if user.Address == "" {
		t.Errorf("Expected a valid wallet address, got empty string")
	}
	if user.PublicKey == "" {
		t.Errorf("Expected a valid public key, got empty string")
	}
	if privateKey == "" {
		t.Errorf("Expected a valid private key, got empty string")
	}
	if user.UUID == "" {
		t.Errorf("Expected a valid UUID, got empty string")
	}
	if user.Balance < 0 {
		t.Errorf("Expected a non-negative balance, got %v", user.Balance)
	}

	// Step 4: Log the created user for debugging
	t.Logf("User created successfully: %+v", user)
	t.Logf("Private Key: %s", privateKey)
}