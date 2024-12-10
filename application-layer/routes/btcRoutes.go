// application-layer/routes/btc_routes.go
package routes

import (
	"application-layer/controllers"

	"github.com/gorilla/mux"
)

func RegisterBtcRoutes(router *mux.Router, controller *controllers.BtcController) {
	btcRouter := router.PathPrefix("/api/btc").Subrouter()

	// btcRouter.HandleFunc("/init", controller.InitHandler).Methods("POST")
	btcRouter.HandleFunc("/login", controller.LoginHandler).Methods("POST")
	btcRouter.HandleFunc("/transaction", controller.TransactionHandler).Methods("POST")
	btcRouter.HandleFunc("/newaddress", controller.GetNewAddressHandler).Methods("POST")
	btcRouter.HandleFunc("/startmining", controller.StartMiningHandler).Methods("POST")
	btcRouter.HandleFunc("/stopmining", controller.StopMiningHandler).Methods("POST")

	btcRouter.HandleFunc("/balance", controller.GetBalanceHandler).Methods("GET")
	btcRouter.HandleFunc("/getreceivedbyaddress", controller.GetReceivedByAddressHandler).Methods("GET")
	btcRouter.HandleFunc("/getblockcount", controller.GetBlockCountHandler).Methods("GET")
	btcRouter.HandleFunc("/listreceivedbyaddress", controller.ListReceivedByAddressHandler).Methods("GET")
	btcRouter.HandleFunc("/listunspent", controller.ListUnspentHandler).Methods("GET")
	btcRouter.HandleFunc("/getminingstatus", controller.GetMiningStatusHandler).Methods("GET")
	btcRouter.HandleFunc("/currentaddress", controller.GetCurrentAddressHandler).Methods("GET")
	btcRouter.HandleFunc("/miningdashboard", controller.GetMiningDashboardHandler).Methods("GET")

}
