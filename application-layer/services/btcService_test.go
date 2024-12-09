package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// Get-Process | Where-Object {$_.Name -like "btc*"}

// C:\dev\workspace\CSE-416\application-layer\services> go test -v -run ^TestStopBtcd$
// go test -v -run ^TestStartBtcdWithNoArgs$
// go test -v -run ^TestStartBtcdWithNoArgs$ -count=1 application-layer/services
// TestStartBtcdWithNoArgs validates starting btcd without arguments.
func TestStartBtcdWithNoArgs(t *testing.T) {
	btcService := NewBtcService()

	// Log test start
	t.Log("Starting TestStartBtcdWithNoArgs...")

	// Call StartBtcd with no arguments
	result := btcService.StartBtcd()
	expectedSuccess := "btcd started successfully"
	expectedFailure := "Error starting btcd"
	expectedAlreadyRunning := "btcd is already running"

	// Log result for debugging
	t.Logf("Test Result: %s", result)

	// Validate results based on OS
	if result != expectedSuccess && result != expectedFailure && result != expectedAlreadyRunning {
		t.Errorf("Unexpected result: %q", result)
	}

	// Additional logging based on results
	if result == expectedSuccess {
		t.Log("btcd started successfully with no arguments.")
	} else if result == expectedFailure {
		t.Log("btcd failed to start.")
	} else if result == expectedAlreadyRunning {
		t.Log("btcd is already running.")
	}

	t.Log("TestStartBtcdWithNoArgs completed.")
}

// go test -v -run ^TestStartBtcdWithNoArgsAndStop$ -count=1 application-layer/services
// TestStartBtcdWithNoArgsAndStop validates starting and stopping btcd without arguments.
func TestStartBtcdWithNoArgsAndStop(t *testing.T) {
	btcService := NewBtcService()

	// Log test start
	t.Log("Starting TestStartBtcdWithNoArgsAndStop...")

	// Call StartBtcd with no arguments
	result := btcService.StartBtcd()
	expectedSuccess := "btcd started successfully"
	expectedFailure := "Error starting btcd"
	expectedAlreadyRunning := "btcd is already running"

	// Log result for debugging
	t.Logf("Test Result: %s", result)

	// Validate results based on OS
	if result != expectedSuccess && result != expectedFailure && result != expectedAlreadyRunning {
		t.Errorf("Unexpected result: %q", result)
	}

	// Additional logging based on results
	if result == expectedSuccess {
		t.Log("btcd started successfully with no arguments.")
	} else if result == expectedFailure {
		t.Log("btcd failed to start.")
	} else if result == expectedAlreadyRunning {
		t.Log("btcd is already running.")
	}

	// Stop btcd process (teardown)
	t.Log("Stopping btcd process...")
	stopResult := btcService.StopBtcd()
	t.Logf("Stop Result: %s", stopResult)

	// Validate stop results
	expectedStopSuccess := "btcd stopped successfully"
	expectedNotRunning := "btcd is not running"

	if stopResult != expectedStopSuccess && stopResult != expectedNotRunning {
		t.Errorf("Unexpected stop result: %q", stopResult)
	} else if stopResult == expectedStopSuccess {
		t.Log("btcd stopped successfully.")
	} else if stopResult == expectedNotRunning {
		t.Log("btcd was not running.")
	}

	t.Log("TestStartBtcdWithNoArgsAndStop completed.")
}

