package main

import (
	"log"
	"net/http"
	"go-server/routes"
)

func main() {
	router := routes.InitRoutes() // Initialize routes

	log.Println("Server starting on port 8080...")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
