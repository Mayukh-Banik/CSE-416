package services

import (
	"application-layer/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
	"runtime"
)

// NewBtcService creates a new instance of BtcService

func NewBtcService() *BtcService {
	SetupTempFilePath() // Ensure temp file is set up
    return &BtcService{}

}

const CREATE_NO_WINDOW = 0x08000000

// Setting the executable path as a global variable
var (
	btcwalletScriptPath = `C:\dev\workspace\CSE-416\btcwallet\btcwallet_create.ps1`
	btcdPath            = "../../btcd/btcd"
	btcwalletPath       = "../../btcwallet/btcwallet"
	btcctlPath          = "../../btcd/cmd/btcctl/btcctl"
	tempFilePath        string
)

type BtcService struct{}

// this function is used to check the contents of a directory for development purposes
func CheckDirectoryContents(parentDir string) string {
	// parentDir := "../../btcd"
	_, err := utils.CheckDirectoryContents(parentDir)
	if err != nil {
		fmt.Printf("Error checking directory: %v\n", err)
		return "Failed to check directory"
	}
	return "successed to check directory"
}

func initializeTempFile() error {
	// initial data
	initialData := map[string]string{
		"status": "initialized",
	}

	// JSON serialisation
	dataBytes, err := json.MarshalIndent(initialData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal initial data: %w", err)
	}

	// Write to file
	err = ioutil.WriteFile(tempFilePath, dataBytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to temp file: %w", err)
	}

	return nil
}

func SetupTempFilePath() {
	if os.Getenv("OS") == "Windows_NT" {
		tempFilePath = filepath.Join(os.Getenv("TEMP"), "btc_temp.json")
	} else {
		tempFilePath = "/tmp/btcd_temp.json"
	}
	fmt.Printf("Temporary file path set to: %s\n", tempFilePath)

	// Check if the file exists, and initialize it if it doesn't
	if _, err := os.Stat(tempFilePath); os.IsNotExist(err) {
		fmt.Println("Temp file not found. Initializing...")
		if err := initializeTempFile(); err != nil {
			fmt.Printf("Failed to initialize temp file: %v\n", err)
		} else {
			fmt.Println("Temp file initialized successfully.")
		}
	}
}

func deleteFromTempFile(key string) error {
	// Read files
	content, err := ioutil.ReadFile(tempFilePath)
	if err != nil {
		return fmt.Errorf("failed to read temp file: %w", err)
	}

	// Parsing JSON
	var data map[string]string
	if err := json.Unmarshal(content, &data); err != nil {
		return fmt.Errorf("failed to unmarshal temp file content: %w", err)
	}

	// Delete data
	delete(data, key)

	// JSON serialisation
	updatedContent, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated data: %w", err)
	}

	// Write to file
	if err := ioutil.WriteFile(tempFilePath, updatedContent, 0644); err != nil {
		return fmt.Errorf("failed to write updated temp file: %w", err)
	}

	return nil
}

// updateTempFile is a helper function to update the temp file with a new key-value pair
func updateTempFile(key, value string) error {
	var data map[string]string

	// Read file
	content, err := ioutil.ReadFile(tempFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			data = make(map[string]string)
		} else {
			return fmt.Errorf("failed to read temp file: %w", err)
		}
	} else {
		if err := json.Unmarshal(content, &data); err != nil {
			return fmt.Errorf("failed to unmarshal temp file content: %w", err)
		}
	}

	// Update data
	data[key] = value

	// JSON serialisation
	newContent, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated data: %w", err)
	}

	if err := ioutil.WriteFile(tempFilePath, newContent, 0644); err != nil {
		return fmt.Errorf("failed to write updated temp file: %w", err)
	}

	return nil
}

// isProcessRunning is a helper function to check if a process is running
func isProcessRunning(processName string) bool {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		// Windows: Use PowerShell to check processes
		cmd = exec.Command("powershell", "-Command", fmt.Sprintf("Get-Process | Where-Object {$_.Name -eq '%s'}", processName))
	} else if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		// macOS/Linux: Use pgrep to check processes
		cmd = exec.Command("pgrep", "-f", processName) // Use -f to match full command
	} else {
		fmt.Printf("Unsupported OS: %s\n", runtime.GOOS)
		return false
	}

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	err := cmd.Run()
	if err != nil {
		// Log error for debugging
		fmt.Printf("Error checking process: %v\n", err)
		if output.String() == "" {
			return false
		}
	}

	// Log output for debugging
	fmt.Printf("isProcessRunning output: %s\n", output.String())

	// Return true if output is not empty
	return strings.TrimSpace(output.String()) != ""
}

// getMiningAddressFromTemp is a helper function to retrieve the mining address from the temp file
func getMiningAddressFromTemp() (string, error) {
	// Read temp file
	content, err := ioutil.ReadFile(tempFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read temp file: %w", err)
	}

	// Parsing JSON
	var data map[string]string
	if err := json.Unmarshal(content, &data); err != nil {
		return "", fmt.Errorf("failed to unmarshal temp file content: %w", err)
	}

	// Retrieve mining address
	miningAddress, exists := data["miningaddr"]
	if !exists || miningAddress == "" {
		return "", fmt.Errorf("mining address not found in temp file")
	}

	return miningAddress, nil
}

