package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/gofrs/flock"
)

var fileMutex sync.Mutex

// use for user uploaded files and user downlaoded files
func updateFile(newFileData DHTMetadata, dirPath string, filePath string, isDelete bool) error {
	// make thread-safe in case multiple nodes send to cloud node
	fileMutex.Lock()
	defer fileMutex.Unlock()

	lock := flock.New(filePath + ".lock")
	defer lock.Unlock()

	locked, err := lock.TryLock()
	if err != nil || !locked {
		return fmt.Errorf("failed to get file lock: %v", err)
	}

	// check if directory and file exist
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.Mkdir(dirPath, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create utils directory: %v", err)
		}
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := os.WriteFile(filePath, []byte("[]"), 0644); err != nil {
			return fmt.Errorf("failed to create files.json: %v", err)
		}
	}

	fmt.Println("file to be added/updated: ", newFileData)

	// read in JSON file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read files.json: %v", err)
	}

	var files []DHTMetadata
	if err := json.Unmarshal(data, &files); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	if fileNameToHashMap == nil {
		fileNameToHashMap = make(map[string][]string)
	}	

	if isDelete {
		// Remove from fileNameToHashMap
		if hashes, exists := fileNameToHashMap[newFileData.Name]; exists {
			// Find the hash in the array and remove it
			for i, hash := range hashes {
				if hash == newFileData.Hash {
					// Remove hash from slice
					fileNameToHashMap[newFileData.Name] = append(hashes[:i], hashes[i+1:]...)
					break
				}
			}

			// If no hashes are left for this file name, delete the key
			if len(fileNameToHashMap[newFileData.Name]) == 0 {
				delete(fileNameToHashMap, newFileData.Name)
			}
		}

		// Remove the file from the files slice
		for i, file := range files {
			if file.Hash == newFileData.Hash {
				// Remove file from slice
				files = append(files[:i], files[i+1:]...)
				break
			}
		}
	} else {
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
			files = append([]DHTMetadata{newFileData}, files...)
		}

		// update fileNameToHashMap
		if _, exists := fileNameToHashMap[newFileData.Name]; !exists {
			fileNameToHashMap[newFileData.Name] = []string{}
		}
		fileNameToHashMap[newFileData.Name] = append(fileNameToHashMap[newFileData.Name], newFileData.Hash)
	}

	// convert updated list of files back to JSON
	updatedData, err := json.MarshalIndent(files, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	err = os.WriteFile(filePath, updatedData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write updated data to files.json: %v", err)
	}

	return nil
}
