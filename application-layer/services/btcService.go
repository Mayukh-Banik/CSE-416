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
)

// 글로벌 변수로 실행 파일 경로 설정
var (
	btcdPath      = "../../btcd/btcd"
	btcwalletPath = "../../btcwallet/btcwallet"
	btcctlPath    = "../../btcd/cmd/btcctl/btcctl"
	tempFilePath  string
)

type BtcService struct{}

func CheckDirectoryContents(parentDir string) string {
	// 경로 확인
	// parentDir := "../../btcd"
	_, err := utils.CheckDirectoryContents(parentDir)
	if err != nil {
		fmt.Printf("Error checking directory: %v\n", err)
		return "Failed to check directory"
	}
	return "successed to check directory"
}

// init 함수 이름을 변경한 초기화 함수
func SetupTempFilePath() {
	if os.Getenv("OS") == "Windows_NT" {
		tempFilePath = filepath.Join(os.Getenv("TEMP"), "btc_temp.json")
	} else {
		tempFilePath = "/tmp/btcd_temp.json"
	}
	fmt.Printf("Temporary file path set to: %s\n", tempFilePath)
}

func deleteFromTempFile(key string) error {
	// 파일 읽기
	content, err := ioutil.ReadFile(tempFilePath)
	if err != nil {
		return fmt.Errorf("failed to read temp file: %w", err)
	}

	// JSON 파싱
	var data map[string]string
	if err := json.Unmarshal(content, &data); err != nil {
		return fmt.Errorf("failed to unmarshal temp file content: %w", err)
	}

	// 키 삭제
	delete(data, key)

	// JSON 직렬화
	updatedContent, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated data: %w", err)
	}

	// 파일 쓰기
	if err := ioutil.WriteFile(tempFilePath, updatedContent, 0644); err != nil {
		return fmt.Errorf("failed to write updated temp file: %w", err)
	}

	return nil
}

// 임시 파일 읽기 및 업데이트 함수
func updateTempFile(key, value string) error {
	var data map[string]string

	// 임시 파일 읽기
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

	// 데이터 업데이트
	data[key] = value

	// 파일 쓰기
	newContent, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated data: %w", err)
	}

	if err := ioutil.WriteFile(tempFilePath, newContent, 0644); err != nil {
		return fmt.Errorf("failed to write updated temp file: %w", err)
	}

	return nil
}

// isProcessRunning는 특정 이름의 프로세스가 실행 중인지 확인하는 함수입니다.
func isProcessRunning(processName string) bool {
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("Get-Process | Where-Object {$_.Name -like '%s'}", processName))
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Run(); err != nil {
		return false
	}

	return output.String() != ""
}

func NewBtcService() *BtcService {
	return &BtcService{}
}

// StartBtcd는 btcd를 시작하고 실행 여부를 확인하는 함수입니다.
func (bs *BtcService) StartBtcd(walletAddress ...string) string {
	// btcd 실행 여부 확인
	if isProcessRunning("btcd") {
		fmt.Println("btcd is already running. Cannot start another instance.")
		return "btcd is already running"
	}

	var cmd *exec.Cmd

	// 인자가 0개인 경우 기본 실행
	if len(walletAddress) == 0 {
		cmd = exec.Command(
			btcdPath,
			"--rpcuser=user",
			"--rpcpass=password",
			"--notls",
		)
	} else if len(walletAddress) == 1 {
		// 인자가 1개인 경우 실행
		cmd = exec.Command(
			btcdPath,
			"--rpcuser=user",
			"--rpcpass=password",
			"--notls",
			fmt.Sprintf("--miningaddr=%s", walletAddress[0]),
		)
	} else {
		// 인자가 1개 초과인 경우 오류 반환
		fmt.Println("Invalid number of arguments. Only 0 or 1 argument is allowed.")
		return "Invalid number of arguments"
	}

	// Detached mode 설정
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP, // Windows에서 프로세스를 독립적으로 실행
	}

	// 표준 출력과 에러 연결
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// btcd 프로세스 실행
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting btcd: %v\n", err)
		return "Error starting btcd"
	}

	fmt.Printf("btcd started with PID: %d\n", cmd.Process.Pid)

	// 실행된 btcd 프로세스 확인
	if !isProcessRunning("btcd") {
		fmt.Println("btcd process not found after starting.")
		return "btcd process not found"
	}

	fmt.Println("btcd is running")

	// 실행 성공 시 임시 파일 업데이트

	if len(walletAddress) == 0 {
		// 마이닝 주소 삭제
		if err := deleteFromTempFile("miningaddr"); err != nil {
			fmt.Printf("Failed to delete miningaddr from temp file: %v\n", err)
			return "btcd started but failed to clear mining address from temp file"
		}
		fmt.Println("Mining address cleared from temporary file.")
	} else if len(walletAddress) == 1 {
		if err := updateTempFile("miningaddr", walletAddress[0]); err != nil {
			fmt.Printf("Failed to update temp file: %v\n", err)
			return "btcd started but failed to update temp file"
		}
		fmt.Println("Temporary file updated successfully.")
	}

	return "btcd started successfully"
}

