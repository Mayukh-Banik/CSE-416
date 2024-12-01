package wallet

import (
	"errors"
)

type MockWalletService struct {
	GenerateNewAddressFn                func() (string, error)
	GetPublicKeyFn                      func(address string) (string, error)
	UnlockWalletFn                      func(passphrase string, duration int) error
	DumpPrivKeyFn                       func(address string) (string, error)
	GenerateNewAddressWithPubKeyAndPrivKeyFn func(passphrase string) (string, string, string, error)
}

func (m *MockWalletService) GenerateNewAddress() (string, error) {
	if m.GenerateNewAddressFn != nil {
		return m.GenerateNewAddressFn()
	}
	return "", errors.New("GenerateNewAddress not implemented")
}

func (m *MockWalletService) GetPublicKey(address string) (string, error) {
	if m.GetPublicKeyFn != nil {
		return m.GetPublicKeyFn(address)
	}
	return "", errors.New("GetPublicKey not implemented")
}

func (m *MockWalletService) UnlockWallet(passphrase string, duration int) error {
	if m.UnlockWalletFn != nil {
		return m.UnlockWalletFn(passphrase, duration)
	}
	return errors.New("UnlockWallet not implemented")
}

func (m *MockWalletService) DumpPrivKey(address string) (string, error) {
	if m.DumpPrivKeyFn != nil {
		return m.DumpPrivKeyFn(address)
	}
	return "", errors.New("DumpPrivKey not implemented")
}

func (m *MockWalletService) GenerateNewAddressWithPubKeyAndPrivKey(passphrase string) (string, string, string, error) {
	if m.GenerateNewAddressWithPubKeyAndPrivKeyFn != nil {
		return m.GenerateNewAddressWithPubKeyAndPrivKeyFn(passphrase)
	}
	return "", "", "", errors.New("GenerateNewAddressWithPubKeyAndPrivKey not implemented")
}
