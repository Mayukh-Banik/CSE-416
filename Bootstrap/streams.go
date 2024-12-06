package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	dht "../application-layer/dht"
	"../application-layer/models"
	"../application-layer/utils"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
)

var (
	dirPath              = filepath.Join("..", "utils")
	marketplaceFilesPath = filepath.Join(dirPath, "marketplaceFiles.json")
	fileMutex            sync.Mutex
	sem                  = make(chan struct{}, 5) //limit to 5 concurrent requests
)

func sendMarketplaceFiles(host host.Host, targetID string) error {
	fmt.Println("sendMarketplaceFiles: sending marketplace files to", targetID)
	fileMutex.Lock()
	defer fileMutex.Unlock()

	// check if file containing all published files exists
	if _, err := os.Stat(marketplaceFilesPath); os.IsNotExist(err) {
		fmt.Printf("sendFile: file %s not found\n", marketplaceFilesPath)
		return err
	}

	fmt.Printf("sending all marketplace files to requester %s \n", targetID)

	// create stream to send the entire json file
	fileStream, err := dht.CreateNewStream(host, targetID, "/marketplaceFiles/p2p")
	if err != nil {
		fmt.Println("Error creating file stream:", err)
		return err
	}
	defer fileStream.Close()

	// open file for reading
	file, err := os.Open(marketplaceFilesPath)
	if err != nil {
		fmt.Printf("error opening file %s: %v", marketplaceFilesPath, err)
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := make([]byte, 4096) // Match the receiver's buffer size

	for {
		n, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Completed sending file data")
				break
			}
			fmt.Printf("Error reading file: %v\n", err)
			return err
		}

		_, err = fileStream.Write(buffer[:n])
		if err != nil {
			fmt.Printf("Error sending data to requester: %v\n", err)
			return err
		}
	}

	return nil
}

// receive and process request for all published nodes
func receiveMarketplaceRequest(node host.Host) error {
	fmt.Println("listening for marketplace requests")
	node.SetStreamHandler("/sendRefreshRequest/p2p", func(s network.Stream) {
		defer s.Close()

		sem <- struct{}{}
		defer func() { <-sem }()

		buf := bufio.NewReader(s)
		data, err := io.ReadAll(buf) //read in requester id

		if err != nil {
			if err == io.EOF {
				log.Printf("Stream closed by peer :%s", s.Conn().RemotePeer())
			} else {
				log.Printf("Error receiving refresh request %v", err)
			}
			return
		}

		sendMarketplaceFiles(node, string(data))
	})
	return nil
}

// used when a file is (re)published or its rating changes
func receiveCloudNodeFiles(node host.Host) error {
	fmt.Println("listening for recently published files")
	node.SetStreamHandler("/nodeFiles/p2p", func(s network.Stream) {
		defer s.Close()

		// Acquire semaphore
		sem <- struct{}{}
		defer func() { <-sem }() // Release semaphore

		buf := bufio.NewReader(s)
		data, err := io.ReadAll(buf) //read in request

		if err != nil {
			if err == io.EOF {
				log.Printf("Stream closed by peer :%s", s.Conn().RemotePeer())
			} else {
				log.Printf("Error receiving refresh request %v", err)
			}
			return
		}

		var fileMetadata models.FileMetadata
		err = json.Unmarshal(data, &fileMetadata)

		utils.SaveOrUpdateFile(fileMetadata, dirPath, marketplaceFilesPath)
	})
	return nil
}

/*
node publishes file to dht -> file metadata is sent to cloud node and stored in json(?) file
marketplace: request for all files from cloud node
	- searching by non file hash??
user upvotes/downvotes -> update in dht and cloud node
refresh cloud node every 24 hours to check if files are still active/published - query dht with GetValue?
	- if not active, mark inactive or delete from file?
	- when users are downloading we check again to make sure providers are actually active


when nodes log in, they automatically republish files to the dht every x hours
*/
