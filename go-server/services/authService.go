package services

import (
    "crypto/rand"
    "encoding/base64"
    "errors"
	"log"
	"strings"
    //"time"

    "go-server/models"
    "go-server/utils"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	
	
)

// Temporary store for challenges (in-memory for now)
var challengeStore = make(map[string]string)

// GenerateChallenge generates a challenge and stores it temporarily
func GenerateChallenge(publicKey string) (string, error) {
    // Log the incoming public key for debugging
    log.Printf("Received challenge request for publicKey: %s", publicKey)

	// Clean up the received public key
    cleanPublicKey := strings.ReplaceAll(publicKey, "\n", "")
    
    // Log the cleaned public key for debugging
    log.Printf("Received cleaned publicKey: %s", cleanPublicKey)

    // Lookup user by public key
    collection := utils.GetCollection("squidcoinDB", "users")
    var user models.User
    err := collection.FindOne(context.TODO(), bson.M{"public_key": cleanPublicKey}).Decode(&user)
    if err != nil {
        log.Printf("User not found for publicKey: %s", cleanPublicKey)  // Log if user not found
        return "", errors.New("user not found")
    }

    // Log that the user was found
    log.Printf("User found for publicKey: %s", cleanPublicKey)

    // Generate a random challenge (16 bytes)
    challengeBytes := make([]byte, 16)
    _, err = rand.Read(challengeBytes)
    if err != nil {
        log.Printf("Failed to generate challenge for publicKey: %s", cleanPublicKey)
        return "", errors.New("failed to generate challenge")
    }
    challenge := base64.StdEncoding.EncodeToString(challengeBytes)

    // Log challenge generation success
    log.Printf("Generated challenge for publicKey: %s", cleanPublicKey)

    // Temporarily store the challenge
    challengeStore[user.UserID] = challenge

    return challenge, nil
}

// VerifySignature verifies if the signature matches the stored challenge
func VerifySignature(userID, signature string) (bool, error) {
    // Fetch user from MongoDB
    collection := utils.GetCollection("squidcoinDB", "users")
    var user models.User
    err := collection.FindOne(context.TODO(), bson.M{"user_id": userID}).Decode(&user)
    if err != nil {
        return false, errors.New("user not found")
    }

    // Get the stored challenge
    challenge, exists := challengeStore[userID]
    if !exists {
        return false, errors.New("challenge not found or expired")
    }

    // Verify the signature using public key and challenge
    publicKey, err := utils.ParsePublicKey(user.PublicKey)
    if err != nil {
        return false, errors.New("invalid public key")
    }

    verified := utils.VerifySignature(publicKey, challenge, signature)
    if !verified {
        return false, errors.New("signature verification failed")
    }

    // Cleanup: remove challenge after successful verification
    delete(challengeStore, userID)

    return true, nil
}
