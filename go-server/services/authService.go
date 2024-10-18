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
    // Clean up the received public key
    cleanPublicKey := strings.ReplaceAll(publicKey, "\n", "")
    
    // Lookup user by public key in MongoDB
    collection := utils.GetCollection("squidcoinDB", "users")
    var user models.User
    err := collection.FindOne(context.TODO(), bson.M{"public_key": cleanPublicKey}).Decode(&user)
    if err != nil {
        log.Printf("User not found for publicKey: %s", cleanPublicKey)
        return "", errors.New("user not found")
    }

    // Generate the challenge (same logic as before)
    challengeBytes := make([]byte, 16)
    _, err = rand.Read(challengeBytes)
    if err != nil {
        log.Printf("Failed to generate challenge for publicKey: %s", cleanPublicKey)
        return "", errors.New("failed to generate challenge")
    }
    challenge := base64.StdEncoding.EncodeToString(challengeBytes)

    // Store the challenge in memory for the user
    challengeStore[user.ID.Hex()] = challenge

    return challenge, nil
}

// VerifySignature verifies if the signature matches the stored challenge
// VerifySignature verifies if the signature matches the stored challenge
func VerifySignature(publicKey, signature string) (bool, error) {
    // Clean up the public key
    cleanPublicKey := strings.ReplaceAll(publicKey, "\n", "")
    
    log.Printf("Received publicKey: %s", cleanPublicKey)
    log.Printf("Received signature: %s", signature)

    // Fetch user by public key
    collection := utils.GetCollection("squidcoinDB", "users")
    var user models.User
    err := collection.FindOne(context.TODO(), bson.M{"public_key": cleanPublicKey}).Decode(&user)
    if err != nil {
        log.Printf("User not found for publicKey: %s", cleanPublicKey)
        return false, errors.New("user not found")
    }

    log.Printf("User found for publicKey: %s", user.PublicKey)

    // Get the stored challenge
    challenge, exists := challengeStore[user.ID.Hex()]
    if !exists {
        log.Printf("Challenge not found for ID: %s", user.ID.Hex())
        return false, errors.New("challenge not found or expired")
    }

    log.Printf("Challenge found for ID: %s", challenge)

    // Verify the signature using public key and challenge
    parsedPublicKey, err := utils.ParsePublicKey(user.PublicKey)
    if err != nil {
        log.Printf("Invalid public key for ID: %s", user.ID.Hex())
        return false, errors.New("invalid public key")
    }

    verified := utils.VerifySignature(parsedPublicKey, challenge, signature)
    if !verified {
        log.Printf("Signature verification failed for ID: %s", user.ID.Hex())
        return false, errors.New("signature verification failed")
    }

    // Cleanup: remove challenge after successful verification
    delete(challengeStore, user.ID.Hex())
    log.Printf("Signature verification succeeded for ID: %s", user.ID.Hex())

    return true, nil
}

// LoginWithWalletID checks if the walletId (public key) exists in the database and logs in the user
func LoginWithWalletID(walletID string) (bool, error) {
	// Clean up the wallet ID (public key)
	cleanWalletID := strings.ReplaceAll(walletID, "\n", "")

	// Fetch user by walletID (public key)
	collection := utils.GetCollection("squidcoinDB", "users")
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"public_key": cleanWalletID}).Decode(&user)
	if err != nil {
		return false, errors.New("user not found")
	}

	// If user is found, return success
	return true, nil
}