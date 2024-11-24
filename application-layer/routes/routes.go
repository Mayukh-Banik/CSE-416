package routes

import (
	"application-layer/controllers"
	"net/http"
)

// RegisterRoutes registers all HTTP routes
func RegisterRoutes(authController *controllers.AuthController) {
	http.HandleFunc("/signup", authController.HandleSignUp)
}

