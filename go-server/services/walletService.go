package services

import (
    "context"
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "errors"
    "log"
    "time"
	"strings"

    "go-server/models"
    "go-server/utils"

    //"go.mongodb.org/mongo-driver/mongo"
)

// WalletResponse struct to send back both public and private keys
type WalletResponse struct {
    PublicKey  string `json:"public_key"`
    PrivateKey string `json:"private_key"`
}

// GenerateWallet creates an RSA key pair, stores the user in MongoDB, and returns the response
func GenerateWallet() (WalletResponse, error) {
    // Generate RSA key pair
    privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        log.Printf("Error generating private key: %v", err)
        return WalletResponse{}, errors.New("failed to generate keys")
    }

    // Encode public key to PEM format
    pubASN1, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
    if err != nil {
        return WalletResponse{}, err
    }
    publicKeyPEM := pem.EncodeToMemory(&pem.Block{
        Type:  "PUBLIC KEY",
        Bytes: pubASN1,
    })

	// Clean the public key (remove newlines for storage consistency)
    cleanPublicKey := strings.ReplaceAll(string(publicKeyPEM), "\n", "")

    // Encode private key to PEM format
    privASN1 := x509.MarshalPKCS1PrivateKey(privateKey)
    privateKeyPEM := pem.EncodeToMemory(&pem.Block{
        Type:  "RSA PRIVATE KEY",
        Bytes: privASN1,
    })

    // Create a new user object
    user := models.User{
        PublicKey:   cleanPublicKey,
        CreatedAt: time.Now(),
    }

    // Store the user in MongoDB
    collection := utils.GetCollection("squidcoinDB", "users")
    _, err = collection.InsertOne(context.TODO(), user)
    if err != nil {
        return WalletResponse{}, err
    }

    // Prepare the response
    walletResponse := WalletResponse{
        PublicKey:  user.PublicKey,
        PrivateKey: string(privateKeyPEM),
    }

    return walletResponse, nil
}