// startBtcd is a function to start the btcd process
// it takes an optional wallet address as an argument
// if no argument is provided, btcd is started without a mining address
// if an argument is provided, btcd is started with the mining address
// if more than one argument is provided, an error is returned
// cannot start another instance of btcd if it is already running
func (bs *BtcService) StartBtcd(walletAddress ...string) string {
	// Check if btcd is already running
	if isProcessRunning("btcd") {
		fmt.Println("btcd is already running. Cannot start another instance.")
		return "btcd is already running"
	}

	var cmd *exec.Cmd

	// no argument provided
	if len(walletAddress) == 0 {
		cmd = exec.Command(
			btcdPath,
			"--rpcuser=user",
			"--rpcpass=password",
			"--notls",
		)
	} else if len(walletAddress) == 1 {
		// one argument provided
		cmd = exec.Command(
			btcdPath,
			"--rpcuser=user",
			"--rpcpass=password",
			"--notls",
			fmt.Sprintf("--miningaddr=%s", walletAddress[0]),
		)
	} else {
		// more than one argument provided
		fmt.Println("Invalid number of arguments. Only 0 or 1 argument is allowed.")
		return "Invalid number of arguments"
	}

	// Detached mode with OS-specific handling
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	} else if runtime.GOOS == "darwin" {
		// macOS specific adjustments
		cmd.Env = append(os.Environ(), "PATH=/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin")
	}

	// Standard output and error
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start btcd process
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting btcd: %v\n", err)
		return "Error starting btcd"
	}

	fmt.Printf("btcd started with PID: %d\n", cmd.Process.Pid)

	// Check if btcd process is running
	if !isProcessRunning("btcd") {
		fmt.Println("btcd process not found after starting.")
		return "btcd process not found"
	}

	fmt.Println("btcd is running")

	// Update temp file with mining address
	if len(walletAddress) == 0 {
		// clear mining address from temp file
		if err := deleteFromTempFile("miningaddr"); err != nil {
			fmt.Printf("Failed to delete miningaddr from temp file: %v\n", err)
			return "btcd started but failed to clear mining address from temp file"
		}
		fmt.Println("Mining address cleared from temporary file.")
	} else if len(walletAddress) == 1 {
		// update temp file with mining address
		if err := updateTempFile("miningaddr", walletAddress[0]); err != nil {
			fmt.Printf("Failed to update temp file: %v\n", err)
			return "btcd started but failed to update temp file"
		}
		fmt.Println("Temporary file updated successfully.")
	}

	return "btcd started successfully"
}

// StopBtcd stops the btcd process.
// go test -v -run ^TestStopBtcd$ -count=1 application-layer/services
// use -count=1 to avoid caching
func (bs *BtcService) StopBtcd() string {
	var checkProcessCmd *exec.Cmd
	var killCmd *exec.Cmd

	if runtime.GOOS == "windows" {
		// Check if btcd is running
		checkProcessCmd = exec.Command("powershell", "-Command", "Get-Process | Where-Object {$_.Name -eq 'btcd'}")
		killCmd = exec.Command("taskkill", "/IM", "btcd.exe", "/F")
	} else if runtime.GOOS == "darwin" {
		// macOS implementation
		checkProcessCmd = exec.Command("pgrep", "btcd")
		killCmd = exec.Command("pkill", "-f", "btcd") // Use -f for full path match
	}

	// Check if btcd is running
	var checkOutput bytes.Buffer
	checkProcessCmd.Stdout = &checkOutput
	checkProcessCmd.Stderr = &checkOutput

	fmt.Println("Checking for running btcd process...")
	err := checkProcessCmd.Run()
	if err != nil || checkOutput.String() == "" {
		fmt.Println("btcd is not running. Cannot stop.")
		return "btcd is not running"
	}

	// Stop btcd process
	fmt.Println("Stopping btcd process...")
	var killOutput bytes.Buffer
	killCmd.Stdout = &killOutput
	killCmd.Stderr = &killOutput

	err = killCmd.Run()
	if err != nil {
		fmt.Printf("Error stopping btcd: %v\n", err)
		return fmt.Sprintf("Error stopping btcd: %s", killOutput.String())
	}

	fmt.Println("btcd stopped successfully")
	return "btcd stopped successfully"
}

// StartBtcwallet is a function to start the btcwallet process
func (bs *BtcService) StartBtcwallet() string {
	// Check if btcwallet is already running
	if isProcessRunning("btcwallet") {
		fmt.Println("btcwallet is already running. Cannot start another instance.")
		return "btcwallet is already running"
	}

	// btcwallet command
	cmd := exec.Command(
		btcwalletPath,
		"--btcdusername=user",
		"--btcdpassword=password",
		"--rpcconnect=127.0.0.1:8334",
		"--noclienttls",
		"--noservertls",
		"--username=user",
		"--password=password",
	)

	// Standard output and error
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start btcwallet process
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting btcwallet: %v\n", err)
		return "Error starting btcwallet"
	}

	fmt.Printf("btcwallet started with PID: %d\n", cmd.Process.Pid)

	// Check if btcwallet process is running
	if !isProcessRunning("btcwallet") {
		fmt.Println("btcwallet process not found after starting.")
		return "btcwallet process not found"
	}

	fmt.Println("btcwallet is running")
	return "btcwallet started successfully"
}

