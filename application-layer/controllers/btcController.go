package controllers

import (
	"application-layer/services"
	"encoding/json"
	"net/http"
)

type BtcController struct {
	Service *services.BtcService
}

func NewBtcController(service *services.BtcService) *BtcController {
	return &BtcController{Service: service}
}

func (bc *BtcController) StartBtcdHandler(w http.ResponseWriter, r *http.Request) {
	status := bc.Service.StartBtcd()
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}

func (bc *BtcController) StopBtcdHandler(w http.ResponseWriter, r *http.Request) {
	status := bc.Service.StopBtcd()
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}

func (bc *BtcController) StartBtcwalletHandler(w http.ResponseWriter, r *http.Request) {
	status := bc.Service.StartBtcwallet()
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}

func (bc *BtcController) StopBtcwalletHandler(w http.ResponseWriter, r *http.Request) {
	status := bc.Service.StopBtcwallet()
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}
