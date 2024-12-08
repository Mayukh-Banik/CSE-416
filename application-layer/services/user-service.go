package services

import (
	"errors"
	"fmt"
)

// WalletService 구조체
type WalletService struct {
	rpcUser string
	rpcPass string
}

// NewWalletService는 WalletService를 초기화합니다.
func NewWalletService(rpcUser, rpcPass string) (*WalletService, error) {
	if rpcUser == "" || rpcPass == "" {
		return nil, errors.New("RPC credentials are required")
	}
	return &WalletService{rpcUser: rpcUser, rpcPass: rpcPass}, nil
}

// StartBtcd는 btcd를 시작합니다.
func (ws *WalletService) StartBtcd() {
	fmt.Println("Starting btcd with RPC credentials...")
	// 실제 btcd 실행 로직을 여기에 추가
}

// StopBtcd는 btcd를 중지합니다.
func (ws *WalletService) StopBtcd() {
	fmt.Println("Stopping btcd...")
	// 실제 btcd 종료 로직을 여기에 추가
}

// StartBtcwallet는 btcwallet을 시작합니다.
func (ws *WalletService) StartBtcwallet() {
	fmt.Println("Starting btcwallet with RPC credentials...")
	// 실제 btcwallet 실행 로직을 여기에 추가
}

// StopBtcwallet는 btcwallet을 중지합니다.
func (ws *WalletService) StopBtcwallet() {
	fmt.Println("Stopping btcwallet...")
	// 실제 btcwallet 종료 로직을 여기에 추가
}

// var jwtKey = []byte("your_secret_key")

// type UserService struct {
// 	Users      map[string]*models.User      // In-memory user store
// 	Challenges map[string]*models.Challenge // In-memory challenge store
// 	Mutex      sync.Mutex                   // Mutex for concurrent access
// }

// type Claims struct {
// 	UUID string `json:"uuid"`
// 	jwt.StandardClaims
// }

// // NewUserService initializes the UserService
// func NewUserService() *UserService {
// 	return &UserService{
// 		Users:      make(map[string]*models.User),
// 		Challenges: make(map[string]*models.Challenge),
// 	}
// }

// // SignUp handles the signup logic
// func (us *UserService) SignUp(walletService WalletService, passphrase string) (*models.User, string, error) {
// 	// Generate a new wallet address, public key, and private key
// 	address, privateKey, err := walletService.GenerateNewAddressWithPubKeyAndPrivKey(passphrase)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("failed to generate wallet address, public key, and private key: %v", err)
// 	}

// 	// Create a UUID for the user
// 	userID := uuid.NewString()

// 	// Create a new User struct
// 	newUser := &models.User{
// 		UUID:        userID,
// 		Address:     address,
// 		CreatedDate: time.Now(),
// 	}

// 	// Store the user in memory
// 	us.Users[userID] = newUser

// 	return newUser, privateKey, nil
// }

// // GetUser retrieves a user by UUID
// func (us *UserService) GetUser(uuid string) (*models.User, error) {
// 	user, exists := us.Users[uuid]
// 	if !exists {
// 		return nil, errors.New("user not found")
// 	}
// 	return user, nil
// }

// // GenerateChallenge creates a new challenge for a given wallet address.
// func (us *UserService) GenerateChallenge(address string) (*models.Challenge, error) {
// 	us.Mutex.Lock()
// 	defer us.Mutex.Unlock()

// 	challenge := uuid.NewString()             // Create a unique string challenge.
// 	expiry := time.Now().Add(5 * time.Minute) // Set challenge expiry time.

// 	newChallenge := &models.Challenge{
// 		Address:   address,
// 		Challenge: challenge,
// 		Expiry:    expiry,
// 	}

// 	us.Challenges[address] = newChallenge

// 	return newChallenge, nil
// }

// // GetChallenge retrieves a stored challenge for the given wallet address.
// func (us *UserService) GetChallenge(address string) (*models.Challenge, error) {
// 	us.Mutex.Lock()
// 	defer us.Mutex.Unlock()

// 	challenge, exists := us.Challenges[address]
// 	if !exists || time.Now().After(challenge.Expiry) {
// 		return nil, errors.New("challenge expired or not found")
// 	}

// 	return challenge, nil
// }

// // RemoveChallenge deletes the challenge after successful verification.
// func (us *UserService) RemoveChallenge(address string) {
// 	us.Mutex.Lock()
// 	defer us.Mutex.Unlock()
// 	delete(us.Challenges, address)
// }