// C:\dev\workspace\CSE-416\application-layer\services> go test -v -run ^TestStopBtcd$
// TestStopBtcd validates stopping the btcd process.
// go test -v -run ^TestStopBtcd$ -count=1 application-layer/services
func TestStopBtcd(t *testing.T) {
	btcService := NewBtcService()

	// Log test start
	t.Log("Starting TestStopBtcd...")

	// Call StopBtcd
	result := btcService.StopBtcd()

	// Determine OS-specific expectations
	expectedNotRunning := "btcd is not running"
	expectedStopped := "btcd stopped successfully"

	// Validate results based on OS
	if runtime.GOOS == "windows" {
		t.Log("Testing on Windows...")
		if result != expectedNotRunning && result != expectedStopped {
			t.Errorf("Unexpected result on Windows: %q", result)
		}
	} else if runtime.GOOS == "darwin" {
		t.Log("Testing on macOS...")
		if result != expectedNotRunning && result != expectedStopped {
			t.Errorf("Unexpected result on macOS: %q", result)
		}
	}

	// Log result for debugging
	t.Logf("Result: %s", result)

	t.Log("TestStopBtcd completed.")
}

// go test -v -run ^TestBtcwalletCreate$ -count=1 application-layer/services
// TestBtcwalletCreate btcwallet implementation test for macOs and Windows Powershell
func TestBtcwalletCreate(t *testing.T) {
	var walletDBPath string

	btcService := NewBtcService()

	// Determine OS-specific wallet.db path
	if runtime.GOOS == "windows" {
		userProfile := os.Getenv("USERPROFILE")
		if userProfile == "" {
			t.Fatal("USERPROFILE environment variable is not set")
		}
		walletDBPath = filepath.Join(userProfile, "AppData", "Local", "Btcwallet", "mainnet", "wallet.db")
	} else if runtime.GOOS == "darwin" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			t.Fatalf("Failed to get home directory: %v", err)
		}
		walletDBPath = filepath.Join(homeDir, "Library", "Application Support", "Btcwallet", "mainnet", "wallet.db")
	} else {
		t.Fatalf("Unsupported OS: %s", runtime.GOOS)
	}

	// Remove existing wallet.db for a clean test
	err := os.Remove(walletDBPath)
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("Failed to delete existing wallet.db: %v", err)
	}

	// Call BtcwalletCreate to test wallet creation
	passphrase := "CSE416"
	err = btcService.BtcwalletCreate(passphrase)
	if err != nil {
		t.Fatalf("Failed to create btcwallet: %v", err)
	}

	// Verify wallet.db was created
	if _, err := os.Stat(walletDBPath); os.IsNotExist(err) {
		t.Fatal("wallet.db does not exist after creation")
	} else {
		fmt.Println("wallet.db successfully created.")
	}
}

// go test -v -run ^TestStartBtcwallet$ -count=1 application-layer/services
// TestStartBtcwallet validates starting the btcwallet process.
func TestStartBtcwallet(t *testing.T) {
	btcService := NewBtcService()

	t.Log("Starting TestStartBtcwallet...")

	// attempt to start btcwallet
	result := btcService.StartBtcwallet()
	expectedSuccess := "btcwallet started successfully"
	expectedAlreadyRunning := "btcwallet is already running"

	// result validation
	if result != expectedSuccess && result != expectedAlreadyRunning {
		t.Errorf("Unexpected result: %q", result)
	}

	// debugging logs
	if result == expectedSuccess {
		t.Log("btcwallet started successfully.")
	} else if result == expectedAlreadyRunning {
		t.Log("btcwallet is already running.")
	}

	t.Log("TestStartBtcwallet completed.")
}

// go test -v -run ^TestStopBtcwallet$ -count=1 application-layer/services
// TestStopBtcwallet validates stopping the btcwallet process.
func TestStopBtcwallet(t *testing.T) {
	btcService := NewBtcService()

	t.Log("Starting TestStopBtcwallet...")

	// attempt to stop btcwallet
	result := btcService.StopBtcwallet()
	expectedNotRunning := "btcwallet is not running"
	expectedStopped := "btcwallet stopped successfully"

	// result validation
	if result != expectedNotRunning && result != expectedStopped {
		t.Errorf("Unexpected result: %q", result)
	}

	// debugging logs
	if result == expectedNotRunning {
		t.Log("btcwallet is not running.")
	} else if result == expectedStopped {
		t.Log("btcwallet stopped successfully.")
	}

	t.Log("TestStopBtcwallet completed.")
}

