package utils

import (
	"application-layer/models"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var (
	dirPath              = filepath.Join("..", "utils")
	proxyHistoryFilePath = filepath.Join(dirPath, "proxyHistory.json")
	uploadedFilePath     = filepath.Join(dirPath, "files.json")
	downloadedFilePath   = filepath.Join(dirPath, "downloadedFiles.json")
	transactionFilePath  = filepath.Join(dirPath, "transactionFiles.json")

	fileCopyPath = filepath.Join("..", "squidcoinFiles")
	fileMutex    sync.Mutex // used by cloud-node
)

// use for user uploaded files and user downlaoded files
func SaveOrUpdateFile(newFileData models.FileMetadata, dirPath, filePath string) (string, error) {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	// check if directory and file exist
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.Mkdir(dirPath, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create utils directory: %v", err)
		}
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := os.WriteFile(filePath, []byte("[]"), 0644); err != nil {
			return "", fmt.Errorf("failed to create files.json: %v", err)
		}
	}

	fmt.Println("file to be added/updated: ", newFileData)

	// read in JSON file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read files.json: %v", err)
	}

	var files []models.FileMetadata
	if err := json.Unmarshal(data, &files); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %v", err)
	}

	// update if file is already in JSON file
	isUpdated := false
	for i := range files {
		if files[i].Hash == newFileData.Hash {
			files[i] = newFileData
			isUpdated = true
			break
		}
	}

	// add file if not already in JSON file
	if !isUpdated {
		files = append([]models.FileMetadata{newFileData}, files...)
	}

	// convert updated list of files back to JSON
	updatedData, err := json.MarshalIndent(files, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %v", err)
	}

	if err := os.WriteFile(filePath, updatedData, 0644); err != nil {
		return "", fmt.Errorf("failed to write updated data to files.json: %v", err)
	}

	action := "updated"
	if !isUpdated {
		action = "added"
	}

	return action, nil
}

func AddOrUpdateTransaction(transaction models.Transaction) error {
	// check if directory and file exist
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.Mkdir(dirPath, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create utils directory: %v", err)
		}
	}

	if _, err := os.Stat(transactionFilePath); os.IsNotExist(err) {
		if err := os.WriteFile(transactionFilePath, []byte("[]"), 0644); err != nil {
			return fmt.Errorf("failed to create transactionFiles.json: %v", err)
		}
	}

	// read in JSON file
	data, err := os.ReadFile(transactionFilePath)
	if err != nil {
		fmt.Errorf("failed to read files.json: %v", err)
	}

	var transactions []models.Transaction
	if err := json.Unmarshal(data, &transactions); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	// update if transaction is already in JSON file
	isUpdated := false
	for i := range transactions {
		if transactions[i].TransactionID == transaction.TransactionID {
			transactions[i] = transaction
			isUpdated = true
			break
		}
	}

	// add file if not already in JSON file
	if !isUpdated {
		transactions = append([]models.Transaction{transaction}, transactions...)
	}

	// convert updated list of files back to JSON
	updatedData, err := json.MarshalIndent(transactions, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	if err := os.WriteFile(transactionFilePath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write updated data to transactionFiles.json: %v", err)
	}

	return nil
}
func AddOrUpdateProxyHistory(newHistoryEntry models.ProxyHistoryEntry) (string, error) {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	// Check if the directory and file exist
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.Mkdir(dirPath, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create utils directory: %v", err)
		}
	}

	if _, err := os.Stat(proxyHistoryFilePath); os.IsNotExist(err) {
		if err := os.WriteFile(proxyHistoryFilePath, []byte("[]"), 0644); err != nil {
			return "", fmt.Errorf("failed to create proxyHistory.json: %v", err)
		}
	}

	// Read the current proxy history from the file
	data, err := os.ReadFile(proxyHistoryFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read proxyHistory.json: %v", err)
	}

	var proxyHistory []models.ProxyHistoryEntry
	if err := json.Unmarshal(data, &proxyHistory); err != nil {
		return "", fmt.Errorf("failed to parse proxyHistory.json: %v", err)
	}

	// Update if the history entry for the client already exists
	isUpdated := false
	for i := range proxyHistory {
		if proxyHistory[i].ClientPeerID == newHistoryEntry.ClientPeerID {
			proxyHistory[i] = newHistoryEntry
			isUpdated = true
			break
		}
	}

	// Add the new history entry if not already present
	if !isUpdated {
		proxyHistory = append([]models.ProxyHistoryEntry{newHistoryEntry}, proxyHistory...)
	}

	// Convert the updated proxy history back to JSON
	updatedData, err := json.MarshalIndent(proxyHistory, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal updated proxy history: %v", err)
	}
	// Write the updated proxy history to the file
	if err := os.WriteFile(proxyHistoryFilePath, updatedData, 0644); err != nil {
		return "", fmt.Errorf("failed to write updated data to proxyHistory.json: %v", err)
	}

	action := "updated"
	if !isUpdated {
		action = "added"
	}
	return action, nil
}
