package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"runtime"
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

//go test -v -run ^TestStartBtcdWithNoArgsAndStop$ -count=1 application-layer/services
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


func TestStartBtcwallet(t *testing.T) {
	btcService := NewBtcService()

	// 인자가 없는 경우 btcwallet 실행
	result := btcService.StartBtcwallet()
	expected := "btcwallet started successfully"

	if result != expected {
		t.Errorf("Expected %q, but got %q", expected, result)
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

// go test -v -run ^TestCreateRawTransaction$
func TestCreateRawTransaction(t *testing.T) {
	// BTC 서비스 초기화
	btcService := NewBtcService()

	SetupTempFilePath()

	// 테스트 인자 설정
	txid := "cf38e31e633110e35a0fd91c16807904a4e38d2acbfd5d5a985d36a7240fe702"
	dst := "1G23RBEZVhePDeTP5gq4be3jEec5mBWorw"
	amount := 20.0

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

func TestUnlockWallet2(t *testing.T) {
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

// go test -v -run ^TestSignRawTransaction$
func TestSignRawTransaction(t *testing.T) {
	// BTC 서비스 초기화
	btcService := NewBtcService()

	// 테스트용 Raw Transaction ID
	rawId := "010000000102e70f24a7365d985a5dfdcb2a8de3a4047980161cd90f5ae31031631ee338cf0000000000ffffffff0280e2eeb2000000001976a9146e9d847156b018b35e6275c150cb2788b24c90a188ac00943577000000001976a914a4bc51f8084fffa6a7432568cf46f1353a0666ca88ac00000000"

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

// go test -v -run ^TestSendRawTransaction$
func TestSendRawTransaction(t *testing.T) {
	// BTC 서비스 초기화
	btcService := NewBtcService()

	// 테스트용 Signed Transaction Hex
	signedTxHex := "010000000171f1ebd5a5a021005107273671642181aaa9457710636efa7b7362b56f5de3f9000000006a47304402200f53775acd4fe46ae4c52fd09b39408798cca6d6e5e709cba17b76451ae282d802204872576f48be483a16bfb5d1272bae4478f34be8d24f20bad5121d6c08af764201210387d22a806b62c919e5f461f6038f583e40bcc591eb2e4e562802ce9e56b9dea2ffffffff02c0512677000000001976a9146e9d847156b018b35e6275c150cb2788b24c90a188ac00943577000000001976a914a4bc51f8084fffa6a7432568cf46f1353a0666ca88ac00000000"

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

// go test -v -run ^TestTransaction$
func TestTransaction(t *testing.T) {

	// Initialize BtcService
	btcService := NewBtcService()

	// Setup temporary file path
	SetupTempFilePath()

	// Test parameters
	passphrase := "CSE416"
	txid := "5dfcbcd73fb25f65c335fef7fe63701d1b0ddb743b2d80c84209515f82f67d0a"
	dst := "1G23RBEZVhePDeTP5gq4be3jEec5mBWorw"
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

// go test -v -run ^TestBtcwalletCreate$ -count=1 application-layer/services
// TestBtcwalletCreate btcwallet implementation test for macOs and Windows Powershell
func TestBtcwalletCreate(t *testing.T) {
	var walletDBPath string

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
	err = BtcwalletCreate(passphrase)
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



