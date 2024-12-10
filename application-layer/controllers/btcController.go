// application-layer/controllers/btcController.go
package controllers

import (
	"application-layer/services"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type BtcController struct {
	Service *services.BtcService
}

// Response is a generic response structure
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// initialize BtcController
func NewBtcController(service *services.BtcService) *BtcController {
	return &BtcController{Service: service}
}

type SignupRequest struct {
	Passphrase string `json:"passphrase"`
}

type SignupResponse struct {
	Address    string `json:"address,omitempty"`
	PrivateKey string `json:"private_key,omitempty"`
	Message    string `json:"message"`
}

// Helper function: respondWithJSON
func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "JSON Encoding Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

// Helper function: respondWithError
func respondWithError(w http.ResponseWriter, status int, message string) {
	resp := Response{
		Status:  "error",
		Message: message,
	}
	respondWithJSON(w, status, resp)
}

// SignupHandler processes signup requests.
func (bc *BtcController) SignupHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SignupHandler called")

	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Check if the passphrase is empty
	if req.Passphrase == "" {
		respondWithError(w, http.StatusBadRequest, "Passphrase is required")
		return
	}

	// Create channels to receive results
	resultCh := make(chan SignupResponse)
	errorCh := make(chan error)

	// Execute wallet creation in a separate goroutine
	go func() {
		newAddress, err := bc.Service.CreateWallet(req.Passphrase)
		if err != nil {
			errorCh <- err
			return
		}
		// Replace with actual private key generation logic
		privateKey := "generated-private-key"
		response := SignupResponse{
			Address:    newAddress,
			PrivateKey: privateKey,
			Message:    "Wallet successfully created.",
		}
		resultCh <- response
	}()

	// Set a 30-second timeout
	select {
	case res := <-resultCh:
		respondWithJSON(w, http.StatusOK, res)
	case err := <-errorCh:
		respondWithError(w, http.StatusInternalServerError, err.Error())
	case <-time.After(30 * time.Second):
		// Stop processes on timeout
		log.Println("SignupHandler: Request timed out.")
		stopWalletMsg := bc.Service.StopBtcwallet()
		stopBtcdMsg := bc.Service.StopBtcd()
		errorMessage := fmt.Sprintf("Request timed out. %s %s", stopWalletMsg, stopBtcdMsg)
		respondWithError(w, http.StatusGatewayTimeout, errorMessage)
	}
}

