// main.go
package main

import (
	"application-layer/controllers"
	"application-layer/routes"
	"application-layer/wallet"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	rpcUser := os.Getenv("RPC_USER")
	rpcPass := os.Getenv("RPC_PASS")
	if rpcUser == "" || rpcPass == "" {
		log.Fatal("RPC_USER and RPC_PASS environment variables are required")
	}

	// WalletService 초기화
	walletService := wallet.NewWalletService(rpcUser, rpcPass)

	// UserService 초기화
	userService := wallet.NewUserService()

	// 컨트롤러 초기화
	walletController := controllers.NewWalletController(walletService)
	authController := controllers.NewAuthController(userService, walletService)

	// 라우트 설정
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, authController, walletController)

	// CORS 설정
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"POST"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	// HTTP 서버 실행
	handler := c.Handler(mux)
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
