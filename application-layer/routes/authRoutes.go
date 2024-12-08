package routes

import (
	"application-layer/controllers"
	"net/http"
)

func AuthRoutes(mux *http.ServeMux, authController *controllers.AuthController) {
	mux.HandleFunc("/api/auth/login", authController.LoginHandler)
	mux.HandleFunc("/api/auth/register", authController.RegisterHandler)
}
