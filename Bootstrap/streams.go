package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-multiaddr"
)

var (
	DirPath              = filepath.Join("..", "utils")
	MarketplaceFilesPath = filepath.Join(DirPath, "marketplaceFiles.json")
	sem                  = make(chan struct{}, 5) //limit to 5 concurrent requests
)

func sendMarketplaceFiles(host host.Host, targetID string) error {
	fmt.Println("sendMarketplaceFiles: sending marketplace files to", targetID)
	fileMutex.Lock()
	defer fileMutex.Unlock()

	// check if file containing all published files exists
	if _, err := os.Stat(MarketplaceFilesPath); os.IsNotExist(err) {
		fmt.Printf("sendFile: file %s not found\n", MarketplaceFilesPath)
		return err
	}

	fmt.Printf("sending all marketplace files to requester %s \n", targetID)

	// create stream to send the entire json file
	fileStream, err := createNewStream(host, targetID, "/marketplaceFiles/p2p")
	if err != nil {
		fmt.Println("Error creating file stream:", err)
		return err
	}
	defer fileStream.Close()

	// open file for reading
	file, err := os.Open(MarketplaceFilesPath)
	if err != nil {
		fmt.Printf("error opening file %s: %v", MarketplaceFilesPath, err)
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
		fmt.Println("receiveMarketplaceRequest: received request from", string(data))
		if err != nil {
			if err == io.EOF {
				log.Printf("Stream closed by peer :%s", s.Conn().RemotePeer())
			} else {
				log.Printf("Error receiving refresh request %v", err)
			}
			return
		}
		fmt.Println("receiveMarketplaceRequest: now entering sendMarketplaceFiles")
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

		// var fileMetadata FileMetadata
		var metadata DHTMetadata
		err = json.Unmarshal(data, &metadata)
		if err != nil {
			fmt.Printf("receiveCloudNodeFiles: failed to unmarshal file metadata: %v", err)
			return
		}

		// update json file containing all dht files - and corresponding variables
		updateFile(metadata, DirPath, MarketplaceFilesPath, false)
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

// adapted from sendDataToPeer
func createNewStream(node host.Host, targetPeerID string, streamProtocol protocol.ID) (network.Stream, error) {
	fmt.Printf("CreateNewStream %v: sending data to peer %v\n", streamProtocol, targetPeerID)

	// Create a context for connection
	var ctx = context.Background()
	targetPeerID = strings.TrimSpace(targetPeerID)

	// Create the relay address
	relayAddr, err := multiaddr.NewMultiaddr(relay_addr)
	if err != nil {
		log.Printf("Failed to create relay multiaddr: %v", err)
		return nil, fmt.Errorf("failed to create relay multiaddr: %v", err)
	}

	// Encapsulate the relay address with the target peer's address
	peerMultiaddr := relayAddr.Encapsulate(multiaddr.StringCast("/p2p-circuit/p2p/" + targetPeerID))

	// Parse the multiaddress into a PeerInfo object
	peerinfo, err := peer.AddrInfoFromP2pAddr(peerMultiaddr)
	if err != nil {
		log.Printf("Failed to parse peer address: %s", err)
		return nil, fmt.Errorf("failed to parse peer address: %v", err)
	}

	// Connect to the target peer
	if err := node.Connect(ctx, *peerinfo); err != nil {
		log.Printf("Failed to connect to peer %s via relay: %v", peerinfo.ID, err)
		return nil, fmt.Errorf("failed to connect to peer %s via relay: %v", peerinfo.ID, err)
	}
	fmt.Printf("connected to node %v, now creating stream %v", targetPeerID, streamProtocol)

	// Create a new stream to the target peer
	// stream, err := node.NewStream(ctx, peerinfo.ID, streamProtocol)
	stream, err := node.NewStream(network.WithAllowLimitedConn(ctx, string(streamProtocol)), peerinfo.ID, streamProtocol)

	if err != nil {
		log.Printf("Failed to open stream to %s: %s", peerinfo.ID, err)
		return nil, fmt.Errorf("failed to open stream to peer %s: %v", peerinfo.ID, err)
	}

	fmt.Printf("Successfully created stream to peer %s\n", peerinfo.ID)
	return stream, nil
}
