package utils

import (
    "context"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

// ConnectMongo establishes a connection to MongoDB
func ConnectMongo(uri string) *mongo.Client {
    clientOptions := options.Client().ApplyURI(uri)

    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatalf("Failed to create MongoDB client: %v", err)
    }

    // Ping the database to verify connection
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }

    log.Println("Connected to MongoDB!")
    MongoClient = client
    return client
}

// GetCollection returns a reference to a MongoDB collection
func GetCollection(dbName, collectionName string) *mongo.Collection {
    return MongoClient.Database(dbName).Collection(collectionName)
}