// placeholder for login handler
func (bc *BtcController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var params struct {
		WalletAddress string `json:"walletAddress"`
		Passphrase    string `json:"passphrase"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	result, err := bc.Service.Login(params.WalletAddress, params.Passphrase)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := Response{
		Status:  "success",
		Message: result,
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func (bc *BtcController) TransactionHandler(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Passphrase string  `json:"passphrase"`
		Txid       string  `json:"txid"`
		Dst        string  `json:"dst"`
		Amount     float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	txIdResult, err := bc.Service.Transaction(params.Passphrase, params.Txid, params.Dst, params.Amount)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := Response{
		Status: "success",
		Data:   map[string]string{"txid": txIdResult},
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func (bc *BtcController) StartBtcdHandler(w http.ResponseWriter, r *http.Request) {
	var params struct {
		WalletAddress string `json:"walletAddress,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var result string
	if params.WalletAddress != "" {
		result = bc.Service.StartBtcd(params.WalletAddress)
	} else {
		result = bc.Service.StartBtcd()
	}

	resp := Response{
		Status:  "success",
		Message: result,
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func (bc *BtcController) StopBtcdHandler(w http.ResponseWriter, r *http.Request) {
	result := bc.Service.StopBtcd()
	resp := Response{
		Status:  "success",
		Message: result,
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func (bc *BtcController) StartBtcwalletHandler(w http.ResponseWriter, r *http.Request) {
	result := bc.Service.StartBtcwallet()
	resp := Response{
		Status:  "success",
		Message: result,
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func (bc *BtcController) StopBtcwalletHandler(w http.ResponseWriter, r *http.Request) {
	result := bc.Service.StopBtcwallet()
	resp := Response{
		Status:  "success",
		Message: result,
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func (bc *BtcController) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	balance, err := bc.Service.GetBalance()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := Response{
		Status: "success",
		Data:   map[string]string{"balance": balance},
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func (bc *BtcController) GetNewAddressHandler(w http.ResponseWriter, r *http.Request) {
	newAddress, err := bc.Service.GetNewAddress()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := Response{
		Status: "success",
		Data:   map[string]string{"newAddress": newAddress},
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func (bc *BtcController) GetReceivedByAddressHandler(w http.ResponseWriter, r *http.Request) {
	walletAddress := r.URL.Query().Get("walletAddress")
	if walletAddress == "" {
		respondWithError(w, http.StatusBadRequest, "walletAddress is required")
		return
	}

	receivedAmount, err := bc.Service.GetReceivedByAddress(walletAddress)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := Response{
		Status: "success",
		Data:   map[string]string{"receivedAmount": receivedAmount},
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func (bc *BtcController) GetBlockCountHandler(w http.ResponseWriter, r *http.Request) {
	blockCount, err := bc.Service.GetBlockCount()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := Response{
		Status: "success",
		Data:   map[string]string{"blockCount": blockCount},
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func (bc *BtcController) ListReceivedByAddressHandler(w http.ResponseWriter, r *http.Request) {
	addresses, err := bc.Service.ListReceivedByAddress()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := Response{
		Status: "success",
		Data:   addresses,
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func (bc *BtcController) ListUnspentHandler(w http.ResponseWriter, r *http.Request) {
	utxos, err := bc.Service.ListUnspent()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := Response{
		Status: "success",
		Data:   utxos,
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func (bc *BtcController) GetCurrentAddressHandler(w http.ResponseWriter, r *http.Request) {
	var tempFilePath string
	if os.Getenv("OS") == "Windows_NT" {
		tempFilePath = filepath.Join(os.Getenv("TEMP"), "btc_temp.json")
	} else {
		tempFilePath = "/tmp/btcd_temp.json"
	}
	content, err := ioutil.ReadFile(tempFilePath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to read temp file: %v", err))
		return
	}

	var data map[string]string
	if err := json.Unmarshal(content, &data); err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to parse JSON: %v", err))
		return
	}

	currentAddress, exists := data["miningaddr"]
	if !exists {
		respondWithError(w, http.StatusNotFound, tempFilePath+" does not contain the current address.")
		return
	}

	resp := Response{
		Status: "success",
		Data:   currentAddress,
	}

	// Debugging
	fmt.Println("Current address: ", currentAddress)
	respondWithJSON(w, http.StatusOK, resp)
}

func (bc *BtcController) GetMiningStatusHandler(w http.ResponseWriter, r *http.Request) {
	status, err := bc.Service.GetMiningStatus()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := Response{
		Status: "success",
		Data:   map[string]bool{"mining": status},
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func (bc *BtcController) StartMiningHandler(w http.ResponseWriter, r *http.Request) {
	var params struct {
		NumBlock int `json:"numBlock"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if params.NumBlock <= 0 {
		respondWithError(w, http.StatusBadRequest, "numBlock must be greater than 0")
		return
	}

	result := bc.Service.StartMining(params.NumBlock)

	resp := Response{
		Status:  "success",
		Message: result,
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func (bc *BtcController) StopMiningHandler(w http.ResponseWriter, r *http.Request) {
	result := bc.Service.StopMining()

	resp := Response{
		Status:  "success",
		Message: result,
	}
	respondWithJSON(w, http.StatusOK, resp)
}

// GetMiningDashboardHandler returns the mining dashboard data
func (bc *BtcController) GetMiningDashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Get balance
	balance, err := bc.Service.GetBalance()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get balance: %v", err))
		return
	}

	// Get mining info
	miningInfoRaw, err := bc.Service.GetMiningInfo()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get mining info: %v", err))
		return
	}

	// Parse mining info JSON
	var miningInfo map[string]interface{}
	if err := json.Unmarshal([]byte(miningInfoRaw), &miningInfo); err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to parse mining info JSON: %v", err))
		return
	}

	// Combine results
	dashboard := map[string]interface{}{
		"balance":    balance,
		"miningInfo": miningInfo,
	}

	// Respond with combined data
	resp := Response{
		Status: "success",
		Data:   dashboard,
	}
	respondWithJSON(w, http.StatusOK, resp)
}

// InitHandler handles the initialization process
func (bc *BtcController) InitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	fmt.Println("InitHandler called")

	result := bc.Service.Init()

	resp := Response{
		Status:  "success",
		Message: result,
	}

	respondWithJSON(w, http.StatusOK, resp)
}
