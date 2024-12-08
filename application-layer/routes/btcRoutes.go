package routes

import (
	"application-layer/controllers"
	"net/http"
)

func BtcRoutes(mux *http.ServeMux, btcController *controllers.BtcController) {
	mux.HandleFunc("/api/wallet/start-btcd", btcController.StartBtcdHandler)
	mux.HandleFunc("/api/wallet/stop-btcd", btcController.StopBtcdHandler)
	mux.HandleFunc("/api/wallet/start-btcwallet", btcController.StartBtcwalletHandler)
	mux.HandleFunc("/api/wallet/stop-btcwallet", btcController.StopBtcwalletHandler)
}
