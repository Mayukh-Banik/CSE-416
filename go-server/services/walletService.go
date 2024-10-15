package services

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
)

// WalletResponse struct to send back both public and private keys
type WalletResponse struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"` // Include private key here
}

// GenerateWallet creates an RSA key pair and returns the public and private keys
func GenerateWallet() (WalletResponse, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Printf("Error generating private key: %v", err)
		return WalletResponse{}, errors.New("failed to generate keys")
	}

	// Extract public key in PEM format
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return WalletResponse{}, err
	}

	// PEM encode the public key
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	})

	// PEM encode the private key
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	// Return both keys
	return WalletResponse{
		PublicKey:  string(publicKeyPEM),
		PrivateKey: string(privateKeyPEM),
	}, nil
}
