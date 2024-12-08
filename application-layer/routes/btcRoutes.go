// application-layer/routes/btc_routes.go
package routes

import (
	"application-layer/controllers"

	"github.com/gorilla/mux"
)

// RegisterBtcRoutes 함수는 모든 BtcController의 핸들러를 라우터에 등록합니다.
func RegisterBtcRoutes(router *mux.Router, controller *controllers.BtcController) {
	router.HandleFunc("/init", controller.InitHandler).Methods("POST")
	router.HandleFunc("/login", controller.LoginHandler).Methods("POST")
	router.HandleFunc("/transaction", controller.TransactionHandler).Methods("POST")
	router.HandleFunc("/startbtcd", controller.StartBtcdHandler).Methods("POST")
	router.HandleFunc("/stopbtcd", controller.StopBtcdHandler).Methods("POST")
	router.HandleFunc("/startbtcwallet", controller.StartBtcwalletHandler).Methods("POST")
	router.HandleFunc("/stopbtcwallet", controller.StopBtcwalletHandler).Methods("POST")
	router.HandleFunc("/balance", controller.GetBalanceHandler).Methods("GET")
	router.HandleFunc("/newaddress", controller.GetNewAddressHandler).Methods("POST")
	router.HandleFunc("/getreceivedbyaddress", controller.GetReceivedByAddressHandler).Methods("GET")
	router.HandleFunc("/getblockcount", controller.GetBlockCountHandler).Methods("GET")
	router.HandleFunc("/listreceivedbyaddress", controller.ListReceivedByAddressHandler).Methods("GET")
	router.HandleFunc("/listunspent", controller.ListUnspentHandler).Methods("GET")
	router.HandleFunc("/getminingstatus", controller.GetMiningStatusHandler).Methods("GET")
	router.HandleFunc("/startmining", controller.StartMiningHandler).Methods("POST")
	router.HandleFunc("/stopmining", controller.StopMiningHandler).Methods("POST")
	router.HandleFunc("/btcwalletcreate", controller.BtcwalletCreateHandler).Methods("POST")
}
