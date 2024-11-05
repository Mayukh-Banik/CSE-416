package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FileMetadata represents the metadata for an uploaded file
type FileMetadata struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Size        int64  `json:"size"`
	Description string `json:"description"`
	Hash        string `json:"hash"`
	IsPublished string `json:"isPublished"`
}

var mongoClient *mongo.Client
var collection *mongo.Collection

func init() {
	// Set up MongoDB connection
	var err error
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	mongoClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	collection = mongoClient.Database("file_db").Collection("files") // Replace "file_db" and "files" with your database and collection names
}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		Size        int64  `json:"size"`
		Description string `json:"description"`
		Hash        string `json:"hash"`
		
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create file metadata from the request
	fileMetadata := FileMetadata{
		Name:        requestBody.Name,
		Type:        requestBody.Type,
		Size:        requestBody.Size,
		Description: requestBody.Description,
		Hash:        requestBody.Hash,
	}

	// Insert the file metadata into the MongoDB collection
	_, err := collection.InsertOne(context.TODO(), fileMetadata)
	if err != nil {
		http.Error(w, "Failed to upload file metadata", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	response := map[string]string{"message": "File metadata uploaded successfully", "hash": requestBody.Hash}
	w.Header().Set("Content-Type", "application/json") // Set content type to JSON
	w.WriteHeader(http.StatusOK)                       // Set response status
	json.NewEncoder(w).Encode(response)                // Encode response as JSON
}

func main() {
	http.HandleFunc("/upload", uploadFileHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
