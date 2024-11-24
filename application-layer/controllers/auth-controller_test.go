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


// func TestHandleLogin(t *testing.T) {
// 	// Step 1: Initialize WalletService with test private key, UserService, and AuthController
// 	testPrivKey := generateTestKey()
// 	walletService := wallet.NewWalletService("user", "password", testPrivKey)
// 	// wallet.MockVerification = false // Disable mock verification
// 	userService := wallet.NewUserService()
// 	authController := NewAuthController(userService, walletService)

// 	// Step 2: Perform signup to create a test user
// 	user, privKey, err := userService.SignUp(*walletService)
// 	if err != nil {
// 		t.Fatalf("Signup failed: %v", err)
// 	}

// 	// Step 3: Generate a challenge for the user
// 	challengeRequest := map[string]string{
// 		"address": user.Address,
// 	}
// 	challengeRequestBody, _ := json.Marshal(challengeRequest)
// 	challengeReq := httptest.NewRequest(http.MethodPost, "/login/request", bytes.NewBuffer(challengeRequestBody))
// 	challengeReq.Header.Set("Content-Type", "application/json")
// 	challengeReqRecorder := httptest.NewRecorder()

// 	authController.HandleLoginRequest(challengeReqRecorder, challengeReq)

// 	// Step 4: Decode the challenge response
// 	challengeResp := challengeReqRecorder.Result()
// 	defer challengeResp.Body.Close()

// 	var challengeData map[string]string
// 	if err := json.NewDecoder(challengeResp.Body).Decode(&challengeData); err != nil {
// 		t.Fatalf("Failed to decode challenge response: %v", err)
// 	}

// 	challenge, exists := challengeData["challenge"]
// 	if !exists {
// 		t.Fatalf("Challenge not found in response")
// 	}

// 	// Step 5: Generate a valid signature for the challenge
// 	signature, err := walletService.SignChallenge(challenge)
// 	if err != nil {
// 		t.Fatalf("Failed to sign challenge: %v", err)
// 	}

// 	// Step 6: Create the login request with the valid signature
// 	loginRequest := map[string]string{
// 		"address":   user.Address,
// 		"signature": signature,
// 	}
// 	loginRequestBody, _ := json.Marshal(loginRequest)

// 	loginReq := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(loginRequestBody))
// 	loginReq.Header.Set("Content-Type", "application/json")
// 	loginRespRecorder := httptest.NewRecorder()

// 	// Step 7: Call HandleLogin
// 	authController.HandleLogin(loginRespRecorder, loginReq)

// 	// Step 8: Validate the response
// 	loginResp := loginRespRecorder.Result()
// 	defer loginResp.Body.Close()

// 	if loginResp.StatusCode != http.StatusOK {
// 		t.Fatalf("Expected status OK, got %v", loginResp.Status)
// 	}

// 	// Optionally, verify the response body
// 	body := new(bytes.Buffer)
// 	body.ReadFrom(loginResp.Body)
// 	if body.String() != "Login successful" {
// 		t.Fatalf("Unexpected login response: %s", body.String())
// 	}

// 	t.Logf("Login test passed for user: %s", user.UUID)
// }

// func TestHandleLogin(t *testing.T) {
// 	// Step 1: Initialize WalletService, UserService, and AuthController
// 	walletService := wallet.NewWalletService("user", "password")
// 	wallet.MockVerification = true // Enable mock verification for testing
// 	userService := wallet.NewUserService()
// 	authController := NewAuthController(userService, walletService)

// 	// Step 2: Perform signup to create a test user
// 	user, err := userService.SignUp(*walletService)
// 	if err != nil {
// 		t.Fatalf("Signup failed: %v", err)
// 	}

// 	// Step 3: Generate a challenge for the user
// 	challengeRequest := map[string]string{
// 		"address": user.Address,
// 	}
// 	challengeRequestBody, _ := json.Marshal(challengeRequest)
// 	challengeReq := httptest.NewRequest(http.MethodPost, "/login/request", bytes.NewBuffer(challengeRequestBody))
// 	challengeReq.Header.Set("Content-Type", "application/json")
// 	challengeReqRecorder := httptest.NewRecorder()

// 	authController.HandleLoginRequest(challengeReqRecorder, challengeReq)

// 	// Step 4: Decode the challenge response
// 	challengeResp := challengeReqRecorder.Result()
// 	defer challengeResp.Body.Close()

// 	var challengeData map[string]string
// 	if err := json.NewDecoder(challengeResp.Body).Decode(&challengeData); err != nil {
// 		t.Fatalf("Failed to decode challenge response: %v", err)
// 	}

// 	// Step 5: Mock a signature for the challenge
// 	signature := challengeData["challenge"]

// 	// Step 6: Create the login request
// 	loginRequest := map[string]string{
// 		"address":   user.Address,
// 		"signature": signature,
// 	}
// 	loginRequestBody, _ := json.Marshal(loginRequest)

// 	loginReq := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(loginRequestBody))
// 	loginReq.Header.Set("Content-Type", "application/json")
// 	loginRespRecorder := httptest.NewRecorder()

// 	// Step 7: Call HandleLogin
// 	authController.HandleLogin(loginRespRecorder, loginReq)

// 	// Step 8: Validate the response
// 	loginResp := loginRespRecorder.Result()
// 	defer loginResp.Body.Close()

// 	if loginResp.StatusCode != http.StatusOK {
// 		t.Fatalf("Expected status OK, got %v", loginResp.Status)
// 	}

// 	t.Logf("Login test passed for user: %s", user.UUID)
// }


// func generateTestKey() *ecdsa.PrivateKey {
// 	privKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
// 	return privKey
// }