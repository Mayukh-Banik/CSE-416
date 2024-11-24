package wallet

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type WalletServiceInterface interface {
	GenerateNewAddress() (string, error)
	GetPublicKey(address string) (string, error)
	UnlockWallet(passphrase string, duration int) error
	DumpPrivKey(address string) (string, error)
	GenerateNewAddressWithPubKeyAndPrivKey(passphrase string) (string, string, string, error)
}

// NewWalletService initializes WalletService with provided RPC credentials
func NewWalletService(rpcUser, rpcPass string) *WalletService {
	return &WalletService{
		BtcctlPath: `../btcd/cmd/btcctl/btcctl.exe`,
		RpcUser:    rpcUser,
		RpcPass:    rpcPass,
		RpcServer:  "127.0.0.1:8332",
	}
}

// GenerateNewAddress generates a new wallet address using btcctl
func (ws *WalletService) GenerateNewAddress() (string, error) {
	cmd := exec.Command(`../btcd/cmd/btcctl/btcctl`, "--wallet", "--rpcuser="+ws.RpcUser, "--rpcpass="+ws.RpcPass, "--rpcserver="+ws.RpcServer, "--notls", "getnewaddress")
	output, err := cmd.CombinedOutput()
	fmt.Printf("@@Command Output: %s\n", output) // Debug: Print the output of the command
	if err != nil {
		return "", fmt.Errorf("failed to generate new address: %v", err)
	}
	address := strings.TrimSpace(string(output))
	return address, nil
}

// GenerateNewAddressWithPubKey generates a new address and retrieves its public key
func (ws *WalletService) GenerateNewAddressWithPubKey() (string, string, error) {
	address, err := ws.GenerateNewAddress()
	if err != nil {
		return "", "", err
	}

	pubKey, err := ws.GetPublicKey(address)
	if err != nil {
		return "", "", err
	}

	return address, pubKey, nil
}

// AddressInfo 구조체는 getaddressinfo 명령어의 JSON 출력을 매핑합니다.
type AddressInfo struct {
	Address      string `json:"address"`
	ScriptPubKey string `json:"scriptPubKey"`
	PubKey       string `json:"pubkey"`
	// 필요에 따라 다른 필드도 추가할 수 있습니다.
}

// GetPublicKey retrieves the public key for a given address using validateaddress
func (ws *WalletService) GetPublicKey(address string) (string, error) {
	args := []string{
		"--wallet",
		"--rpcuser=" + ws.RpcUser,
		"--rpcpass=" + ws.RpcPass,
		"--rpcserver=" + ws.RpcServer,
		"--notls",
		"validateaddress",
		address,
	}
	cmd := exec.Command(ws.BtcctlPath, args...)
	output, err := cmd.CombinedOutput()
	fmt.Printf("GetAddressInfo Output: %s\n", output) // Debug: Print the output of the command
	if err != nil {
		fmt.Printf("Command failed with error: %v\n", err)
		fmt.Printf("Command Output: %s\n", output)
		return "", fmt.Errorf("failed to get address info: %v", err)
	}

	var info AddressInfo
	if err := json.Unmarshal(output, &info); err != nil {
		return "", fmt.Errorf("failed to parse address info: %v", err)
	}

	if info.PubKey == "" {
		return "", fmt.Errorf("public key not found for address: %s", address)
	}

	return info.PubKey, nil
}

// UnlockWallet unlocks the wallet to allow access to private keys
func (ws *WalletService) UnlockWallet(passphrase string, duration int) error {
	cmd := exec.Command(ws.BtcctlPath, "--wallet", "--rpcuser="+ws.RpcUser, "--rpcpass="+ws.RpcPass, "--rpcserver="+ws.RpcServer, "--notls", "walletpassphrase", passphrase, fmt.Sprintf("%d", duration))
	output, err := cmd.CombinedOutput() // Need to make note of duration here
	if err != nil {
		return fmt.Errorf("failed to unlock wallet: %v\nOutput: %s", err, output)
	}
	return nil
}

// DumpPrivKey retrieves the private key for a given wallet address
func (ws *WalletService) DumpPrivKey(address string) (string, error) {
	cmd := exec.Command(ws.BtcctlPath, "--wallet", "--rpcuser="+ws.RpcUser, "--rpcpass="+ws.RpcPass, "--rpcserver="+ws.RpcServer, "--notls", "dumpprivkey", address)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve private key: %v\nOutput: %s", err, output)
	}
	privateKey := strings.TrimSpace(string(output))
	return privateKey, nil
}

// GenerateNewAddressWithPubKeyAndPrivKey generates a new address, retrieves its public key, and retrieves the private key
func (ws *WalletService) GenerateNewAddressWithPubKeyAndPrivKey(passphrase string) (string, string, string, float64, error) {
	address, err := ws.GenerateNewAddress()
	if err != nil {
		return "", "", "", 0, err
	}

	pubKey, err := ws.GetPublicKey(address)
	if err != nil {
		return "", "", "", 0, err
	}

	// Unlock the wallet before dumping the private key
	err = ws.UnlockWallet(passphrase, 600) // Unlock for 10 minutes (600 seconds)
	if err != nil {
		return "", "", "", 0, err
	}

	// Retrieve the private key for the generated address
	privateKey, err := ws.DumpPrivKey(address)
	if err != nil {
		return "", "", "", 0, err
	}

	// Retrieve the balance for the address
	balance, err := ws.GetBalance(address)
	if err != nil {
		return "", "", "", 0, err
	}

	return address, pubKey, privateKey, balance, nil
}

// Function to get the balance for a specific address
func (ws *WalletService) GetBalance(address string) (float64, error) {
	cmd := exec.Command(ws.BtcctlPath, "--wallet", "--rpcuser="+ws.RpcUser, "--rpcpass="+ws.RpcPass, "--rpcserver="+ws.RpcServer, "--notls", "getreceivedbyaddress", address, "1")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve balance: %v\nOutput: %s", err, output)
	}
	balance := strings.TrimSpace(string(output))

	// Convert balance to float64
	var balanceValue float64
	_, err = fmt.Sscanf(balance, "%f", &balanceValue)
	if err != nil {
		return 0, fmt.Errorf("failed to parse balance: %v", err)
	}

	return balanceValue, nil
}

// btcctl --wallet --rpcuser=user --rpcpass=password --rpcserver=127.0.0.1:8332 --notls help
