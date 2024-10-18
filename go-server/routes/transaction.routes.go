package routes

import (
	// "go-server/controllers"
	"github.com/gorilla/mux"
)

// TransactionRoutes initializes transaction-related routes.
func TransactionRoutes() *mux.Router {
	router := mux.NewRouter()

	// 거래 내역 조회 경로
	// router.HandleFunc("/api/transactions/", controllers.GetTransactionHistory).Methods("GET")

	return router
}
