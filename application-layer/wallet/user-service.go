package wallet

import (
	"application-layer/models"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type UserService struct {
	Users  map[string]*models.User // In-memory user store
	Challenges map[string]*models.Challenge // In-memory challenge store
	Mutex     sync.Mutex                    // Mutex for concurrent access
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
        Challenges: make(map[string]*models.Challenge),
	}
}

// SignUp handles the signup logic
func (us *UserService) SignUp(walletService WalletService, passphrase string) (*models.User, string, error) {
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

// GenerateChallenge creates a new challenge for a given wallet address.
func (us *UserService) GenerateChallenge(address string) (*models.Challenge, error) {
    us.Mutex.Lock()
    defer us.Mutex.Unlock()

    challenge := uuid.NewString() // Create a unique string challenge.
    expiry := time.Now().Add(5 * time.Minute) // Set challenge expiry time.

    newChallenge := &models.Challenge{
        Address:   address,
        Challenge: challenge,
        Expiry:    expiry,
    }

    us.Challenges[address] = newChallenge

    return newChallenge, nil
}

// GetChallenge retrieves a stored challenge for the given wallet address.
func (us *UserService) GetChallenge(address string) (*models.Challenge, error) {
    us.Mutex.Lock()
    defer us.Mutex.Unlock()

    challenge, exists := us.Challenges[address]
    if !exists || time.Now().After(challenge.Expiry) {
        return nil, errors.New("challenge expired or not found")
    }

    return challenge, nil
}

// RemoveChallenge deletes the challenge after successful verification.
func (us *UserService) RemoveChallenge(address string) {
    us.Mutex.Lock()
    defer us.Mutex.Unlock()
    delete(us.Challenges, address)
}

