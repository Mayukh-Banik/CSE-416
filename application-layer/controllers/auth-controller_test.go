package controllers

import (
	"application-layer/models"
	"application-layer/wallet"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleSignUp(t *testing.T) {
	// Step 1: Initialize WalletService, UserService, and AuthController
	walletService := wallet.NewWalletService("user", "password")
	userService := wallet.NewUserService()
	authController := &AuthController{
		UserService:   userService,
		WalletService: walletService,
	}

	// Step 2: Create the signup request
	reqBody := map[string]string{
		"passphrase": "CSE416",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	respRecorder := httptest.NewRecorder()

	// Step 3: Call HandleSignUp
	authController.HandleSignUp(respRecorder, req)

	// Step 4: Validate the response
	resp := respRecorder.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", resp.Status)
	}

	var responseData struct {
		User       *models.User `json:"user"`
		PrivateKey string       `json:"private_key"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if responseData.User == nil || responseData.User.Address == "" || responseData.PrivateKey == "" {
		t.Fatalf("Invalid user data or private key in response: %+v", responseData)
	}

	t.Logf("SignUp test passed. User: %+v, PrivateKey: %s", responseData.User, responseData.PrivateKey)
}

// TestHandleLoginRequest tests the generation of a login challenge for a given wallet address.
func TestHandleLoginRequest(t *testing.T) {
	// Step 1: Initialize services
	walletService := wallet.NewWalletService("user", "password")
	userService := wallet.NewUserService()
	authController := NewAuthController(userService, walletService)

	// Step 2: Create a test wallet address
	address, _, _, _, err := walletService.GenerateNewAddressWithPubKeyAndPrivKey("CSE416")
	if err != nil {
		t.Fatalf("Failed to generate new address: %v", err)
	}

	// Step 3: Create the login request
	loginRequest := map[string]string{
		"address": address,
	}
	loginRequestBody, _ := json.Marshal(loginRequest)

	// Step 4: Create an HTTP request to generate a challenge
	req := httptest.NewRequest(http.MethodPost, "/login/request", bytes.NewBuffer(loginRequestBody))
	req.Header.Set("Content-Type", "application/json")
	respRecorder := httptest.NewRecorder()

	// Step 5: Call HandleLoginRequest to generate the challenge
	authController.HandleLoginRequest(respRecorder, req)

	// Step 6: Decode the challenge response
	resp := respRecorder.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", resp.Status)
	}

	var challengeData map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&challengeData); err != nil {
		t.Fatalf("Failed to decode challenge response: %v", err)
	}

	if challengeData["challenge"] == "" {
		t.Fatalf("Expected a challenge string, got empty string")
	}

	t.Logf("Challenge generated successfully for address %s: %s", address, challengeData["challenge"])
}

// TestHandleLogin tests the verification of a signed challenge.
func TestHandleLogin(t *testing.T) {
    // Step 1: Initialize services
    walletService := wallet.NewWalletService("user", "password")
    userService := wallet.NewUserService()
    authController := NewAuthController(userService, walletService)

    // Step 2: Create a test wallet address and generate a challenge
    address, _, _, _, err := walletService.GenerateNewAddressWithPubKeyAndPrivKey("CSE416")
    if err != nil {
        t.Fatalf("Failed to generate new address with private key: %v", err)
    }

    // Step 3: Generate a challenge for the address
    challenge, err := userService.GenerateChallenge(address)
    if err != nil {
        t.Fatalf("Failed to generate challenge: %v", err)
    }

    // Step 4: Sign the challenge using btcctl (using the SignMessage function)
    signature, err := walletService.SignMessage(address, challenge.Challenge, "CSE416")
    if err != nil {
        t.Fatalf("Failed to sign challenge: %v", err)
    }

    // Step 5: Create the login verification request
    loginVerification := map[string]string{
        "address":   address,
        "signature": signature,
    }
    loginVerificationBody, _ := json.Marshal(loginVerification)

    // Step 6: Create an HTTP request for login verification
    req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(loginVerificationBody))
    req.Header.Set("Content-Type", "application/json")
    respRecorder := httptest.NewRecorder()

    // Step 7: Call HandleLogin to verify the signed challenge
    authController.HandleLogin(respRecorder, req)

    // Step 8: Validate the response
    resp := respRecorder.Result()
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Fatalf("Expected status OK, got %v", resp.Status)
    }

    t.Logf("Login test passed for address %s with signature %s", address, signature)
}