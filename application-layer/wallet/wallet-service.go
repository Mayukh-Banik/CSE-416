package wallet

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"path/filepath"
)

type WalletServiceInterface interface {
	GenerateNewAddress() (string, error)
	GetPublicKey(address string) (string, error)
	UnlockWallet(passphrase string, duration int) error
	DumpPrivKey(address string) (string, error)
	GenerateNewAddressWithPubKeyAndPrivKey(passphrase string) (string, string, string, float64, error)
}

// WalletService defines wallet-related operations and configurations
type WalletService struct {
	BtcctlPath string
	RpcUser    string
	RpcPass    string
	RpcServer  string
}

// NewWalletService initializes WalletService with provided RPC credentials
func NewWalletService(rpcUser, rpcPass string) *WalletService {

	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Failed to get current working directory: %v", err)
	} else {
		log.Printf("Current working directory: %s", cwd)
	}

	var btcctlPath string

	// 운영 체제에 따라 btcctl 경로 설정
	switch runtime.GOOS {
	case "windows":
		// btcctlPath = `E:\.code-workspace\CSE-416\application-layer\btcd\cmd\btcctl\btcctl.exe` // 절대 경로로 변경
		btcctlPath = `./btcd/cmd/btcctl/btcctl.exe`
	default: // macOS 및 Linux
		btcctlPath = `./btcd/cmd/btcctl/btcctl`
	}

	// 상대 경로가 실제로 존재하는지 확인
	if _, err := os.Stat(btcctlPath); os.IsNotExist(err) {
		log.Printf("btcctl.exe does not exist at path: %s", btcctlPath)
	} else if err != nil {
		log.Printf("Error checking btcctl.exe path: %v", err)
	} else {
		log.Printf("btcctl.exe exists at path: %s", btcctlPath)
	}

	// 상대 경로를 절대 경로로 변환
	absPath, err := filepath.Abs(btcctlPath)
	if err != nil {
		log.Printf("Failed to get absolute path for btcctlPath '%s': %v", btcctlPath, err)
	} else {
		log.Printf("Absolute btcctlPath: %s", absPath)
	}

	// 운영 체제에 따라 경로 구분자 조정 (Windows에서는 필요 시)
	switch runtime.GOOS {
	case "windows":
		absPath = filepath.FromSlash(absPath)
	default:
		// macOS 및 Linux는 별도 처리 필요 없음
	}

	// btcctlPath가 실제로 존재하는지 다시 확인
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		log.Printf("btcctl.exe does not exist at absolute path: %s", absPath)
	} else if err != nil {
		log.Printf("Error checking btcctl.exe at absolute path: %v", err)
	} else {
		log.Printf("btcctl.exe exists at absolute path: %s", absPath)
	}

	return &WalletService{
		BtcctlPath: btcctlPath,
		RpcUser:    rpcUser,
		RpcPass:    rpcPass,
		RpcServer:  "127.0.0.1:8332",
	}
}

// GenerateNewAddress generates a new wallet address using btcctl
func (ws *WalletService) GenerateNewAddress() (string, error) {
	cmd := exec.Command(ws.BtcctlPath, "--wallet", "--rpcuser="+ws.RpcUser, "--rpcpass="+ws.RpcPass, "--rpcserver="+ws.RpcServer, "--notls", "getnewaddress")
	output, err := cmd.CombinedOutput()
	log.Printf("Command Output for GenerateNewAddress: %s\n", output) // 디버그 로그
	if err != nil {
		log.Printf("Command execution failed: %v\nOutput: %s", err, output)
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
	log.Println("Generating new address with passphrase:", passphrase)

	address, err := ws.GenerateNewAddress()
	if err != nil {
		return "", "", "", 0, err
	}
	log.Println("Generated Address:", address)

	pubKey, err := ws.GetPublicKey(address)
	if err != nil {
		return "", "", "", 0, err
	}
	log.Println("Retrieved Public Key:", pubKey)

	// Unlock the wallet before dumping the private key
	err = ws.UnlockWallet(passphrase, 600) // Unlock for 10 minutes (600 seconds)
	if err != nil {
		return "", "", "", 0, err
	}
	log.Println("Wallet unlocked successfully")

	// Retrieve the private key for the generated address
	privateKey, err := ws.DumpPrivKey(address)
	if err != nil {
		return "", "", "", 0, err
	}
	log.Println("Retrieved Private Key:", privateKey)

	// Retrieve the balance for the address
	balance, err := ws.GetBalance(address)
	if err != nil {
		return "", "", "", 0, err
	}
	log.Println("Retrieved Balance:", balance)

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

// SignMessage signs a given message using the specified wallet address
func (ws *WalletService) SignMessage(address string, message string, passphrase string) (string, error) {
	// Unlock the wallet first
	err := ws.UnlockWallet(passphrase, 600)
	if err != nil {
		return "", fmt.Errorf("failed to unlock wallet: %v", err)
	}

	// Run the btcctl command to sign the message
	cmd := exec.Command(ws.BtcctlPath, "--wallet", "--rpcuser="+ws.RpcUser, "--rpcpass="+ws.RpcPass, "--rpcserver="+ws.RpcServer, "--notls", "signmessage", address, message)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to sign message: %v\nOutput: %s", err, output)
	}

	signature := strings.TrimSpace(string(output))
	return signature, nil
}

// VerifySignature uses btcctl to verify the signed challenge.
func (ws *WalletService) VerifySignature(address, challenge, signature string) (bool, error) {
	cmd := exec.Command(ws.BtcctlPath, "--wallet", "--rpcuser="+ws.RpcUser, "--rpcpass="+ws.RpcPass, "--rpcserver="+ws.RpcServer, "--notls", "verifymessage", address, signature, challenge)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to verify signature: %v\nOutput: %s", err, output)
	}

	result := strings.TrimSpace(string(output))
	return result == "true", nil
}

// btcctl --wallet --rpcuser=user --rpcpass=password --rpcserver=127.0.0.1:8332 --notls help
