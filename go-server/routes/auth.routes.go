package routes

import (
	"go-server/controllers"
	"github.com/gorilla/mux"
)

// AuthRoutes initializes authentication-related routes.
func AuthRoutes() *mux.Router {
	router := mux.NewRouter()

	// 회원가입 경로
	router.HandleFunc("/users/signup", controllers.Signup).Methods("POST")

	// 로그인 (지갑 ID를 사용한 로그인)
	router.HandleFunc("/users/login", controllers.LoginWithWalletID).Methods("POST")

	// 로그인 도전 과제 요청
	router.HandleFunc("/users/login/challenge", controllers.RequestChallenge).Methods("POST")

	// 로그인 서명 검증
	router.HandleFunc("/users/login/verify", controllers.VerifyLogin).Methods("POST")

	return router
}
