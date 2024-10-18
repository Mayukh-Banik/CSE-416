package main

import (
    "context"
    "log"
    "net/http"
	"os"	

    "go-server/routes"
    "go-server/utils"
	"github.com/rs/cors"
	"github.com/joho/godotenv"
)

func main() {
    // Initialize MongoDB connection
    // mongoURI := "mongodb://localhost:27017" // Update if necessary
	// MONGODB_URI := os.Getenv("MONGODB_URI")
    // utils.ConnectMongo(MONGODB_URI)

    // defer func() {
    //     if err := utils.MongoClient.Disconnect(context.TODO()); err != nil {
    //         log.Fatalf("Error disconnecting from MongoDB: %v", err)
    //     }
    // }()

	// Load the .env file to get environment variables
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    // Initialize MongoDB connection using the environment variable
    MONGODB_URI := os.Getenv("MONGODB_URI")
    if MONGODB_URI == "" {
        log.Fatalf("MONGODB_URI environment variable is not set")
    }

    utils.ConnectMongo(MONGODB_URI)

    defer func() {
        if err := utils.MongoClient.Disconnect(context.TODO()); err != nil {
            log.Fatalf("Error disconnecting from MongoDB: %v", err)
        }
    }()


    // Initialize routes
    router := http.NewServeMux()
    apiRouter := routes.InitRoutes()

    // "/api" 경로에 대해서 라우터 연결
    router.Handle("/api/", http.StripPrefix("/api", apiRouter))

    // Add CORS middleware
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000"}, // Frontend origin
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type", "Authorization"},
        AllowCredentials: true,
    })

    // Use the CORS handler directly in the `ListenAndServe` function
    log.Println("Server starting on port 8080...")
    err = http.ListenAndServe(":8080", c.Handler(router)) // Pass the handler here
    if err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}