// StopBtcd는 btcd 프로세스를 종료하는 함수입니다.
func (bs *BtcService) StopBtcd() string {
	// btcd 실행 여부 확인
	checkProcessCmd := exec.Command("powershell", "-Command", "Get-Process | Where-Object {$_.Name -eq 'btcd'}")
	var checkOutput bytes.Buffer
	checkProcessCmd.Stdout = &checkOutput
	checkProcessCmd.Stderr = &checkOutput

	// 실행 여부를 확인하고 출력
	fmt.Println("Checking for running btcd process...")
	err := checkProcessCmd.Run()
	fmt.Printf("Check Process Output: %s\n", checkOutput.String())
	if err != nil || checkOutput.String() == "" {
		fmt.Println("btcd is not running. Cannot stop.")
		return "btcd is not running"
	}

	// btcd 프로세스 종료
	fmt.Println("Stopping btcd process...")
	cmd := exec.Command("taskkill", "/IM", "btcd.exe", "/F")
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	err = cmd.Run()
	fmt.Printf("Taskkill Output: %s\n", output.String())
	if err != nil {
		fmt.Printf("Error stopping btcd: %v\n", err)
		return fmt.Sprintf("Error stopping btcd: %s", output.String())
	}

	fmt.Println("btcd stopped successfully")
	return "btcd stopped successfully"
}

// StartBtcwallet는 btcwallet 프로세스를 시작하는 함수입니다.
func (bs *BtcService) StartBtcwallet() string {
	// btcwallet 실행 여부 확인
	if isProcessRunning("btcwallet") {
		fmt.Println("btcwallet is already running. Cannot start another instance.")
		return "btcwallet is already running"
	}

	// btcwallet 실행 커맨드
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

	// 표준 출력과 에러 연결
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// btcwallet 프로세스 실행
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting btcwallet: %v\n", err)
		return "Error starting btcwallet"
	}

	fmt.Printf("btcwallet started with PID: %d\n", cmd.Process.Pid)

	// 실행된 btcwallet 프로세스 확인
	if !isProcessRunning("btcwallet") {
		fmt.Println("btcwallet process not found after starting.")
		return "btcwallet process not found"
	}

	fmt.Println("btcwallet is running")
	return "btcwallet started successfully"
}

// StopBtcwallet는 btcwallet 프로세스를 종료하는 함수입니다.
func (bs *BtcService) StopBtcwallet() string {
	// btcwallet 실행 여부 확인
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

func initializeTempFile() error {
	// 초기 데이터
	initialData := map[string]string{
		"status": "initialized",
	}

	// JSON 직렬화
	dataBytes, err := json.MarshalIndent(initialData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal initial data: %w", err)
	}

	// 임시 파일 쓰기
	err = ioutil.WriteFile(tempFilePath, dataBytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to temp file: %w", err)
	}

	return nil
}

// Init은 btcd와 btcwallet을 시작하고 TA 서버에 연결한 뒤 종료하는 초기화 함수입니다.
func (bs *BtcService) Init() string {
	SetupTempFilePath()

	// 임시 파일 초기화
	err := initializeTempFile()
	if err != nil {
		fmt.Printf("Failed to initialize temp file: %v\n", err)
		return "Failed to initialize temp file"
	}
	fmt.Println("Temporary file initialized successfully.")

	// 1. btcd 시작
	btcdResult := bs.StartBtcd()
	if btcdResult != "btcd started successfully" {
		stopBtcd()
		stopBtcwallet()
		fmt.Println("Failed to start btcd.")
		return btcdResult
	}

	// 2. btcwallet 시작
	btcwalletResult := bs.StartBtcwallet()
	if btcwalletResult != "btcwallet started successfully" {
		stopBtcd()
		stopBtcwallet()
		fmt.Println("Failed to start btcwallet.")
		return btcwalletResult
	}

	// 3. TA 서버 연결 (btcctl 명령 실행)
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

	// 4. btcd 종료
	stopBtcdResult := bs.StopBtcd()
	if stopBtcdResult != "btcd stopped successfully" {
		fmt.Println("Failed to stop btcd.")
		return stopBtcdResult
	}

	// 5. btcwallet 종료
	stopBtcwalletResult := bs.StopBtcwallet()
	if stopBtcwalletResult != "btcwallet stopped successfully" {
		fmt.Println("Failed to stop btcwallet.")
		return stopBtcwalletResult
	}

	fmt.Println("Initialization and cleanup completed successfully.")
	return "Initialization and cleanup completed successfully"
}

// getMiningStatus는 마이닝이 활성화되어 있는지 확인하는 함수입니다.
func (bs *BtcService) GetMiningStatus() (bool, error) {
	// btcctl 명령 실행
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

	// 명령 실행
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error executing btcctl getgenerate: %v\n", err)
		return false, fmt.Errorf("error executing btcctl getgenerate: %w", err)
	}

	// 출력 결과 분석
	result := output.String()
	fmt.Printf("getgenerate output: %s\n", result)

	// 결과에 따라 true 또는 false 반환
	if result == "true\n" {
		return true, nil
	}
	return false, nil
}