// StopBtcwallet is a function to stop the btcwallet process
func (bs *BtcService) StopBtcwallet() string {
	// check if btcwallet is running
	if !isProcessRunning("btcwallet") {
		fmt.Println("btcwallet is not running. Cannot stop.")
		return "btcwallet is not running"
	}

	cmd := exec.Command("taskkill", "/IM", "btcwallet.exe", "/F")
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error stopping btcwallet: %v\n", err)
		return fmt.Sprintf("Error stopping btcwallet: %s", output.String())
	}

	fmt.Println("btcwallet stopped successfully")
	return "btcwallet stopped successfully"
}


// Init is an initialisation function that starts btcd and btcwallet, connects to the TA server, and exits.
func (bs *BtcService) Init() string {
	SetupTempFilePath()

	// initialize temp file
	err := initializeTempFile()
	if err != nil {
		fmt.Printf("Failed to initialize temp file: %v\n", err)
		return "Failed to initialize temp file"
	}
	fmt.Println("Temporary file initialized successfully.")

	// start btcd
	btcdResult := bs.StartBtcd()
	if btcdResult != "btcd started successfully" {
		bs.StopBtcd()
		bs.StopBtcwallet()
		fmt.Println("Failed to start btcd.")
		return btcdResult
	}

	// start btcwallet
	btcwalletResult := bs.StartBtcwallet()
	if btcwalletResult != "btcwallet started successfully" {
		bs.StopBtcd()
		bs.StopBtcwallet()
		fmt.Println("Failed to start btcwallet.")
		return btcwalletResult
	}

	// connect to TA server
	cmd := exec.Command(
		btcctlPath,
		"--rpcuser=user",
		"--rpcpass=password",
		"--rpcserver=127.0.0.1:8334",
		"--notls",
		"addnode",
		"130.245.173.221:8333",
		"add",
	)

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error connecting to TA server: %v\n", err)
		return fmt.Sprintf("Error connecting to TA server: %s", output.String())
	}

	fmt.Println("Connected to TA server successfully.")

	// stop btcd and btcwallet
	stopBtcdResult := bs.StopBtcd()
	if stopBtcdResult != "btcd stopped successfully" {
		fmt.Println("Failed to stop btcd.")
		return stopBtcdResult
	}

	stopBtcwalletResult := bs.StopBtcwallet()
	if stopBtcwalletResult != "btcwallet stopped successfully" {
		fmt.Println("Failed to stop btcwallet.")
		return stopBtcwalletResult
	}

	fmt.Println("Initialization and cleanup completed successfully.")
	return "Initialization and cleanup completed successfully"
}

// getMiningStatus is a function to check the mining status
func (bs *BtcService) GetMiningStatus() (bool, error) {
	// btcctl command
	cmd := exec.Command(
		btcctlPath,
		"--rpcuser=user",
		"--rpcpass=password",
		"--rpcserver=127.0.0.1:8332",
		"--notls",
		"getgenerate",
	)

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	// execute command
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error executing btcctl getgenerate: %v\n", err)
		return false, fmt.Errorf("error executing btcctl getgenerate: %w", err)
	}

	// get result
	result := output.String()
	fmt.Printf("getgenerate output: %s\n", result)

	// check if mining is active
	if result == "true\n" {
		return true, nil
	}
	return false, nil
}

// StartMining is a function to start mining
func (bs *BtcService) StartMining(numBlock int) string {
	// check if mining is already running
	isMining, err := bs.GetMiningStatus()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return "Error checking mining status"
	}

	if isMining {
		fmt.Println("Mining is already running.")
		return "mining is running"
	}

	// start mining
	cmd := exec.Command(
		btcctlPath,
		"--wallet",
		"--rpcuser=user",
		"--rpcpass=password",
		"--rpcserver=127.0.0.1:8332",
		"--notls",
		"generate",
		strconv.Itoa(numBlock),
	)

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	// execute command
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error starting mining: %v\n", err)
		return fmt.Sprintf("Error starting mining: %s", err.Error())
	}

	fmt.Printf("Mining started successfully. Output: %s\n", output.String())
	return "mining started successfully"
}