// go test -v -run ^TestCreateWallet$ -count=1 application-layer/services
// TestCreateWallet validates the workflow of creating a wallet, generating a new address, and cleaning up resources.
// TestCreateWallet validates the creation of a new wallet, address generation, and proper cleanup of btcd and btcwallet processes.
func TestCreateWallet(t *testing.T) {
	btcService := NewBtcService()

	// Define test passphrase
	passphrase := "CSE416"

	// Call CreateWallet
	newAddress, err := btcService.CreateWallet(passphrase)
	if err != nil {
		t.Fatalf("CreateWallet failed: %v", err)
	}

	// Validate the generated address
	if newAddress == "" {
		t.Errorf("Expected a valid new address, but got an empty string")
	} else {
		t.Logf("Generated new address: %s", newAddress)
	}

	// Ensure btcd and btcwallet are not running after function execution
	if isProcessRunning("btcd") {
		t.Errorf("btcd process is still running after CreateWallet execution")
	}
	if isProcessRunning("btcwallet") {
		t.Errorf("btcwallet process is still running after CreateWallet execution")
	}

	// Log success
	t.Log("CreateWallet executed successfully, and processes were cleaned up.")
}


// go test -v -run ^TestInit$ -count=1 application-layer/services
// TestInit validates the initialization of the BtcService.
func TestInit(t *testing.T) {
	btcService := NewBtcService()

	// call Init function
	result := btcService.Init()

	// Small delay to allow process state updates
	time.Sleep(500 * time.Millisecond)

	expected := "Initialization and cleanup completed successfully"

	// results validation
	if result != expected {
		t.Errorf("Expected %q, but got %q", expected, result)
	}

	// debugging log
	t.Logf("Init function result: %s", result)
}

// go test -v -run ^TestUnlockWallet$ -count=1 application-layer/services
// TestUnlockWallet validates unlocking the wallet.
func TestUnlockWallet(t *testing.T) {
	btcService := NewBtcService()

	// Test passphrase
	passphrase := "CSE416"

	// Ensure btcd and btcwallet are running
	if !isProcessRunning("btcd") || !isProcessRunning("btcwallet") {
		t.Fatalf("btcd or btcwallet is not running. Please start both processes before testing")
	}

	// Call UnlockWallet
	result, err := btcService.UnlockWallet(passphrase)
	if err != nil {
		t.Errorf("Failed to unlock wallet: %v", err)
	}

	// Verify the result
	expected := "" // btcctl command typically outputs an empty string
	if result != expected {
		t.Errorf("Expected %q, but got %q", expected, result)
	}

	t.Logf("UnlockWallet Result: %s", result)
}

// go test -v -run ^TestLockWallet$ -count=1 application-layer/services
// TestLockWallet validates locking the wallet.
func TestLockWallet(t *testing.T) {
	btcService := NewBtcService()

	// Ensure btcd and btcwallet are running
	if !isProcessRunning("btcd") || !isProcessRunning("btcwallet") {
		t.Fatalf("btcd or btcwallet is not running. Please start both processes before testing")
	}

	// Call LockWallet
	result, err := btcService.LockWallet()
	if err != nil {
		t.Errorf("Failed to lock wallet: %v", err)
	}

	// Verify the result
	expected := "" // btcctl command typically outputs an empty string
	if result != expected {
		t.Errorf("Expected %q, but got %q", expected, result)
	}

	t.Logf("LockWallet Result: %s", result)
}

// go test -v -run ^TestGetNewAddress$ -count=1 application-layer/services
// TestGetNewAddress validates getting a new address.
func TestGetNewAddress(t *testing.T) {
	btcService := NewBtcService()

	// Ensure btcd and btcwallet are running
	if !isProcessRunning("btcd") || !isProcessRunning("btcwallet") {
		t.Fatalf("btcd or btcwallet is not running. Please start both processes before testing")
	}

	// Call GetNewAddress
	newAddress, err := btcService.GetNewAddress()
	if err != nil {
		t.Fatalf("Failed to get new address: %v", err)
	}

	// Validate the generated address
	if newAddress == "" {
		t.Fatalf("Generated address is empty")
	}

	// Log and print the result
	t.Logf("Generated new address: %s", newAddress)
	fmt.Printf("Generated new address (from test): %s\n", newAddress)
}

