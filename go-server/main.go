package main

import (
    "context"
    "log"
    "net/http"
	

    "go-server/routes"
    "go-server/utils"
	"github.com/rs/cors"
)

func main() {
    // Initialize MongoDB connection
    mongoURI := "mongodb://localhost:27017" // Update if necessary
    utils.ConnectMongo(mongoURI)

    defer func() {
        if err := utils.MongoClient.Disconnect(context.TODO()); err != nil {
            log.Fatalf("Error disconnecting from MongoDB: %v", err)
        }
    }()

    // Initialize routes
    router := routes.InitRoutes()

    // Add CORS middleware
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000"}, // Frontend origin
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type", "Authorization"},
        AllowCredentials: true,
    })

    // Use the CORS handler directly in the `ListenAndServe` function
    log.Println("Server starting on port 8080...")
    err := http.ListenAndServe(":8080", c.Handler(router)) // Pass the handler here
    if err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}