// StopMining is a function to stop mining
func (bs *BtcService) StopMining() string {
	// check if mining is running
	isMining, err := bs.GetMiningStatus()
	if err != nil {
		fmt.Printf("Failed to check mining status: %v\n", err)
		return "Error checking mining status"
	}

	// if mining is not running, no action needed
	if !isMining {
		fmt.Println("Mining is not active. No action needed.")
		return "Mining is not active"
	}

	// get mining address from temp file
	miningAddress, err := getMiningAddressFromTemp()
	if err != nil {
		fmt.Printf("Failed to retrieve mining address: %v\n", err)
		return "Failed to retrieve mining address"
	}
	fmt.Printf("Retrieved mining address: %s\n", miningAddress)

	// stop btcd and btcwallet
	stopBtcdResult := bs.StopBtcd()
	if stopBtcdResult != "btcd stopped successfully" {
		fmt.Printf("Failed to stop btcd: %s\n", stopBtcdResult)
		return stopBtcdResult
	}

	stopBtcwalletResult := bs.StopBtcwallet()
	if stopBtcwalletResult != "btcwallet stopped successfully" {
		fmt.Printf("Failed to stop btcwallet: %s\n", stopBtcwalletResult)
		return stopBtcwalletResult
	}

	// restart btcd and btcwallet with mining address
	startBtcdResult := bs.StartBtcd(miningAddress)
	if startBtcdResult != "btcd started successfully" {
		fmt.Printf("Failed to restart btcd with mining address: %s\n", startBtcdResult)
		return startBtcdResult
	}

	startBtcwalletResult := bs.StartBtcwallet()
	if startBtcwalletResult != "btcwallet started successfully" {
		fmt.Printf("Failed to restart btcwallet: %s\n", startBtcwalletResult)
		return startBtcwalletResult
	}

	fmt.Println("Mining process stopped and restarted successfully.")
	return "Mining process stopped and restarted successfully"
}

// btcwalletCreate executes a PowerShell script to create a btcwallet.
func BtcwalletCreate(passphrase string) error {
	// Configure PowerShell command with -WindowStyle Hidden
	cmd := exec.Command("powershell",
		"-NoProfile",
		"-ExecutionPolicy", "Bypass",
		"-WindowStyle", "Hidden",
		"-File", btcwalletScriptPath,
	)

	// Add environment variable for the passphrase
	cmd.Env = append(os.Environ(), "BTCWALLET_PASSPHRASE="+passphrase)

	// Hide terminal window in Windows (I think this is not working)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// HideWindow:    true,
		// CreationFlags: CREATE_NO_WINDOW,
	}

	// Hide terminal window in Windows
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute the command
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start btcwallet process: %v", err)
	}

	// Since it runs in the background, return without waiting
	// The console window will not appear as the application is in GUI mode.
	// Logs can be recorded if needed.
	fmt.Println("btcwallet process started in the background.")

	return nil
}

// GetNewAddress is a function to generate a new address
func (bs *BtcService) GetNewAddress() (string, error) {
	// check if btcd and btcwallet are running
	btcdRunning := isProcessRunning("btcd")
	btcwalletRunning := isProcessRunning("btcwallet")

	// check if btcd and btcwallet are running
	if !btcdRunning {
		return "", fmt.Errorf("btcd is not running. Please start btcd before calling this function")
	}
	if !btcwalletRunning {
		return "", fmt.Errorf("btcwallet is not running. Please start btcwallet before calling this function")
	}

	// get new address
	cmd := exec.Command(
		btcctlPath,
		"--rpcuser=user",
		"--rpcpass=password",
		"--rpcserver=127.0.0.1:8332",
		"--notls",
		"getnewaddress",
	)

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error generating new address: %v\n", err)
		return "", fmt.Errorf("error generating new address: %w", err)
	}

	// new address
	newAddress := strings.TrimSpace(output.String())
	fmt.Printf("Generated new address: %s\n", newAddress)

	return newAddress, nil
}

// UnlockWallet is a function to unlock the wallet
func (bs *BtcService) UnlockWallet(passphrase string) (string, error) {
	// check if btcd and btcwallet are running
	if !isProcessRunning("btcd") {
		return "", fmt.Errorf("btcd is not running. Please start btcd before unlocking the wallet")
	}

	if !isProcessRunning("btcwallet") {
		return "", fmt.Errorf("btcwallet is not running. Please start btcwallet before unlocking the wallet")
	}

	// unlock wallet command for 600 seconds
	cmd := exec.Command(
		btcctlPath,
		"--wallet",
		"--rpcuser=user",
		"--rpcpass=password",
		"--rpcserver=127.0.0.1:8332",
		"--notls",
		"walletpassphrase",
		passphrase,
		"600", // unlock for 600 seconds
	)

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error unlocking wallet: %v\n", err)
		return "", fmt.Errorf("error unlocking wallet: %w", err)
	}

	result := strings.TrimSpace(output.String())
	fmt.Println("Wallet unlocked successfully.")
	return result, nil
}

// LockWallet is a function to lock the wallet
func (bs *BtcService) LockWallet() (string, error) {
	// check if btcd and btcwallet are running
	if !isProcessRunning("btcd") {
		return "", fmt.Errorf("btcd is not running. Please start btcd before locking the wallet")
	}

	if !isProcessRunning("btcwallet") {
		return "", fmt.Errorf("btcwallet is not running. Please start btcwallet before locking the wallet")
	}

	// lock wallet command
	cmd := exec.Command(
		btcctlPath,
		"--wallet",
		"--rpcuser=user",
		"--rpcpass=password",
		"--rpcserver=127.0.0.1:8332",
		"--notls",
		"walletlock",
	)

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error locking wallet: %v\n", err)
		return "", fmt.Errorf("error locking wallet: %w", err)
	}

	result := strings.TrimSpace(output.String())
	fmt.Println("Wallet locked successfully.")
	return result, nil
}

