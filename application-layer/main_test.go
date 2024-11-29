package main_test

import (
	"application-layer/controllers"
	"application-layer/models"
	"application-layer/wallet"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestHandleSignUp(t *testing.T) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	// Step 1: Initialize WalletService, UserService, and AuthController
	passphrase := os.Getenv("WALLET_PASSPHRASE")
	if passphrase == "" {
		t.Fatalf("Environment variable WALLET_PASSPHRASE is not set")
	}
	t.Logf("Loaded WALLET_PASSPHRASE: %s", passphrase)

	walletService := wallet.NewWalletService("user", "password")
	userService := wallet.NewUserService()
	authController := &controllers.AuthController{
		UserService:   userService,
		WalletService: walletService,
	}

	// Step 2: Create the signup request
	reqBody := map[string]string{
		"passphrase": passphrase,
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	respRecorder := httptest.NewRecorder()

	// Step 3: Call HandleSignUp
	authController.HandleSignUp(respRecorder, req)

	// Step 4: Validate the response
	resp := respRecorder.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Expected status OK, got %v. Additionally, failed to read response body: %v", resp.Status, err)
		}
		t.Fatalf("Expected status OK, got %v. Response body: %s", resp.Status, string(bodyBytes))
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
	authController := controllers.NewAuthController(userService, walletService)

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
	authController := controllers.NewAuthController(userService, walletService)

	// Step 2: Create a test wallet address and generate a challenge
	address, _, _, _, err := walletService.GenerateNewAddressWithPubKeyAndPrivKey(os.Getenv("WALLET_PASSPHRASE"))
	if err != nil {
		t.Fatalf("Failed to generate new address with private key: %v", err)
	}

	// Step 3: Generate a challenge for the address
	challenge, err := userService.GenerateChallenge(address)
	if err != nil {
		t.Fatalf("Failed to generate challenge: %v", err)
	}

	// Step 4: Sign the challenge using btcctl (using the SignMessage function)
	signature, err := walletService.SignMessage(address, challenge.Challenge, os.Getenv("WALLET_PASSPHRASE"))
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

func TestSignUp(t *testing.T) {
	// Step 1: Initialize services
	walletService := wallet.NewWalletService("user", "password")
	userService := wallet.NewUserService()
	passphrase := os.Getenv("WALLET_PASSPHRASE")

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

func TestGenerateNewAddress(t *testing.T) {
	ws := wallet.NewWalletService("user", "password") // No parameters needed as it loads from .env
	address, pubKey, err := ws.GenerateNewAddressWithPubKey()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if address == "" {
		t.Fatalf("Expected a new address, got empty string")
	}
	t.Logf("Generated Address: %s\nPublic Key: %s\n", address, pubKey)
}

func TestGenerateNewAddressWithPubKeyAndPrivKey(t *testing.T) {
	ws := wallet.NewWalletService("user", "password") // No parameters needed as it loads from .env
	passphrase := os.Getenv("WALLET_PASSPHRASE")      // Replace with a valid passphrase
	address, pubKey, privateKey, balance, err := ws.GenerateNewAddressWithPubKeyAndPrivKey(passphrase)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if address == "" {
		t.Fatalf("Expected a new address, got empty string")
	}
	if pubKey == "" {
		t.Fatalf("Expected a public key, got empty string")
	}
	if privateKey == "" {
		t.Fatalf("Expected a private key, got empty string")
	}
	if balance < 0 {
		t.Fatalf("Expected a non-negative balance, got %f", balance)
	}
	t.Logf("Generated Address: %s\nPublic Key: %s\nPrivate Key: %s\nBalance: %f\n", address, pubKey, privateKey, balance)
}
