package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// Get-Process | Where-Object {$_.Name -like "btc*"}

// C:\dev\workspace\CSE-416\application-layer\services> go test -v -run ^TestStopBtcd$
// go test -v -run ^TestStartBtcdWithNoArgs$
func TestStartBtcdWithNoArgs(t *testing.T) {

	btcService := NewBtcService()
	// 초기화 호출
	btcService.Init()

	fmt.Println("Starting TestStartBtcdWithNoArgs...") // 테스트 시작 메시지 출력

	result := btcService.StartBtcd()
	expectedSuccess := "btcd started successfully"
	expectedFailure := "Error starting btcd"
	expectedAlreadyRunning := "btcd is already running"

	fmt.Printf("Test Result: %s\n", result) // 실행 결과 출력

	if result != expectedSuccess && result != expectedFailure && result != expectedAlreadyRunning {
		t.Errorf("Unexpected result: %q", result)
	}

	if result == expectedSuccess {
		t.Log("btcd started successfully with no arguments.")
		fmt.Println("btcd started successfully.") // 추가 출력
	} else if result == expectedFailure {
		fmt.Println("btcd failed to start.")
	} else if result == expectedAlreadyRunning {
		fmt.Println("btcd is already running.")
	}

	fmt.Println("TestStartBtcdWithNoArgs completed.") // 테스트 종료 메시지 출력
}