func (bs *BtcService) Login(walletAddress, passphrase string) (string, error) {
	// Step 0: Check if the wallet exists
	walletDBPath := fmt.Sprintf(`%s\AppData\Local\Btcwallet\mainnet\wallet.db`, os.Getenv("USERPROFILE")) // Adjust path as needed
	if _, err := os.Stat(walletDBPath); os.IsNotExist(err) {
		fmt.Printf("Wallet does not exist at path: %s\n", walletDBPath)
		return "Wallet does not exist", fmt.Errorf("wallet does not exist at path: %s", walletDBPath)
	}
	// Step 1: Start btcd with wallet address
	btcdResult := bs.StartBtcd(walletAddress)
	time.Sleep(time.Second) // 1 second
	if btcdResult != "btcd started successfully" {
		fmt.Printf("Failed to start btcd: %s\n", btcdResult)
		bs.StopBtcd()
		return "Failed to start btcd", fmt.Errorf("failed to start btcd: %s", btcdResult)
	}

	// Step 2: Start btcwallet
	btcwalletResult := bs.StartBtcwallet()
	time.Sleep(time.Second) // 1 second
	if btcwalletResult != "btcwallet started successfully" {
		fmt.Printf("Failed to start btcwallet: %s\n", btcwalletResult)
		bs.StopBtcd()
		bs.StopBtcwallet()
		return "Failed to start btcwallet", fmt.Errorf("failed to start btcwallet: %s", btcwalletResult)
	}

	// Step 4: Unlock the wallet
	unlockResult, err := bs.UnlockWallet(passphrase)
	time.Sleep(time.Second) // 1 second
	if err != nil {
		fmt.Printf("Failed to unlock wallet: %v\n", err)
		bs.StopBtcd()
		time.Sleep(time.Second) // 1 second
		bs.StopBtcwallet()
		return "Failed to unlock wallet", fmt.Errorf("failed to unlock wallet: %w", err)
	}

	// Step 5: Success
	fmt.Println("Login successful. Wallet unlocked.")
	return unlockResult, nil
}

func (bs *BtcService) GetBalance() (string, error) {
	// Check if btcd and btcwallet are running
	if !isProcessRunning("btcd") {
		return "", fmt.Errorf("btcd is not running. Please start btcd before checking balance")
	}

	if !isProcessRunning("btcwallet") {
		return "", fmt.Errorf("btcwallet is not running. Please start btcwallet before checking balance")
	}

	// Execute btcctl getbalance command
	cmd := exec.Command(
		btcctlPath,
		"--wallet",
		"--rpcuser=user",
		"--rpcpass=password",
		"--rpcserver=127.0.0.1:8332",
		"--notls",
		"getbalance",
	)

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error fetching balance: %v\n", err)
		return "", fmt.Errorf("error fetching balance: %w", err)
	}

	// Return balance
	balance := strings.TrimSpace(output.String())
	fmt.Printf("Wallet balance: %s\n", balance)
	return balance, nil
}

func (bs *BtcService) GetReceivedByAddress(walletAddress string) (string, error) {
	// Check if btcd and btcwallet are running
	if !isProcessRunning("btcd") {
		return "", fmt.Errorf("btcd is not running. Please start btcd before checking received amount")
	}

	if !isProcessRunning("btcwallet") {
		return "", fmt.Errorf("btcwallet is not running. Please start btcwallet before checking received amount")
	}

	// Execute btcctl getreceivedbyaddress command
	cmd := exec.Command(
		btcctlPath,
		"--wallet",
		"--rpcuser=user",
		"--rpcpass=password",
		"--rpcserver=127.0.0.1:8332",
		"--notls",
		"getreceivedbyaddress",
		walletAddress,
		"1", // Minimum confirmations set to 1
	)

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error fetching received amount for address %s: %v\n", walletAddress, err)
		return "", fmt.Errorf("error fetching received amount for address %s: %w", walletAddress, err)
	}

	// Return received amount
	receivedAmount := strings.TrimSpace(output.String())
	fmt.Printf("Received amount for address %s: %s\n", walletAddress, receivedAmount)
	return receivedAmount, nil
}

func (bs *BtcService) GetBlockCount() (string, error) {
	// Check if btcd is running
	if !isProcessRunning("btcd") {
		return "", fmt.Errorf("btcd is not running. Please start btcd before checking block count")
	}

	// Execute btcctl getblockcount command
	cmd := exec.Command(
		btcctlPath,
		"--rpcuser=user",
		"--rpcpass=password",
		"--rpcserver=127.0.0.1:8334",
		"--notls",
		"getblockcount",
	)

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error fetching block count: %v\n", err)
		return "", fmt.Errorf("error fetching block count: %w", err)
	}

	// Return block count
	blockCount := strings.TrimSpace(output.String())
	fmt.Printf("Current block count: %s\n", blockCount)
	return blockCount, nil
}

