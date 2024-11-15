package wallet

import (
    "testing"
)

func TestGenerateNewAddress(t *testing.T) {
	PrintCurrentDir()
    ws := NewWalletService() // No parameters needed as it loads from .env
    address, err := ws.GenerateNewAddress()
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if address == "" {
        t.Fatalf("Expected a new address, got empty string")
    }
    t.Logf("Generated Address: %s", address)
}
