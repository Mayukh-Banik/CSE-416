package dht_kad

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multihash"
)

func ConnectToPeer(node host.Host, peerAddr string) {
	addr, err := multiaddr.NewMultiaddr(peerAddr)
	if err != nil {
		log.Printf("Failed to parse peer address: %s", err)
		return
	}

	info, err := peer.AddrInfoFromP2pAddr(addr)
	if err != nil {
		log.Printf("Failed to get AddrInfo from address: %s", err)
		return
	}

	node.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)
	err = node.Connect(GlobalCtx, *info)
	if err != nil {
		log.Printf("Failed to connect to peer: %s", err)
		return
	}

	fmt.Println("Connected to:", info.ID)
}

func ConnectToPeerUsingRelay(node host.Host, targetPeerID string) error {
	ctx := GlobalCtx
	targetPeerID = strings.TrimSpace(targetPeerID)
	relayAddr, err := multiaddr.NewMultiaddr(Relay_node_addr)
	if err != nil {
		log.Printf("Failed to create relay multiaddr: %v", err)
	}
	fmt.Println("--------target peer id:", targetPeerID)
	peerMultiaddr := relayAddr.Encapsulate(multiaddr.StringCast("/p2p-circuit/p2p/" + targetPeerID))

	relayedAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMultiaddr)
	if err != nil {
		return fmt.Errorf("failed to get relayed AddrInfo: %w", err)
	}
	// Connect to the peer through the relay
	err = node.Connect(ctx, *relayedAddrInfo)
	if err != nil {
		return fmt.Errorf("failed to connect to peer through relay: %w", err)
	}

	fmt.Printf("connected to peer via relay: %s\n", targetPeerID)
	return nil
}

func ReceiveDataFromPeer(node host.Host) {
	fmt.Println("listening for data from peer")
	// Set a stream handler to listen for incoming streams on the "/senddata/p2p" protocol
	node.SetStreamHandler("/senddata/p2p", func(s network.Stream) {
		defer s.Close()
		// Create a buffered reader to read data from the stream
		buf := bufio.NewReader(s)
		// Read data from the stream
		data, err := buf.ReadBytes('\n') // Reads until a newline character
		if err != nil {
			if err == io.EOF {
				log.Printf("Stream closed by peer: %s", s.Conn().RemotePeer())
			} else {
				log.Printf("Error reading from stream: %v", err)
			}
			return
		}
		// Print the received data
		log.Printf("Received data: %s", data)
	})
}

func SendDataToPeer(node host.Host, targetpeerid string) {
	fmt.Println("sending data to peer: ", targetpeerid)
	var ctx = context.Background()
	targetPeerID := strings.TrimSpace(targetpeerid)
	relayAddr, err := multiaddr.NewMultiaddr(Relay_node_addr)
	if err != nil {
		log.Printf("Failed to create relay multiaddr: %v", err)
	}
	peerMultiaddr := relayAddr.Encapsulate(multiaddr.StringCast("/p2p-circuit/p2p/" + targetPeerID))

	peerinfo, err := peer.AddrInfoFromP2pAddr(peerMultiaddr)
	if err != nil {
		log.Fatalf("Failed to parse peer address: %s", err)
	}
	if err := node.Connect(ctx, *peerinfo); err != nil {
		log.Printf("Failed to connect to peer %s via relay: %v", peerinfo.ID, err)
		return
	}
	s, err := node.NewStream(network.WithAllowLimitedConn(ctx, "/senddata/p2p"), peerinfo.ID, "/senddata/p2p")
	if err != nil {
		log.Printf("Failed to open stream to %s: %s", peerinfo.ID, err)
		return
	}
	defer s.Close()

	_, err = s.Write([]byte("sending hello to peer\n"))
	if err != nil {
		log.Fatalf("Failed to write to stream: %s", err)
	}
}

