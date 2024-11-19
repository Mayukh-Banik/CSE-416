package wallet

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// NewWalletService initializes WalletService with provided RPC credentials
func NewWalletService(rpcUser, rpcPass string) *WalletService {
	return &WalletService{
		BtcctlPath: `../btcd/cmd/btcctl/btcctl`,
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

// btcctl --wallet --rpcuser=user --rpcpass=password --rpcserver=127.0.0.1:8332 --notls help

