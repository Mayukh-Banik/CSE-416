package wallet

import (
	"testing"
)

func TestGenerateNewAddress(t *testing.T) {
	ws, err := NewWalletService("user", "password") // No parameters needed as it loads from .env
	if err != nil {
		t.Fatalf("NewWalletService error got %v", err)
	}
	// address, err := ws.GenerateNewAddress()    // Match the return values
	address, pubKey, err := ws.GenerateNewAddressWithPubKey()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if address == "" {
		t.Fatalf("Expected a new address, got empty string")
	}
	t.Logf("Generated Address: %s\nPublic Key: %s\n", address, pubKey)
}
