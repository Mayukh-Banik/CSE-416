package controllers

import (
	"application-layer/services"
	"encoding/json"
	"net/http"
)

// define Response struct
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// define BtcController struct
type BtcController struct {
	Service *services.BtcService
}

// initialize BtcController
func NewBtcController(service *services.BtcService) *BtcController {
	return &BtcController{Service: service}
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

func (bc *BtcController) InitHandler(w http.ResponseWriter, r *http.Request) {
	result := bc.Service.Init()
	resp := Response{
		Status:  "success",
		Message: result,
	}
	respondWithJSON(w, http.StatusOK, resp)
}

// LoginHandler 핸들러 함수
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

func (bc *BtcController) BtcwalletCreateHandler(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Passphrase string `json:"passphrase"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err := services.BtcwalletCreate(params.Passphrase)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := Response{
		Status:  "success",
		Message: "btcwallet process started in the background.",
	}
	respondWithJSON(w, http.StatusOK, resp)
}
