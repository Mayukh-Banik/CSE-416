// application-layer/routes/btc_routes.go
package routes

import (
	"application-layer/controllers"

	"github.com/gorilla/mux"
)

// RegisterBtcRoutes 함수는 모든 BtcController의 핸들러를 라우터에 등록합니다.
func RegisterBtcRoutes(router *mux.Router, controller *controllers.BtcController) {
	btcRouter := router.PathPrefix("/api/btc").Subrouter() // /api/btc 경로를 기준으로 묶음

	// POST 요청 핸들러
	// btcRouter.HandleFunc("/init", controller.InitHandler).Methods("POST")
	btcRouter.HandleFunc("/login", controller.LoginHandler).Methods("POST")
	btcRouter.HandleFunc("/transaction", controller.TransactionHandler).Methods("POST")
	btcRouter.HandleFunc("/startbtcd", controller.StartBtcdHandler).Methods("POST")
	btcRouter.HandleFunc("/stopbtcd", controller.StopBtcdHandler).Methods("POST")
	btcRouter.HandleFunc("/startbtcwallet", controller.StartBtcwalletHandler).Methods("POST")
	btcRouter.HandleFunc("/stopbtcwallet", controller.StopBtcwalletHandler).Methods("POST")
	btcRouter.HandleFunc("/newaddress", controller.GetNewAddressHandler).Methods("POST")
	btcRouter.HandleFunc("/startmining", controller.StartMiningHandler).Methods("POST")
	btcRouter.HandleFunc("/stopmining", controller.StopMiningHandler).Methods("POST")

	// GET 요청 핸들러
	btcRouter.HandleFunc("/balance", controller.GetBalanceHandler).Methods("GET")
	btcRouter.HandleFunc("/getreceivedbyaddress", controller.GetReceivedByAddressHandler).Methods("GET")
	btcRouter.HandleFunc("/getblockcount", controller.GetBlockCountHandler).Methods("GET")
	btcRouter.HandleFunc("/listreceivedbyaddress", controller.ListReceivedByAddressHandler).Methods("GET")
	btcRouter.HandleFunc("/listunspent", controller.ListUnspentHandler).Methods("GET")
	btcRouter.HandleFunc("/getminingstatus", controller.GetMiningStatusHandler).Methods("GET")
}
