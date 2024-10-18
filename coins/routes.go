// package routes

// import (
// 	"go-server/controllers"
// 	"github.com/gorilla/mux"
// )

// // InitRoutes initializes all the application routes
// func InitRoutes() *mux.Router {
// 	router := mux.NewRouter()

// 	// Auth routes (Signup/Login)
// 	router.HandleFunc("/signup", controllers.Signup).Methods("POST")
// 	//router.HandleFunc("/login", controllers.Login).Methods("POST")
// 	router.HandleFunc("/login", controllers.LoginWithWalletID).Methods("POST")
// 	router.HandleFunc("/login/challenge", controllers.RequestChallenge).Methods("POST")
// 	router.HandleFunc("/login/verify", controllers.VerifyLogin).Methods("POST")

// 	return router
// }
