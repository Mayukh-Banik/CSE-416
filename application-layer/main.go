// main.go
package main

import (
	"application-layer/controllers"
	dht_kad "application-layer/dht"
	"application-layer/download"
	"application-layer/files"
	"application-layer/routes"
	"application-layer/services"
	"application-layer/websocket"
	"fmt"
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

	// 서비스 초기화
	btcService := services.NewBtcService()

	// initialization
	// walletService, err := wallet.NewWalletService(rpcUser, rpcPass)
	// if err != nil {
	// 	log.Fatalf("Failed to initialize WalletService: %v", err)
	// }
	// userService := wallet.NewUserService()
	// walletController := controllers.NewWalletController(walletService)
	// authController := controllers.NewAuthController(userService, walletService)

	fmt.Println("Main server started")

	fileRouter := files.InitFileRoutes()
	downloadRouter := download.InitDownloadRoutes()
	btcController := controllers.NewBtcController(btcService)

	// 라우터 초기화
	mux := http.NewServeMux()
	routes.BtcRoutes(mux, btcController)

	go dht_kad.StartDHTService()

	// CORS handler
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},        // Frontend's origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"}, // Allowed HTTP methods
		AllowedHeaders:   []string{"Content-Type", "Hash"},         // Allowed headers
		AllowCredentials: true,                                     // Allow credentials (cookies, auth headers)
	})

	// Combine both routers on the same port
	http.Handle("/files/", c.Handler(fileRouter))        // File routes under /files
	http.Handle("/download/", c.Handler(downloadRouter)) // Download routes under /download
	http.Handle("/ws", http.HandlerFunc(websocket.WsHandler))

	port := ":8080"
	handler := c.Handler(mux)
	fmt.Printf("Starting server for file routes and DHT on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, handler))
}
