// application-layer/routes/authRoutes.go
package routes

import (
	"application-layer/controllers"

	"fmt"
	"github.com/gorilla/mux"
)

// RegisterAuthRoutes는 인증 관련 라우트를 등록합니다.
func RegisterAuthRoutes(router *mux.Router, authController *controllers.BtcController) {
	fmt.Println("Registering /api/auth routes...")

	authRouter := router.PathPrefix("/api/auth").Subrouter()
	authRouter.HandleFunc("/signup", authController.SignupHandler).Methods("POST")
	authRouter.HandleFunc("/login", authController.LoginHandler).Methods("POST") // Add this line
}