func handlePeerExchange(node host.Host) {
	relayInfo, _ := peer.AddrInfoFromString(Relay_node_addr)
	node.SetStreamHandler("/orcanet/p2p", func(s network.Stream) {
		defer s.Close()

		buf := bufio.NewReader(s)
		peerAddr, err := buf.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Printf("error reading from stream: %v", err)
			}
		}
		peerAddr = strings.TrimSpace(peerAddr)
		var data map[string]interface{}
		err = json.Unmarshal([]byte(peerAddr), &data)
		if err != nil {
			fmt.Printf("error unmarshaling JSON: %v", err)
		}
		if knownPeers, ok := data["known_peers"].([]interface{}); ok {
			for _, peer := range knownPeers {
				fmt.Println("Peer:")
				if peerMap, ok := peer.(map[string]interface{}); ok {
					if peerID, ok := peerMap["peer_id"].(string); ok {
						if string(peerID) != string(relayInfo.ID) {
							ConnectToPeerUsingRelay(node, peerID)
						}
					}
				}
			}
		}
	})
}

// find all providers for a file
func FindProviders(fileHash string) {
	fmt.Printf("looking for providers for file: %v", fileHash)
	data := []byte(fileHash)
	hash := sha256.Sum256(data)
	mh, err := multihash.EncodeName(hash[:], "sha2-256")
	if err != nil {
		fmt.Printf("Error encoding multihash: %v\n", err)
	}
	c := cid.NewCidV1(cid.Raw, mh)
	providers := DHT.FindProvidersAsync(GlobalCtx, c, 20)

	fmt.Println("Searching for providers...")
	for p := range providers {
		if p.ID == peer.ID("") {
			break
		}
		fmt.Printf("Found provider: %s\n", p.ID.String())
		for _, addr := range p.Addrs {
			fmt.Printf(" - Address: %s\n", addr.String())
		}
	}
}

// get addr for a specific provider of file
func FindSpecificProvider(fileHash string, targetProviderID peer.ID) (*peer.AddrInfo, error) {
	fmt.Printf("looking for providers for file: %v", fileHash)
	data := []byte(fileHash)
	hash := sha256.Sum256(data)

	// Encode to multihash
	mh, err := multihash.EncodeName(hash[:], "sha2-256")
	if err != nil {
		return nil, fmt.Errorf("error encoding multihash: %v", err)
	}

	// Create CID from multihash
	c := cid.NewCidV1(cid.Raw, mh)

	// Start asynchronous provider search
	providers := DHT.FindProvidersAsync(GlobalCtx, c, 20)
	targetPeerID := peer.ID(targetProviderID)

	fmt.Println("Searching for specific provider...")
	for p := range providers {
		if p.ID == targetPeerID {
			fmt.Printf("Found target provider: %s\n", p.ID.String())
			for _, addr := range p.Addrs {
				fmt.Printf(" - Address: %s\n", addr.String())
			}
			// Return the matching provider's AddrInfo
			return &p, nil
		}
	}

	return nil, fmt.Errorf("provider with ID %s not found for file hash %s", targetProviderID, fileHash)
}

// adapted from sendDataToPeer
func CreateNewStream(node host.Host, targetPeerID string, streamProtocol protocol.ID) (network.Stream, error) {
	// Use a timeout context for the stream connection attempt
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Parse the target peer ID and set up the relay multiaddr
	targetPeerID = strings.TrimSpace(targetPeerID)
	relayNodeAddr, err := multiaddr.NewMultiaddr(Relay_node_addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create relay node multiaddr: %v", err)
	}

	// Create a multiaddress that encapsulates the relay address with the target peer ID
	peerMultiaddr := relayNodeAddr.Encapsulate(multiaddr.StringCast("/p2p-circuit/p2p/" + targetPeerID))
	peerinfo, err := peer.AddrInfoFromP2pAddr(peerMultiaddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get target peer address: %s", err)
	}

	// Connect to the target peer through the relay node
	if err := node.Connect(ctx, *peerinfo); err != nil {
		return nil, fmt.Errorf("failed to connect to peer %s via relay: %v", peerinfo.ID, err)
	}

	// Open a new stream to the target peer
	stream, err := node.NewStream(ctx, peerinfo.ID, streamProtocol)
	if err != nil {
		return nil, fmt.Errorf("failed to open stream to target peer %s: %s", peerinfo.ID, err)
	}

	fmt.Printf("Successfully created stream to peer %s\n", peerinfo.ID)
	return stream, nil
}
