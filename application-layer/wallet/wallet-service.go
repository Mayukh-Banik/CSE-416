package wallet

import (
	"os"
    "fmt"
    "os/exec"
)

// WalletService defines the structure for wallet-related operations
type WalletService struct {
    BtcctlPath string
    RpcUser    string
    RpcPass    string
    RpcServer  string
}

func NewWalletService() *WalletService {
    
    return &WalletService{
        BtcctlPath: `../btcd/cmd/btcctl/btcctl`,
        RpcUser:    "user",
        RpcPass:    "password",
        RpcServer:  "127.0.0.1:8332",
    }
}


// GenerateNewAddress generates a new wallet address using btcctl
func (ws *WalletService) GenerateNewAddress() (string, error) {
    cmd := exec.Command(`../btcd/cmd/btcctl/btcctl`, "--wallet", "--rpcuser="+ws.RpcUser, "--rpcpass="+ws.RpcPass, "--rpcserver="+ws.RpcServer, "--notls", "getnewaddress")
    output, err := cmd.CombinedOutput()
	fmt.Printf("Command Output: %s\n", output) // Debug: Print the output of the command
    if err != nil {
        return "", fmt.Errorf("failed to generate new address: %v", err)
    }
    return string(output), nil
}

func PrintCurrentDir() {
    dir, err := os.Getwd()
    if err != nil {
        fmt.Println("Error getting current directory:", err)
    }
    fmt.Println("Current Directory:", dir)
}
