package wallet

import (
	"application-layer/models"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var jwtKey = []byte("your_secret_key")

type UserService struct {
	Users  map[string]*models.User // In-memory user store
	Nonces map[string]string       // Nonce store for login challenges
	mutex  sync.Mutex              // Mutex for concurrent access
}

// WalletService defines wallet-related operations and configurations
type WalletService struct {
	BtcctlPath string
	RpcUser    string
	RpcPass    string
	RpcServer  string
}

type Claims struct {
	UUID string `json:"uuid"`
	jwt.StandardClaims
}
 
// NewUserService initializes the UserService
func NewUserService() *UserService {
	return &UserService{
		Users: make(map[string]*models.User),
	}
}

// SignUp handles the signup logic
func (us *UserService) SignUp(walletService WalletService, passphrase string) (*models.User, string, error) {
	// Generate a new wallet address and public key
	// address, pubKey, err := walletService.GenerateNewAddressWithPubKey()

	// if err != nil {
	// 	return nil, fmt.Errorf("failed to generate wallet address and public key: %v", err)
	// }

    // Generate a new wallet address, public key, and private key
    address, pubKey, privateKey, balance, err := walletService.GenerateNewAddressWithPubKeyAndPrivKey(passphrase)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate wallet address, public key, and private key: %v", err)
	}

	// Create a UUID for the user
	userID := uuid.NewString()

	// Create a new User struct
	newUser := &models.User{
		UUID:         userID,
		Address:      address,
		PublicKey:    pubKey,
		CreatedDate:  time.Now(),
		Metadata:     make(map[string]string),
		Balance:      balance,
		Transactions: []string{},
	}

	// Store the user in memory
	us.Users[userID] = newUser

	return newUser, privateKey, nil
}

// GetUser retrieves a user by UUID
func (us *UserService) GetUser(uuid string) (*models.User, error) {
	user, exists := us.Users[uuid]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// InitiateLogin generates a nonce for the user to sign
func (us *UserService) InitiateLogin(userUUID string) (string, error) {
	us.mutex.Lock()
	defer us.mutex.Unlock()

	// Check if user exists
	if _, exists := us.Users[userUUID]; !exists {
		return "", errors.New("user not found")
	}

	// Generate a random nonce
	nonceBytes := make([]byte, 16)
	_, err := rand.Read(nonceBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate nonce: %v", err)
	}
	nonce := hex.EncodeToString(nonceBytes)

	// Store the nonce with an expiration time (e.g., 5 minutes)
	us.Nonces[userUUID] = nonce

	return nonce, nil
}

// CompleteLogin verifies the signed nonce and returns a JWT token if successful
func (us *UserService) CompleteLogin(userUUID, signature, message string) (string, error) {
	us.mutex.Lock()
	nonce, exists := us.Nonces[userUUID]
	us.mutex.Unlock()

	if !exists {
		return "", errors.New("no login initiation found for this user")
	}

	// Verify that the message is the nonce
	if message != nonce {
		return "", errors.New("invalid message")
	}

	// Get the user's public key
	user, err := us.GetUser(userUUID)
	if err != nil {
		return "", err
	}

	// Decode the public key
	pubKeyBytes, err := hex.DecodeString(user.PublicKey)
	if err != nil {
		return "", errors.New("invalid public key format")
	}

	if len(pubKeyBytes) != 64 {
		return "", errors.New("invalid public key length")
	}

	pubKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     new(big.Int).SetBytes(pubKeyBytes[:32]),
		Y:     new(big.Int).SetBytes(pubKeyBytes[32:]),
	}

	// Decode the signature
	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return "", errors.New("invalid signature format")
	}

	if len(signatureBytes) != 64 {
		return "", errors.New("invalid signature length")
	}

	rSig := new(big.Int).SetBytes(signatureBytes[:32])
	sSig := new(big.Int).SetBytes(signatureBytes[32:])

	// Hash the message
	messageHash := sha256.Sum256([]byte(message))

	// Verify the signature
	valid := ecdsa.Verify(&pubKey, messageHash[:], rSig, sSig)
	if !valid {
		return "", errors.New("invalid signature")
	}

	// Remove the nonce as it's single-use
	us.mutex.Lock()
	delete(us.Nonces, userUUID)
	us.mutex.Unlock()

	// Create JWT token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UUID: userUUID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return tokenString, nil
}
