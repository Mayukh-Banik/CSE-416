package routes

import (
	"go-server/controllers"

	"github.com/gorilla/mux"
)

func AuthRoutes(router *mux.Router) {
	router.HandleFunc("/api/auth/request-challenge", controllers.RequestChallenge).Methods("GET")
	router.HandleFunc("/api/auth/signup", controllers.Signup).Methods("POST")
	router.HandleFunc("/api/auth/login", controllers.Login).Methods("POST")
	router.HandleFunc("/api/auth/verify-challenge", controllers.VerifyChallenge).Methods("POST")
	router.HandleFunc("/api/auth/login-with-wallet", controllers.LoginWithWalletID).Methods("POST")
	// Check auth status (JWT validation)
	// router.HandleFunc("/api/auth/status", controllers.AuthStatus).Methods("GET")
}