// go test -v -run ^TestListReceivedByAddress$ -count=1 application-layer/services
// TestListReceivedByAddress validates listing received addresses.
func TestListReceivedByAddress(t *testing.T) {
	btcService := NewBtcService()

	// Ensure btcd and btcwallet are running
	if !isProcessRunning("btcd") || !isProcessRunning("btcwallet") {
		t.Fatalf("btcd or btcwallet is not running. Please start both processes before testing")
	}

	// Call ListReceivedByAddress
	addressList, err := btcService.ListReceivedByAddress()
	if err != nil {
		t.Errorf("Failed to list received addresses: %v", err)
		return
	}

	// Validate that addresses are returned
	if len(addressList) == 0 {
		t.Fatalf("No addresses received from ListReceivedByAddress")
	}

	// Log and print the result
	t.Logf("Received Addresses: %v", addressList)
	fmt.Printf("Received Addresses (from test): %v\n", addressList)
}

// go test -v -run ^TestStartBtcdWithArgs$ -count=1 application-layer/services
// TestStartBtcdWithArgs validates starting btcd with arguments.
func TestStartBtcdWithArgs(t *testing.T) {
	btcService := NewBtcService()

	// 초기화 호출
	SetupTempFilePath()

	result := btcService.StartBtcd("14QnrKvCS9cskoMjfKkCe7xaWkQwdWCbJc")
	expectedSuccess := "btcd started successfully"
	expectedFailure := "Error starting btcd"
	expectedAlreadyRunning := "btcd is already running"

	if result != expectedSuccess && result != expectedFailure && result != expectedAlreadyRunning {
		t.Errorf("Unexpected result: %q", result)
	}

	if result == expectedSuccess {
		t.Log("btcd started successfully with wallet address.")
	}
}

// go test -v -run ^TestStartBtcdWithInvalidArgs$
// TestStartBtcdWithInvalidArgs validates behavior for invalid arguments.
func TestStartBtcdWithInvalidArgs(t *testing.T) {
	btcService := NewBtcService()

	// Call StartBtcd with invalid arguments
	result := btcService.StartBtcd("1ExampleWalletAddress", "AnotherArgument")
	expected := "Invalid number of arguments"

	// Validate result
	if result != expected {
		t.Errorf("Expected %q, but got %q", expected, result)
	}
}

// go test -v -run ^TestGetMiningStatus$ -count=1 application-layer/services
// TestGetMiningStatus validates checking the mining status.
func TestGetMiningStatus(t *testing.T) {
	btcService := NewBtcService()

	// Check mining status
	status, err := btcService.GetMiningStatus()
	if err != nil {
		t.Errorf("Error checking mining status: %v", err)
	}

	// Log the mining status
	t.Logf("Mining status: %t", status)
}

// go test -v -run ^TestStartMining$ -count=1 application-layer/services
// TestStartMining validates starting the mining process.
func TestStartMining(t *testing.T) {
	btcService := NewBtcService()

	result := btcService.StartMining(5) // Request to generate 5 blocks

	// Expected outcomes
	expected := "mining started successfully"
	expectedAlreadyRunning := "mining is running"
	expectedError := "Error checking mining status"

	// Validate result
	if result != expected && result != expectedAlreadyRunning && result != expectedError {
		t.Errorf("Unexpected result: %q", result)
	}

	// Log based on result
	switch result {
	case expected:
		t.Log("Mining started successfully.")
	case expectedAlreadyRunning:
		t.Log("Mining is already running.")
	case expectedError:
		t.Log("Error occurred while checking mining status.")
	}
}

