package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"strings"
	"time"

	"context"
	"go-server/models"
	"go-server/utils"

	"go.mongodb.org/mongo-driver/bson"
)

// Temporary store for challenges (in-memory for now)
var challengeStore = make(map[string]models.ChallengeData)

func GetChallenge(publicKey string) (models.ChallengeData, bool) {
	challenge, exists := challengeStore[publicKey]
	return challenge, exists
}

// FindUserByPublicKey checks if a user exists in the database by their public key.
// It also cleans the public key by removing any newlines.
func FindUserByPublicKey(publicKey string) (*models.User, error) {

	cleanPublicKey := strings.ReplaceAll(publicKey, "\n", "")
	collection := utils.GetCollection("squidcoinDB", "users")
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"public_key": cleanPublicKey}).Decode(&user)
	if err != nil {
		log.Printf("User not found for publicKey: %s", cleanPublicKey)
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// GenerateChallenge generates a random challenge
func GenerateChallenge() (string, error) {
	challengeBytes := make([]byte, 16)
	_, err := rand.Read(challengeBytes) // crypto/rand 사용
	if err != nil {
		log.Printf("Failed to generate challenge")
		return "", errors.New("failed to generate challenge")
	}

	challenge := base64.StdEncoding.EncodeToString(challengeBytes)

	// 생성된 챌린지를 로그로 출력
	log.Printf("Generated challenge: %s", challenge)

	return challenge, nil
}

// StoreChallenge stores the generated challenge in memory with creation time
func StoreChallenge(publicKey, challenge string) error {
	challengeStore[publicKey] = models.ChallengeData{
		Challenge: challenge,
		CreatedAt: time.Now(),
	}
	log.Printf("Challenge stored for publicKey: %s", publicKey)
	return nil
}

// Delete expired Challenges periodically
func CleanupExpiredChallenges(expirationTime time.Duration) {
	for publicKey, data := range challengeStore {
		if time.Since(data.CreatedAt) > expirationTime {
			log.Printf("Challenge expired for publicKey: %s", publicKey)
			delete(challengeStore, publicKey)
		}
	}
}

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

	verified := utils.VerifySignature(parsedPublicKey, challenge.Challenge, signature)

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