// go test -v -run ^TestStartBtcdWithArgs$
func TestStartBtcdWithArgs(t *testing.T) {
	btcService := NewBtcService()

	// 초기화 호출
	SetupTempFilePath()

	result := btcService.StartBtcd("1B5t2bk3BtCw88uveEFbvFERotX6adGY6w")
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

func TestStartBtcdWithArgs2(t *testing.T) {
	btcService := NewBtcService()

	// 초기화 호출
	SetupTempFilePath()

	result := btcService.StartBtcd("13NPW1mgHkJv3tAkogtd3hMvAwoid2YgkU")
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
func TestStartBtcdWithInvalidArgs(t *testing.T) {
	btcService := NewBtcService()

	result := btcService.StartBtcd("1ExampleWalletAddress", "AnotherArgument")
	expected := "Invalid number of arguments"

	if result != expected {
		t.Errorf("Expected %q, but got %q", expected, result)
	}
}

func TestStartBtcwallet(t *testing.T) {
	btcService := NewBtcService()

	// 인자가 없는 경우 btcwallet 실행
	result := btcService.StartBtcwallet()
	expected := "btcwallet started successfully"

	if result != expected {
		t.Errorf("Expected %q, but got %q", expected, result)
	}
}

// C:\dev\workspace\CSE-416\application-layer\services> go test -v -run ^TestStopBtcd$
func TestStopBtcd(t *testing.T) {
	btcService := NewBtcService()

	result := btcService.StopBtcd()
	t.Logf("Result: %s", result)
	if result != "btcd is not running" {
		t.Errorf("Expected 'btcd is not running', but got %q", result)
	}

	t.Logf("Result: %s", result)
	if result != "btcd stopped successfully" {
		t.Errorf("Expected 'btcd stopped successfully', but got %q", result)
	}
}
func TestStopBtcwallet(t *testing.T) {
	btcService := NewBtcService()

	// StopBtcwallet 호출
	result := btcService.StopBtcwallet()

	// 예상 가능한 결과
	expectedNotRunning := "btcwallet is not running"
	expectedStopped := "btcwallet stopped successfully"

	// 실행 결과 확인
	if result != expectedNotRunning && result != expectedStopped {
		t.Errorf("Unexpected result: %q", result)
	}

	// 추가적으로 예상된 로그를 출력
	if result == expectedNotRunning {
		t.Log("btcwallet is not running.")
	} else if result == expectedStopped {
		t.Log("btcwallet stopped successfully.")
	}
}

// go test -v -run ^TestInit$
func TestInit(t *testing.T) {
	btcService := NewBtcService()

	// Init 호출
	result := btcService.Init()

	// 예상 결과
	expected := "Initialization and cleanup completed successfully"

	// 결과 검증
	if result != expected {
		t.Errorf("Expected %q, but got %q", expected, result)
	}

	// 추가 로그
	t.Logf("Init function result: %s", result)
}

// go test -v -run ^TestGetMiningStatus$
func TestGetMiningStatus(t *testing.T) {
	btcService := NewBtcService()

	// 마이닝 상태 확인
	status, err := btcService.GetMiningStatus()
	if err != nil {
		t.Errorf("Error checking mining status: %v", err)
	}

	// 상태 출력
	t.Logf("Mining status: %t", status)
}

// go test -v -run ^TestStartMining$
func TestStartMining(t *testing.T) {
	btcService := NewBtcService()

	result := btcService.StartMining(5) // 블록 5개 생성 요청

	// 예상 가능한 결과
	expected := "mining started successfully"
	expectedAlreadyRunning := "mining is running"
	expectedError := "Error checking mining status"

	// 결과 검증
	if result != expected && result != expectedAlreadyRunning && result != expectedError {
		t.Errorf("Unexpected result: %q", result)
	}

	// 결과에 따라 로그 출력
	switch result {
	case expected:
		t.Log("Mining started successfully.")
	case expectedAlreadyRunning:
		t.Log("Mining is already running.")
	case expectedError:
		t.Log("Error occurred while checking mining status.")
	}
}

// go test -v -run ^TestStopMining$
func TestStopMining(t *testing.T) {
	btcService := NewBtcService()

	// 1. temp 파일 초기화
	SetupTempFilePath()

	// temp 파일에 저장된 마이닝 주소 확인
	miningAddress, err := getMiningAddressFromTemp()
	if err != nil {
		t.Fatalf("Failed to retrieve mining address from temp file: %v", err)
	}

	// 마이닝 주소 출력
	fmt.Printf("Mining address retrieved from temp file: %s\n", miningAddress)

	// 2. StopMining 호출
	result := btcService.StopMining()

	// 결과 검증
	expected := "Mining process stopped and restarted successfully"
	if result != expected {
		t.Errorf("Expected %q, but got %q", expected, result)
	}

	// 추가 로그
	t.Logf("StopMining function result: %s", result)
}

// go test -v -run ^TestGetNewAddress$
func TestGetNewAddress(t *testing.T) {
	btcService := NewBtcService()

	// GetNewAddress 호출
	newAddress, err := btcService.GetNewAddress()
	if err != nil {
		t.Fatalf("Failed to get new address: %v", err)
	}

	// 생성된 주소 출력
	t.Logf("Generated new address: %s", newAddress)
}

// go test -v -run ^TestUnlockWallet$
func TestUnlockWallet(t *testing.T) {
	btcService := NewBtcService()

	// 테스트용 passphrase 설정
	passphrase := "CSE416"

	// btcd와 btcwallet이 실행 중인지 확인
	if !isProcessRunning("btcd") || !isProcessRunning("btcwallet") {
		t.Fatalf("btcd or btcwallet is not running. Please start both processes before testing")
	}

	// UnlockWallet 호출
	result, err := btcService.UnlockWallet(passphrase)
	if err != nil {
		t.Errorf("Failed to unlock wallet: %v", err)
	}

	// 결과 검증
	expected := "" // btcctl 명령어는 일반적으로 출력이 없으므로 빈 문자열 예상
	if result != expected {
		t.Errorf("Expected %q, but got %q", expected, result)
	}

	// 로그 출력
	t.Logf("UnlockWallet Result: %s", result)
}

func TestLockWallet(t *testing.T) {
	btcService := NewBtcService()

	// btcd와 btcwallet이 실행 중인지 확인
	if !isProcessRunning("btcd") || !isProcessRunning("btcwallet") {
		t.Fatalf("btcd or btcwallet is not running. Please start both processes before testing")
	}

	// LockWallet 호출
	result, err := btcService.LockWallet()
	if err != nil {
		t.Errorf("Failed to lock wallet: %v", err)
	}

	// 결과 검증
	expected := "" // btcctl 명령어는 일반적으로 출력이 없으므로 빈 문자열 예상
	if result != expected {
		t.Errorf("Expected %q, but got %q", expected, result)
	}

	// 로그 출력
	t.Logf("LockWallet Result: %s", result)
}

// go test -v -run ^TestLogin$
func TestLogin(t *testing.T) {
	btcService := NewBtcService()

	// 초기화 호출
	SetupTempFilePath()

	// 테스트용 walletAddress와 passphrase
	walletAddress := "1B5t2bk3BtCw88uveEFbvFERotX6adGY6w"
	passphrase := "CSE416"

	// Login 호출
	result, err := btcService.Login(walletAddress, passphrase)
	if err != nil {
		t.Errorf("Login failed: %v", err)
	}

	// 결과 검증
	expected := "Wallet unlocked successfully"
	if result != expected {
		t.Errorf("Expected %q, but got %q", expected, result)
	}

	// 로그 출력
	t.Logf("Login Result: %s", result)
}

// go test -v -run ^TestGetBalance$
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

// go test -v -run ^TestGetReceivedByAddress$
func TestGetReceivedByAddress(t *testing.T) {
	btcService := NewBtcService()

	// 테스트용 walletAddress 설정
	walletAddress := "1B5t2bk3BtCw88uveEFbvFERotX6adGY6w"

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

// go test -v -run ^TestGetBlockCount$
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

// go test -v -run ^TestListReceivedByAddress$
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

	// Log and Print address list
	t.Logf("Received Addresses: %v", addressList)
	fmt.Printf("Received Addresses (from test): %v\n", addressList)
}

// go test -v -run ^TestListUnspent$
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

// go test -v -run ^TestCreateRawTransactionWithValidation$
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

// go test -v -run ^TestTransaction$
func TestTransaction(t *testing.T) {
	// Initialize BtcService
	btcService := NewBtcService()

	// Setup temporary file path
	SetupTempFilePath()

	// Test parameters
	passphrase := "CSE416"
	txid := "3a085fd6c155d244c98dc9ca720242a7bb438fd543e6c51578a7424987ad508b"
	dst := "1G23RBEZVhePDeTP5gq4be3jEec5mBWorw"
	amount := 1.2

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
	err = BtcwalletCreate("CSE416")
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