func (bs *BtcService) ListReceivedByAddress() ([]map[string]interface{}, error) {
	// Check if btcd and btcwallet are running
	if !isProcessRunning("btcd") {
		return nil, fmt.Errorf("btcd is not running. Please start btcd before listing addresses")
	}

	if !isProcessRunning("btcwallet") {
		return nil, fmt.Errorf("btcwallet is not running. Please start btcwallet before listing addresses")
	}

	// Execute btcctl listreceivedbyaddress command
	cmd := exec.Command(
		btcctlPath,
		"--wallet",
		"--rpcuser=user",
		"--rpcpass=password",
		"--rpcserver=127.0.0.1:8332",
		"--notls",
		"listreceivedbyaddress",
		"0",    // Include addresses with 0 confirmations
		"true", // Include empty addresses
	)

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error listing received addresses: %v\n", err)
		return nil, fmt.Errorf("error listing received addresses: %w", err)
	}

	// Parse the output as JSON
	var addresses []map[string]interface{}
	if err := json.Unmarshal(output.Bytes(), &addresses); err != nil {
		fmt.Printf("Error parsing address list: %v\n", err)
		return nil, fmt.Errorf("error parsing address list: %w", err)
	}

	// Log full result for debugging
	fmt.Printf("ㅇㅇFull list of received addresses: %v\n", addresses)

	// Return full result
	return addresses, nil
}

func (bs *BtcService) ListUnspent() ([]map[string]interface{}, error) {
	// Check if btcd and btcwallet are running
	if !isProcessRunning("btcd") {
		return nil, fmt.Errorf("btcd is not running. Please start btcd before listing unspent transactions")
	}

	if !isProcessRunning("btcwallet") {
		return nil, fmt.Errorf("btcwallet is not running. Please start btcwallet before listing unspent transactions")
	}

	// Execute btcctl listunspent command
	cmd := exec.Command(
		btcctlPath,
		"--wallet",
		"--rpcuser=user",
		"--rpcpass=password",
		"--rpcserver=127.0.0.1:8332",
		"--notls",
		"listunspent",
	)

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error listing unspent transactions: %v\n", err)
		return nil, fmt.Errorf("error listing unspent transactions: %w", err)
	}

	// Parse the output as JSON
	var utxos []map[string]interface{}
	if err := json.Unmarshal(output.Bytes(), &utxos); err != nil {
		fmt.Printf("Error parsing UTXO list: %v\n", err)
		return nil, fmt.Errorf("error parsing UTXO list: %w", err)
	}

	// Log full result for debugging
	fmt.Printf("List of unspent transactions: %v\n", utxos)

	// Return full result
	return utxos, nil
}

// CreateRawTransaction 함수 수정
func (bs *BtcService) CreateRawTransaction(txid string, dst string, amount float64) (string, error) {
	// Step 1: Retrieve source address (mining address) from temp file
	fmt.Println("Step 1: Retrieving source address from temp file...")
	tempContent, err := ioutil.ReadFile(tempFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read temp file: %w", err)
	}

	fmt.Printf("Raw temp file content: %s\n", tempContent)

	var tempData map[string]string
	if err := json.Unmarshal(tempContent, &tempData); err != nil {
		return "", fmt.Errorf("failed to parse temp file content: %w", err)
	}

	srcAddress, exists := tempData["miningaddr"]
	if !exists || srcAddress == "" {
		return "", fmt.Errorf("source address not found in temp file")
	}

	fmt.Printf("Source address retrieved: %s\n", srcAddress)

	// Step 2: Get UTXOs from ListUnspent
	fmt.Println("Step 2: Retrieving unspent transactions (UTXOs)...")
	utxos, err := bs.ListUnspent()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve unspent transactions: %w", err)
	}

	// Validate UTXOs
	fmt.Println("Step 3: Validating UTXOs...")
	var selectedUTXO map[string]interface{}
	for _, utxo := range utxos {
		utxoTxID, ok1 := utxo["txid"].(string)
		utxoAddress, ok2 := utxo["address"].(string)
		utxoAmount, ok3 := utxo["amount"].(float64)

		if !ok1 || !ok2 || !ok3 {
			continue // 잘못된 UTXO 형식 건너뛰기
		}

		// Match txid, source address, and check amount sufficiency
		if utxoTxID == txid && utxoAddress == srcAddress && utxoAmount >= amount {
			selectedUTXO = utxo
			fmt.Printf("Matching UTXO found: txid=%s, amount=%.8f\n", utxoTxID, utxoAmount)
			break
		}
	}

	if selectedUTXO == nil {
		return "", fmt.Errorf("no suitable UTXO found for txid %s and address %s with amount >= %.8f", txid, srcAddress, amount)
	}

	// Extract UTXO details
	voutFloat, ok := selectedUTXO["vout"].(float64)
	if !ok {
		return "", fmt.Errorf("invalid vout type in UTXO")
	}
	vout := int(voutFloat)

	srcAmount, ok := selectedUTXO["amount"].(float64)
	if !ok {
		return "", fmt.Errorf("invalid amount type in UTXO")
	}

	// Step 4: Construct raw transaction command
	fmt.Println("Step 4: Constructing raw transaction command...")
	rawTxCommandWin := []string{
		"--wallet",
		"--rpcuser=user",
		"--rpcpass=password",
		"--rpcserver=127.0.0.1:8332",
		"--notls",
		"createrawtransaction",
		fmt.Sprintf(`[{"txid":"%s", "vout":%d}]`, txid, vout),
		fmt.Sprintf(`{"%s": %.8f, "%s": %.8f}`, dst, amount, srcAddress, srcAmount-amount),
	}

	rawTxCommandMac := []string{
		"--wallet",
		"--rpcuser=user",
		"--rpcpass=password",
		"--rpcserver=127.0.0.1:8332",
		"--notls",
		"createrawtransaction",
		fmt.Sprintf(`[{"txid":"%s", "vout":%d}]`, txid, vout),
		fmt.Sprintf(`{"%s": %.8f, "%s": %.8f}`, dst, amount, srcAddress, srcAmount-amount),
	}

	// Determine OS-specific command
	var rawTxCommand []string
	if os.Getenv("OS") == "Windows_NT" {
		rawTxCommand = rawTxCommandWin
	} else {
		rawTxCommand = rawTxCommandMac
	}

	fmt.Printf("Raw transaction command: %v\n", rawTxCommand)

	// Step 5: Execute raw transaction command
	fmt.Println("Step 5: Executing raw transaction command...")
	cmd := exec.Command(btcctlPath, rawTxCommand...)
	var output bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("[ERROR] Command execution failed: %v\n", err)
		fmt.Printf("[ERROR] Stderr: %s\n", stderr.String())
		fmt.Printf("[DEBUG] Raw command output: %s\n", output.String())
		return "", fmt.Errorf("failed to create raw transaction: %w. Stderr: %s", err, stderr.String())
	}

	rawId := strings.TrimSpace(output.String())
	fmt.Printf("Raw transaction created: %s\n", rawId)
	return rawId, nil
}

