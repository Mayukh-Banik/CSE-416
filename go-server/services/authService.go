package services

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"context"
	"go-server/models"
	"go-server/utils"

	"go.mongodb.org/mongo-driver/bson"
)

// Temporary store for challenges (in-memory for now)
var challengeStore = make(map[string]models.ChallengeData)

var (
	currentChallenge string
	challengeMutex   sync.Mutex
)

// GetChallenge retrieves the current challenge
func GetChallenge() (string, error) {
	challengeMutex.Lock()
	defer challengeMutex.Unlock()
	if currentChallenge == "" {
		return "", errors.New("no challenge stored")
	}
	return currentChallenge, nil
}

// FindUserByPublicKey checks if a user exists in the database by their public key.
// It also cleans the public key by removing any newlines.
func FindUserByPublicKey(publicKey string) (*models.User, error) {
	collection := utils.GetCollection("squidcoinDB", "users")
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"public_key": publicKey}).Decode(&user)
	if err != nil {
		log.Printf("User not found for publicKey")
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// GenerateChallenge generates a new challenge
func GenerateChallenge() (string, error) {
	challengeBytes := make([]byte, 32)
	_, err := rand.Read(challengeBytes)
	if err != nil {
		return "", err
	}

	challenge := base64.StdEncoding.EncodeToString(challengeBytes)
	return challenge, nil
}

// StoreChallenge stores the current challenge
func StoreChallenge(challenge string) {
	challengeMutex.Lock()
	defer challengeMutex.Unlock()
	currentChallenge = challenge
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
func formatPublicKey(key string) string {
	key = strings.TrimSpace(key)
	header := "-----BEGIN PUBLIC KEY-----"
	footer := "-----END PUBLIC KEY-----"

	// 헤더와 푸터 제거
	key = strings.ReplaceAll(key, header, "")
	key = strings.ReplaceAll(key, footer, "")

	// 공백 및 줄바꿈 제거
	key = strings.ReplaceAll(key, " ", "")
	key = strings.ReplaceAll(key, "\n", "")
	key = strings.ReplaceAll(key, "\r", "")

	// 64자마다 줄바꿈 추가
	var formattedKey strings.Builder
	formattedKey.WriteString(header + "\n")
	for i := 0; i < len(key); i += 64 {
		end := i + 64
		if end > len(key) {
			end = len(key)
		}
		formattedKey.WriteString(key[i:end] + "\n")
	}
	formattedKey.WriteString(footer + "\n")

	return formattedKey.String()
}

// VerifySignature verifies if the signature matches the stored challenge
func VerifySignature(publicKeyPem, signatureBase64 string) (bool, error) {

	// 퍼블릭 키 형식 보정
	publicKeyPem = formatPublicKey(publicKeyPem)

	// 현재 저장된 챌린지 가져오기
	challenge, err := GetChallenge()
	if err != nil {
		return false, err
	}

	// 서명 디코딩
	signatureBytes, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return false, errors.New("invalid signature format")
	}

	// 퍼블릭 키 디코딩
	block, _ := pem.Decode([]byte(publicKeyPem))
	if block == nil {
		return false, errors.New("invalid public key PEM")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, errors.New("invalid public key")
	}

	pubKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return false, errors.New("public key is not RSA")
	}

	// 챌린지 디코딩
	challengeBytes, err := base64.StdEncoding.DecodeString(challenge)
	if err != nil {
		return false, errors.New("invalid challenge data")
	}

	// 챌린지 해싱
	hashed := sha256.Sum256(challengeBytes)

	// 서명 검증
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], signatureBytes)
	if err != nil {
		return false, errors.New("signature verification failed")
	}

	// 챌린지 사용 후 삭제
	challengeMutex.Lock()
	currentChallenge = ""
	challengeMutex.Unlock()

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