// go test -v -run ^TestStopMining$ -count=1 application-layer/services
// TestStopMining validates stopping the mining process.
func TestStopMining(t *testing.T) {
	btcService := NewBtcService()

	// Initialize temp file
	SetupTempFilePath()

	// Retrieve mining address from temp file
	miningAddress, err := getMiningAddressFromTemp()
	if err != nil {
		t.Fatalf("Failed to retrieve mining address from temp file: %v", err)
	}

	// Log mining address
	fmt.Printf("Mining address retrieved from temp file: %s\n", miningAddress)

	// Call StopMining
	result := btcService.StopMining()

	// Validate result
	expected := "Mining process stopped and restarted successfully"
	if result != expected {
		t.Errorf("Expected %q, but got %q", expected, result)
	}

	// Log result
	t.Logf("StopMining function result: %s", result)

	// Cleanup: Stop btcd and btcwallet
	stopBtcdResult := btcService.StopBtcd()
	if stopBtcdResult != "btcd stopped successfully" {
		t.Logf("Failed to stop btcd: %s", stopBtcdResult)
	}

	stopBtcwalletResult := btcService.StopBtcwallet()
	if stopBtcwalletResult != "btcwallet stopped successfully" {
		t.Logf("Failed to stop btcwallet: %s", stopBtcwalletResult)
	}
}

// go test -v -run ^TestLogin$ -count=1 application-layer/services
// TestLogin validates logging into the wallet.
func TestLogin(t *testing.T) {
	btcService := NewBtcService()

	// Initialize temporary file path
	SetupTempFilePath()

	// Test wallet address and passphrase
	walletAddress := "15RMzownS37XMpPhGExBJoSQFHkePenw39"
	passphrase := "CSE416"

	// Ensure wallet exists for testing
	var walletDBPath string
	if runtime.GOOS == "windows" {
		walletDBPath = fmt.Sprintf(`%s\AppData\Local\Btcwallet\mainnet\wallet.db`, os.Getenv("USERPROFILE"))
	} else if runtime.GOOS == "darwin" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			t.Fatalf("Failed to get home directory: %v", err)
		}
		walletDBPath = filepath.Join(homeDir, "Library", "Application Support", "Btcwallet", "mainnet", "wallet.db")
	} else {
		t.Fatalf("Unsupported OS: %s", runtime.GOOS)
	}

	if _, err := os.Stat(walletDBPath); os.IsNotExist(err) {
		t.Fatalf("Wallet does not exist at path: %s. Please set up a test wallet first.", walletDBPath)
	}

	// Call Login
	result, err := btcService.Login(walletAddress, passphrase)
	if err != nil {
		t.Errorf("Login failed: %v", err)
	}

	// Validate results
	expected := "Wallet unlocked successfully"
	if result != expected {
		t.Errorf("Expected %q, but got %q", expected, result)
	}

	// Log output
	t.Logf("Login Result: %s", result)

	// Cleanup: Stop btcd and btcwallet
	stopBtcdResult := btcService.StopBtcd()
	if stopBtcdResult != "btcd stopped successfully" {
		t.Logf("Failed to stop btcd: %s", stopBtcdResult)
	}

	stopBtcwalletResult := btcService.StopBtcwallet()
	if stopBtcwalletResult != "btcwallet stopped successfully" {
		t.Logf("Failed to stop btcwallet: %s", stopBtcwalletResult)
	}
}

// go test -v -run ^TestGetBalance$ -count=1 application-layer/services
// TestGetBalance validates getting the wallet balance.
func TestGetBalance(t *testing.T) {
	btcService := NewBtcService()

	// Call GetBalance
	balance, err := btcService.GetBalance()
	if err != nil {
		t.Errorf("Failed to get balance: %v", err)
	}

	// Log and Print balance
	t.Logf("Wallet Balance: %s", balance)
	fmt.Printf("Wallet Balance (from test): %s\n", balance)
}