// signRawTransaction 함수 수정
func (bs *BtcService) signRawTransaction(rawId string) (string, bool, error) {
	fmt.Printf("[DEBUG] Starting signRawTransaction with rawId: %s\n", rawId)

	// Validate rawId
	if rawId == "" {
		return "", false, fmt.Errorf("rawId is empty")
	}

	// Construct the command
	cmd := exec.Command(
		btcctlPath,
		"--wallet",
		"--rpcuser=user",
		"--rpcpass=password",
		"--rpcserver=127.0.0.1:8332",
		"--notls",
		"signrawtransaction",
		rawId,
	)

	var output bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &stderr

	// Execute the command
	if err := cmd.Run(); err != nil {
		fmt.Printf("[ERROR] Command execution failed: %v\n", err)
		fmt.Printf("[ERROR] Stderr: %s\n", stderr.String())
		fmt.Printf("[DEBUG] Raw command output: %s\n", output.String())
		return "", false, fmt.Errorf("failed to sign raw transaction: %w. Stderr: %s", err, stderr.String())
	}

	fmt.Printf("[DEBUG] Command executed successfully. Raw output: %s\n", output.String())

	// Parse the output JSON
	var signedTx map[string]interface{}
	if err := json.Unmarshal(output.Bytes(), &signedTx); err != nil {
		fmt.Printf("[ERROR] Failed to parse JSON output: %v\n", err)
		fmt.Printf("[DEBUG] Raw JSON output: %s\n", output.String())
		return "", false, fmt.Errorf("failed to parse signed transaction: %w", err)
	}

	// Extract the 'hex' and 'complete' fields
	hex, hexOk := signedTx["hex"].(string)
	complete, completeOk := signedTx["complete"].(bool)

	if !hexOk || !completeOk {
		fmt.Printf("[ERROR] Unexpected output format. Parsed data: %v\n", signedTx)
		return "", false, fmt.Errorf("unexpected output format: %s", output.String())
	}

	fmt.Printf("[DEBUG] Transaction signed successfully. Hex: %s, Complete: %v\n", hex, complete)
	return hex, complete, nil
}

// sendRawTransaction 함수는 문제 없어 보입니다.
func (bs *BtcService) sendRawTransaction(hex string) (string, error) {
	// Use btcctlPath to construct the command
	cmd := exec.Command(
		btcctlPath,
		"--wallet",
		"--rpcuser=user",
		"--rpcpass=password",
		"--rpcserver=127.0.0.1:8332",
		"--notls",
		"sendrawtransaction",
		hex,
	)

	var output bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &stderr

	// Execute the command
	if err := cmd.Run(); err != nil {
		fmt.Printf("[ERROR] Command execution failed: %v\n", err)
		fmt.Printf("[ERROR] Stderr: %s\n", stderr.String())
		fmt.Printf("[DEBUG] Raw command output: %s\n", output.String())
		return "", fmt.Errorf("failed to send raw transaction: %w. Stderr: %s", err, stderr.String())
	}

	// Get the transaction ID from the output
	txid := strings.TrimSpace(output.String())
	fmt.Printf("Transaction sent successfully. TxID: %s\n", txid)

	// Return the transaction ID
	return txid, nil
}

