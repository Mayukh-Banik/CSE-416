// application-layer/routes/authRoutes.go
package routes

import (
	"application-layer/controllers"

	"github.com/gorilla/mux"
)

// RegisterAuthRoutes는 인증 관련 라우트를 등록합니다.
func RegisterAuthRoutes(router *mux.Router, authController *controllers.BtcController) {
	authRouter := router.PathPrefix("/api/auth").Subrouter()
	authRouter.HandleFunc("/signup", authController.SignupHandler).Methods("POST")
}