// go test -v -run ^TestGetReceivedByAddress$ -count=1 application-layer/services
// TestGetReceivedByAddress validates getting the received amount for an address.
func TestGetReceivedByAddress(t *testing.T) {
	btcService := NewBtcService()

	// Test wallet addres
	walletAddress := "14QnrKvCS9cskoMjfKkCe7xaWkQwdWCbJc"

	// Ensure btcd and btcwallet are running
	if !isProcessRunning("btcd") || !isProcessRunning("btcwallet") {
		t.Fatalf("btcd or btcwallet is not running. Please start both processes before testing")
	}

	// Call GetReceivedByAddress
	receivedAmount, err := btcService.GetReceivedByAddress(walletAddress)
	if err != nil {
		t.Errorf("Failed to get received amount for address %s: %v", walletAddress, err)
		return
	}

	// Log and Print received amount
	t.Logf("Received amount for address %s: %s", walletAddress, receivedAmount)
	fmt.Printf("Received amount for address %s (from test): %s\n", walletAddress, receivedAmount)
}

// go test -v -run ^TestGetBlockCount$ -count=1 application-layer/services
// TestGetBlockCount validates getting the current block count.
func TestGetBlockCount(t *testing.T) {
	btcService := NewBtcService()

	// Ensure btcd is running
	if !isProcessRunning("btcd") {
		t.Fatalf("btcd is not running. Please start btcd before testing")
	}

	// Call GetBlockCount
	blockCount, err := btcService.GetBlockCount()
	if err != nil {
		t.Errorf("Failed to get block count: %v", err)
		return
	}

	// Log and Print block count
	t.Logf("Current block count: %s", blockCount)
	fmt.Printf("Current block count (from test): %s\n", blockCount)
}

// go test -v -run ^TestListUnspent$ -count=1 application-layer/services
// TestListUnspent validates listing unspent transactions.
func TestListUnspent(t *testing.T) {
	btcService := NewBtcService()

	// Ensure btcd and btcwallet are running
	if !isProcessRunning("btcd") || !isProcessRunning("btcwallet") {
		t.Fatalf("btcd or btcwallet is not running. Please start both processes before testing")
	}

	// Call ListUnspent
	utxos, err := btcService.ListUnspent()
	if err != nil {
		t.Errorf("Failed to list unspent transactions: %v", err)
		return
	}

	// Log and Print full result
	t.Logf("List of unspent transactions: %v", utxos)
	fmt.Printf("List of unspent transactions (from test): %v\n", utxos)
}