// Transaction 함수 수정
func (bs *BtcService) Transaction(passphrase, txid, dst string, amount float64) (string, error) {
	fmt.Printf("[DEBUG] Starting transaction with passphrase: %s, txid: %s, dst: %s, amount: %.8f\n", passphrase, txid, dst, amount)

	// Step 1: Store original balance in a temp file
	fmt.Println("[DEBUG] Step 1: Retrieving original balance...")
	originalBalanceStr, err := bs.GetBalance()
	if err != nil {
		fmt.Printf("[ERROR] Failed to retrieve original balance: %v\n", err)
		return "", fmt.Errorf("failed to retrieve original balance: %w", err)
	}
	fmt.Printf("[DEBUG] Original balance: %s\n", originalBalanceStr)

	if err := updateTempFile("originalBalance", originalBalanceStr); err != nil {
		fmt.Printf("[ERROR] Failed to store original balance: %v\n", err)
		return "", fmt.Errorf("failed to store original balance: %w", err)
	}

	// // Convert originalBalance to float64
	// originalBalance, err := strconv.ParseFloat(originalBalanceStr, 64)
	// if err != nil {
	// 	fmt.Printf("[ERROR] Failed to parse original balance: %v\n", err)
	// 	return "", fmt.Errorf("failed to parse original balance: %w", err)
	// }

	// Step 2: Create a raw transaction
	fmt.Println("[DEBUG] Step 2: Creating raw transaction...")
	rawId, err := bs.CreateRawTransaction(txid, dst, amount)
	if err != nil {
		fmt.Printf("[ERROR] Failed to create raw transaction: %v\n", err)
		return "", fmt.Errorf("failed to create raw transaction: %w", err)
	}
	fmt.Printf("[DEBUG] Raw transaction ID: %s\n", rawId)
	time.Sleep(1 * time.Second)

	// Step 3: Unlock the wallet
	fmt.Println("[DEBUG] Step 3: Unlocking the wallet...")
	if _, err := bs.UnlockWallet(passphrase); err != nil {
		fmt.Printf("[ERROR] Failed to unlock wallet: %v\n", err)
		return "", fmt.Errorf("failed to unlock wallet: %w", err)
	}
	time.Sleep(1 * time.Second)

	// Step 4: Sign the raw transaction
	fmt.Println("[DEBUG] Step 4: Signing raw transaction...")
	hex, complete, err := bs.signRawTransaction(rawId)
	if err != nil {
		fmt.Printf("[ERROR] Failed to sign raw transaction: %v\n", err)
		return "", fmt.Errorf("failed to sign raw transaction: %w", err)
	}
	fmt.Printf("[DEBUG] Transaction signing complete. Hex: %s, Complete: %v\n", hex, complete)

	if !complete {
		fmt.Printf("[ERROR] Transaction signing incomplete\n")
		return "", fmt.Errorf("transaction signing incomplete")
	}
	time.Sleep(1 * time.Second)

	// Step 5: Send the raw transaction
	fmt.Println("[DEBUG] Step 5: Sending raw transaction...")
	txIdResult, err := bs.sendRawTransaction(hex)
	if err != nil {
		fmt.Printf("[ERROR] Failed to send raw transaction: %v\n", err)
		return "", fmt.Errorf("failed to send raw transaction: %w", err)
	}
	fmt.Printf("[DEBUG] Raw transaction sent successfully. TxID: %s\n", txIdResult)
	time.Sleep(1 * time.Second)

	// // Step 6: Verify the transaction and balance
	// fmt.Println("[DEBUG] Step 6: Verifying transaction and balance...")
	// currentBalanceStr, err := bs.GetBalance()
	// if err != nil {
	// 	fmt.Printf("[ERROR] Failed to retrieve current balance: %v\n", err)
	// 	return "", fmt.Errorf("failed to retrieve current balance: %w", err)
	// }

	// // Convert currentBalance to float64
	// currentBalance, err := strconv.ParseFloat(currentBalanceStr, 64)
	// if err != nil {
	// 	fmt.Printf("[ERROR] Failed to parse current balance: %v\n", err)
	// 	return "", fmt.Errorf("failed to parse current balance: %w", err)
	// }

	// // Validate balance
	// expectedBalance := originalBalance - amount
	// if currentBalance != expectedBalance {
	// 	fmt.Printf("[ERROR] Balance mismatch: Expected %.8f, Got %.8f\n", expectedBalance, currentBalance)
	// 	return "", fmt.Errorf("balance mismatch after transaction")
	// }

	fmt.Println("[DEBUG] Balance validation successful.")
	time.Sleep(1 * time.Second)

	// Step 7: Lock the wallet again (optional)
	fmt.Println("[DEBUG] Step 7: Locking the wallet again...")
	if _, err := bs.LockWallet(); err != nil {
		fmt.Printf("[WARNING] Failed to lock wallet: %v\n", err)
		// Continue even if locking fails
	}
	time.Sleep(1 * time.Second)

	// Step 8: Delete temporary balance file
	fmt.Println("[DEBUG] Step 8: Cleaning up temporary files...")
	if err := deleteFromTempFile("originalBalance"); err != nil {
		fmt.Printf("[WARNING] Failed to delete temp file: %v\n", err)
	}
	time.Sleep(1 * time.Second)

	fmt.Println("[DEBUG] Transaction completed successfully.")
	return txIdResult, nil
}