// StartMining는 지정된 수의 블록을 생성하는 함수입니다.
func (bs *BtcService) StartMining(numBlock int) string {
	// 현재 마이닝 상태 확인
	isMining, err := bs.GetMiningStatus()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return "Error checking mining status"
	}

	if isMining {
		fmt.Println("Mining is already running.")
		return "mining is running"
	}

	// btcctl 명령 실행
	cmd := exec.Command(
		btcctlPath,
		"--wallet",
		"--rpcuser=user",
		"--rpcpass=password",
		"--rpcserver=127.0.0.1:8332",
		"--notls",
		"generate",
		strconv.Itoa(numBlock), // numBlock을 문자열로 변환
	)

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	// 명령 실행
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error starting mining: %v\n", err)
		return fmt.Sprintf("Error starting mining: %s", err.Error())
	}

	fmt.Printf("Mining started successfully. Output: %s\n", output.String())
	return "mining started successfully"
}

func getMiningAddressFromTemp() (string, error) {
	// 파일 읽기
	content, err := ioutil.ReadFile(tempFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read temp file: %w", err)
	}

	// JSON 파싱
	var data map[string]string
	if err := json.Unmarshal(content, &data); err != nil {
		return "", fmt.Errorf("failed to unmarshal temp file content: %w", err)
	}

	// 마이닝 주소 가져오기
	miningAddress, exists := data["miningaddr"]
	if !exists || miningAddress == "" {
		return "", fmt.Errorf("mining address not found in temp file")
	}

	return miningAddress, nil
}

func (bs *BtcService) StopMining() string {
	// 1. 마이닝 상태 확인
	isMining, err := bs.GetMiningStatus()
	if err != nil {
		fmt.Printf("Failed to check mining status: %v\n", err)
		return "Error checking mining status"
	}

	// 2. 마이닝이 활성화되지 않은 경우
	if !isMining {
		fmt.Println("Mining is not active. No action needed.")
		return "Mining is not active"
	}

	// 3. 임시 파일에서 마이닝 주소 가져오기
	miningAddress, err := getMiningAddressFromTemp()
	if err != nil {
		fmt.Printf("Failed to retrieve mining address: %v\n", err)
		return "Failed to retrieve mining address"
	}
	fmt.Printf("Retrieved mining address: %s\n", miningAddress)

	// 4. btcd 및 btcwallet 중지
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

	// 5. btcd 및 btcwallet 재시작
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

	// 6. 작업 완료 메시지
	fmt.Println("Mining process stopped and restarted successfully.")
	return "Mining process stopped and restarted successfully"
}

func (bs *BtcService) GetNewAddress() (string, error) {
	// btcd와 btcwallet이 실행 중인지 확인
	btcdRunning := isProcessRunning("btcd")
	btcwalletRunning := isProcessRunning("btcwallet")

	// btcd가 실행 중인지 확인
	if !btcdRunning {
		return "", fmt.Errorf("btcd is not running. Please start btcd before calling this function")
	}

	// btcwallet이 실행 중인지 확인
	if !btcwalletRunning {
		return "", fmt.Errorf("btcwallet is not running. Please start btcwallet before calling this function")
	}

	// 3. 새 주소 생성
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
		return "", fmt.Errorf("Error generating new address: %w", err)
	}

	// 생성된 주소
	newAddress := strings.TrimSpace(output.String())
	fmt.Printf("Generated new address: %s\n", newAddress)

	// 최종적으로 새 주소 반환
	return newAddress, nil
}