// go test -v -run ^TestCreateRawTransactionWithValidation$ -count=1 application-layer/services
// this test is incomplete code block of the original function
func TestCreateRawTransactionWithValidation(t *testing.T) {
	btcService := NewBtcService()

	SetupTempFilePath()

	// Step 1: Retrieve source address (mining address) from temp file
	fmt.Println("Step 1: Reading temp file to retrieve mining address...")
	tempContent, err := ioutil.ReadFile(tempFilePath)
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}

	fmt.Printf("Raw temp file content: %s\n", tempContent)
	t.Logf("Raw temp file content: %s", tempContent)

	var tempData map[string]string
	err = json.Unmarshal(tempContent, &tempData)
	if err != nil {
		t.Fatalf("Failed to parse temp file content: %v", err)
	}

	fmt.Printf("Parsed temp file data: %v\n", tempData)
	t.Logf("Parsed temp file data: %v", tempData)

	miningAddress, exists := tempData["miningaddr"]
	if !exists || miningAddress == "" {
		t.Fatalf("Mining address not found in temp file")
	}

	fmt.Printf("Mining address retrieved: %s\n", miningAddress)
	t.Logf("Mining address retrieved: %s", miningAddress)

	// Step 2: Get UTXOs from ListUnspent
	fmt.Println("Step 2: Retrieving unspent transactions (UTXOs)...")
	utxos, err := btcService.ListUnspent()
	if err != nil {
		t.Fatalf("Failed to retrieve unspent transactions: %v", err)
	}

	fmt.Printf("List of unspent transactions: %v\n", utxos)
	t.Logf("List of unspent transactions: %v", utxos)

	// Step 3: Validate the chosen txid and amount
	chosenTxid := "9e97af2efb366d06bf6a75a3f5493b74a77fabf7805768f5b56eb70cca4eb3c1"
	requestedAmount := 50.0

	fmt.Println("Step 3: Validating chosen txid and amount...")
	isValid := false
	for _, utxo := range utxos {
		utxoTxid := utxo["txid"].(string)
		utxoAmount := utxo["amount"].(float64)

		if utxoTxid == chosenTxid {
			fmt.Printf("Found txid: %s with amount: %.8f\n", utxoTxid, utxoAmount)
			t.Logf("Found txid: %s with amount: %.8f", utxoTxid, utxoAmount)

			if utxoAmount >= requestedAmount {
				fmt.Printf("Valid transaction: txid %s has sufficient amount %.8f for requested %.8f\n", utxoTxid, utxoAmount, requestedAmount)
				t.Logf("Valid transaction: txid %s has sufficient amount %.8f for requested %.8f", utxoTxid, utxoAmount, requestedAmount)
				isValid = true
			} else {
				fmt.Printf("Insufficient amount: txid %s has %.8f but requested %.8f\n", utxoTxid, utxoAmount, requestedAmount)
				t.Logf("Insufficient amount: txid %s has %.8f but requested %.8f", utxoTxid, utxoAmount, requestedAmount)
			}
			break
		}
	}

	if !isValid {
		t.Fatalf("Invalid transaction: either txid %s not found or amount insufficient", chosenTxid)
	}

	fmt.Println("Transaction validation successful.")
	t.Log("Transaction validation successful.")
}

// go test -v -run ^TestCreateRawTransaction$ -count=1 application-layer/services
func TestCreateRawTransaction(t *testing.T) {
	// BTC 서비스 초기화
	btcService := NewBtcService()

	SetupTempFilePath()

	// 테스트 인자 설정
	txid := "224e895ef356e72f4b0809e2ccfe9d2eef8c83b157285624822929ec949797a1"
	dst := "15RMzownS37XMpPhGExBJoSQFHkePenw39"
	amount := 10.0

	// 함수 호출
	rawTx, err := btcService.CreateRawTransaction(txid, dst, amount)
	if err != nil {
		t.Fatalf("Failed to create raw transaction: %v", err)
	}

	// 결과 검증
	if rawTx == "" {
		t.Errorf("Expected a valid raw transaction, but got an empty string")
	}

	t.Logf("Raw transaction created successfully: %s", rawTx)
}

// go test -v -run ^TestSignRawTransaction$ -count=1 application-layer/services
func TestSignRawTransaction(t *testing.T) {
	// BTC 서비스 초기화
	btcService := NewBtcService()

	// 테스트용 Raw Transaction ID
	rawId := "0100000001a1979794ec29298224562857b1838cef2e9dfecce209084b2fe756f35e894e220000000000ffffffff0200286bee000000001976a914256837efb737ab0023119c429de6bb4b96546a2888ac00ca9a3b000000001976a914307c036da6d21136926389930a8d40c7d6850b5588ac00000000"

	// 함수 호출
	hex, complete, err := btcService.signRawTransaction(rawId)
	if err != nil {
		t.Fatalf("Failed to sign raw transaction: %v", err)
	}

	// 결과 검증
	if hex == "" {
		t.Errorf("Expected a valid signed transaction hex, but got an empty string")
	}
	if !complete {
		t.Errorf("Expected the transaction to be complete, but got incomplete")
	}

	t.Logf("Signed transaction hex: %s", hex)
	t.Logf("Transaction complete status: %v", complete)
}

