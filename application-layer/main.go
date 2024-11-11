package main

import (
	dht_kad "application-layer/dht"
	"application-layer/files"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	fmt.Println("Main server started")

	fileRouter := files.InitFileRoutes()
	go dht_kad.StartDHTService()

	// CORS handler
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},        // Frontend's origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"}, // Allowed HTTP methods
		AllowedHeaders:   []string{"Content-Type"},                 // Allowed headers
		AllowCredentials: true,                                     // Allow credentials (cookies, auth headers)
	})

	go startFileAndDHTServer(c.Handler(fileRouter))

	select {}
}

func startFileAndDHTServer(fileRouter http.Handler) {
	port := ":8081"
	http.Handle("/", fileRouter)
	// fmt.Printf("Starting server for file routes and DHT on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