// go test -v -run ^TestSendRawTransaction$ -count=1 application-layer/services
func TestSendRawTransaction(t *testing.T) {
	// BTC 서비스 초기화
	btcService := NewBtcService()

	// 테스트용 Signed Transaction Hex
	signedTxHex := "0100000001a1979794ec29298224562857b1838cef2e9dfecce209084b2fe756f35e894e22000000006b483045022100b177aa508b5e73a29b643db7c952591664815467c3cd6bd4c338fcfa97244193022057e0fe59c442bd8539f14639f9ba9f0031cffb53c5e0c4d27ea9cab29100bab401210300a4d118b08eff4d98c1a86bc6d19a563a78d0456869c65b7d0f97e4057691a9ffffffff0200286bee000000001976a914256837efb737ab0023119c429de6bb4b96546a2888ac00ca9a3b000000001976a914307c036da6d21136926389930a8d40c7d6850b5588ac00000000"

	// 함수 호출
	txid, err := btcService.sendRawTransaction(signedTxHex)
	if err != nil {
		t.Fatalf("Failed to send raw transaction: %v", err)
	}

	// 결과 검증
	if txid == "" {
		t.Errorf("Expected a valid transaction ID, but got an empty string")
	}

	t.Logf("Transaction ID: %s", txid)
}

// go test -v -run ^TestTransaction$ -count=1 application-layer/services
func TestTransaction(t *testing.T) {

	// Initialize BtcService
	btcService := NewBtcService()

	// Setup temporary file path
	SetupTempFilePath()

	// Test parameters
	passphrase := "CSE416"
	txid := "6badadaa3cd1c6ac32e04b37d036b3b10fa9b330dc558ac4123eade72a8ba681"
	dst := "15RMzownS37XMpPhGExBJoSQFHkePenw39"
	amount := 5.0

	// Step 1: Call the Transaction function
	fmt.Println("Starting transaction...")
	txIdResult, err := btcService.Transaction(passphrase, txid, dst, amount)
	if err != nil {
		t.Fatalf("Transaction failed: %v", err)
	}

	// Log the transaction ID result
	fmt.Printf("Transaction completed successfully. TxID: %s\n", txIdResult)
	t.Logf("Transaction completed successfully. TxID: %s", txIdResult)

	// Additional validation (if required)
	currentBalance, err := btcService.GetBalance()
	if err != nil {
		t.Fatalf("Failed to retrieve current balance: %v", err)
	}

	// Log current balance
	fmt.Printf("Current balance after transaction: %s\n", currentBalance)
	t.Logf("Current balance after transaction: %s", currentBalance)

	// Validate balance logic here if necessary
}

// go test -v -run ^TestBtcwalletCreate_PowerShell$
// TestBtcwalletCreate_PowerShell는 btcwallet 생성 테스트를 수행합니다.
func TestBtcwalletCreate_PowerShell(t *testing.T) {

	btcService := NewBtcService()
	// %USERPROFILE% 환경 변수에서 사용자 홈 경로 가져오기
	userProfile := os.Getenv("USERPROFILE")
	if userProfile == "" {
		t.Fatal("USERPROFILE environment variable is not set")
	}

	// wallet.db 경로 동적으로 설정
	walletDBPath := filepath.Join(userProfile, "AppData", "Local", "Btcwallet", "mainnet", "wallet.db")

	// 기존 지갑 삭제 (테스트를 위해 기존 지갑이 있다면 삭제)
	err := os.Remove(walletDBPath)
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("Failed to delete existing wallet.db: %v", err)
	}

	// btcwallet 생성 함수 호출
	err = btcService.BtcwalletCreate("CSE416")
	if err != nil {
		t.Fatalf("Failed to create btcwallet via PowerShell: %v", err)
	}

	// 지갑이 생성되었는지 확인
	if _, err := os.Stat(walletDBPath); os.IsNotExist(err) {
		t.Fatal("wallet.db does not exist after creation")
	} else {
		fmt.Println("wallet.db successfully created.")
	}
}
