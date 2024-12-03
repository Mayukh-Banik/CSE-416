package dht_kad

import (
	"application-layer/models"
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-kad-dht/providers"
	record "github.com/libp2p/go-libp2p-record"

	// new imports

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/client"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
	"github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multihash"
)

var (
	Node_id         string
	Relay_node_addr = "/ip4/130.245.173.221/tcp/4001/p2p/12D3KooWDpJ7As7BWAwRMfu1VU2WCqNjvq387JEYKDBj4kx6nXTN"
	// Bootstrap_node_addr = "/ip4/130.245.173.222/tcp/61000/p2p/12D3KooWQd1K1k8XA9xVEzSAu7HUCodC7LJB6uW5Kw4VwkRdstPE"
	// Bootstrap_node_addr = "/ip4/192.168.86.218/tcp/61000/p2p/12D3KooWPs4FtjU4YmGoFgnd225gj3XKBD6QZpWFK5Pq1yEp87kx"
	Bootstrap_node_addr = "/ip4/192.168.1.169/tcp/61000/p2p/12D3KooWAZv5dC3xtzos2KiJm2wDqiLGJ5y4gwC7WSKU5DvmCLEL"
	GlobalCtx           context.Context
	PeerID              string
	DHT                 *dht.IpfsDHT
	ProviderStore       providers.ProviderStore
	Host                host.Host
)

func generatePrivateKeyFromSeed(seed []byte) (crypto.PrivKey, error) {
	hash := sha256.Sum256(seed) // Generate deterministic key material
	// Create an Ed25519 private key from the hash
	privKey, _, err := crypto.GenerateEd25519Key(
		bytes.NewReader(hash[:]),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}
	return privKey, nil
}

func setupDHT(ctx context.Context, h host.Host) *dht.IpfsDHT {
	// Set up the DHT instance
	kadDHT, err := dht.New(ctx, h, dht.Mode(dht.ModeClient))
	if err != nil {
		log.Fatal(err)
	}

	// Bootstrap the DHT (connect to other peers to join the DHT network)
	err = kadDHT.Bootstrap(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Configure the DHT to use the custom validator - idk if this even goes here
	kadDHT.Validator = record.NamespacedValidator{
		"orcanet": &CustomValidator{}, // Add a custom validator for the "orcanet" namespace
	}

	return kadDHT
}

func createNode() (host.Host, *dht.IpfsDHT, error) {
	ctx := context.Background()
	seed := []byte(Node_id)
	customAddr, err := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/0")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse multiaddr: %w", err)
	}
	privKey, err := generatePrivateKeyFromSeed(seed)
	if err != nil {
		log.Fatal(err)
	}
	relayAddr, err := multiaddr.NewMultiaddr(Relay_node_addr)
	if err != nil {
		log.Fatalf("Failed to create relay multiaddr: %v", err)
	}

	// Convert the relay multiaddress to AddrInfo
	relayInfo, err := peer.AddrInfoFromP2pAddr(relayAddr)
	if err != nil {
		log.Fatalf("Failed to create AddrInfo from relay multiaddr: %v", err)
	}

	node, err := libp2p.New(
		libp2p.ListenAddrs(customAddr),
		libp2p.Identity(privKey),
		libp2p.NATPortMap(),
		libp2p.EnableNATService(),
		libp2p.EnableAutoRelayWithStaticRelays([]peer.AddrInfo{*relayInfo}),
		libp2p.EnableRelayService(),
		libp2p.EnableHolePunching(),
	)

	if err != nil {
		return nil, nil, err
	}
	_, err = relay.New(node)
	if err != nil {
		log.Printf("Failed to instantiate the relay: %v", err)
	}

	dhtRouting, err := dht.New(ctx, node, dht.Mode(dht.ModeClient))
	if err != nil {
		return nil, nil, err
	}
	namespacedValidator := record.NamespacedValidator{
		"orcanet": &CustomValidator{}, // Add a custom validator for the "orcanet" namespace
	}

	dhtRouting.Validator = namespacedValidator // Configure the DHT to use the custom validator

	err = dhtRouting.Bootstrap(ctx)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("DHT bootstrap complete.")

	// Set up notifications for new connections
	node.Network().Notify(&network.NotifyBundle{
		ConnectedF: func(n network.Network, conn network.Conn) {
			fmt.Printf("Notification: New peer connected %s\n", conn.RemotePeer().String())
		},
	})

	return node, dhtRouting, nil
}

func handleInput(ctx context.Context, dht *dht.IpfsDHT) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("User Input \n ")
	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n') // Read input from keyboard
		input = strings.TrimSpace(input)    // Trim any trailing newline or spaces
		args := strings.Split(input, " ")
		if len(args) < 1 {
			fmt.Println("No command provided")
			continue
		}
		command := args[0]
		command = strings.ToUpper(command)
		switch command {
		// case "LIST_ALL":
		// 	if len(args) >= 2 {
		// 		fmt.Println("We don't want arguments")
		// 		continue
		// 	}
		// 	listKeys(ctx, dht)
		case "GET":
			if len(args) < 2 {
				fmt.Println("Expected key")
				continue
			}
			key := args[1]
			dhtKey := "/orcanet/" + key
			res, err := dht.GetValue(ctx, dhtKey)
			if err != nil {
				fmt.Printf("Failed to get record: %v\n", err)
				continue
			}
			fmt.Printf("Record: %s\n", res)

		case "GET_PROVIDERS":
			if len(args) < 2 {
				fmt.Println("Expected key")
				continue
			}
			key := args[1]
			data := []byte(key)
			hash := sha256.Sum256(data)
			mh, err := multihash.EncodeName(hash[:], "sha2-256")
			if err != nil {
				fmt.Printf("Error encoding multihash: %v\n", err)
				continue
			}
			c := cid.NewCidV1(cid.Raw, mh)
			providers := dht.FindProvidersAsync(ctx, c, 20)

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

		case "PUT":
			if len(args) < 3 {
				fmt.Println("Expected key and value")
				continue
			}
			key := args[1]
			value := args[2]
			dhtKey := "/orcanet/" + key
			log.Println(dhtKey)
			err := dht.PutValue(ctx, dhtKey, []byte(value))
			if err != nil {
				fmt.Printf("Failed to put record: %v\n", err)
				continue
			}
			// ProvideKey(ctx, dht, key)
			fmt.Println("Record stored successfully")

		case "PUT_PROVIDER":
			if len(args) < 2 {
				fmt.Println("Expected key")
				continue
			}
			key := args[1]
			ProvideKey(ctx, dht, key)

			// doesnt work
		// case "GET_FILE":
		// 	if len(args) < 2 {
		// 		fmt.Println("Expected file request")
		// 		continue
		// 	}
		// 	fmt.Println("download request", args[1])
		// 	jsonInput := os.Args[1]

		// 	// Create a Transaction struct
		// 	var transaction models.Transaction

		// 	// Parse the JSON string into the struct
		// 	err := json.Unmarshal([]byte(jsonInput), &transaction)
		// 	if err != nil {
		// 		fmt.Println("Error parsing JSON:", err)
		// 		return
		// 	}
		// 	SendDownloadRequest(transaction)

		default:
			fmt.Println("Expected GET, GET_PROVIDERS, PUT or PUT_PROVIDER")
		}
	}
}

func ProvideKey(ctx context.Context, dht *dht.IpfsDHT, key string) error {
	data := []byte(key)
	hash := sha256.Sum256(data)
	mh, err := multihash.EncodeName(hash[:], "sha2-256")
	if err != nil {
		return fmt.Errorf("error encoding multihash: %v", err)
	}
	c := cid.NewCidV1(cid.Raw, mh)

	// Start providing the key
	err = dht.Provide(ctx, c, true)
	if err != nil {
		return fmt.Errorf("failed to start providing key: %v", err)
	}
	return nil
}

func makeReservation(node host.Host) {
	ctx := GlobalCtx
	relayInfo, err := peer.AddrInfoFromString(Relay_node_addr)
	if err != nil {
		log.Fatalf("Failed to create addrInfo from string representation of relay multiaddr: %v", err)
	}
	_, err = client.Reserve(ctx, node, *relayInfo)
	if err != nil {
		log.Fatalf("Failed to make reservation on relay: %v", err)
	}
	fmt.Printf("Reservation successfull \n")
}

func refreshReservation(node host.Host, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			makeReservation(node)
		case <-GlobalCtx.Done():
			fmt.Println("Context done, stopping reservation refresh.")
			return
		}
	}
}

func getNodeId() {
	for {
		fmt.Print("Enter SBU ID: ")
		reader := bufio.NewReader(os.Stdin)
		id, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue // Retry on error
		}

		id = strings.TrimSpace(id)
		if id == "" {
			fmt.Println("SBU ID cannot be empty. Please try again.")
			continue // Retry on empty input
		}

		Node_id = id
		fmt.Println("Your SBUID is:", Node_id)
		break // Exit the loop on valid input
	}
}

// move to files package?
func UpdateFileInDHT(currentInfo models.FileMetadata) error {
	// Retrieve the current metadata for the file, if it exists
	var currentMetadata models.DHTMetadata
	existingData, err := DHT.GetValue(GlobalCtx, "/orcanet/"+currentInfo.Hash)
	if err == nil { // If data exists, unmarshal it
		err = json.Unmarshal(existingData, &currentMetadata)
		if err != nil {
			log.Fatal("Failed to unmarshal existing DHTMetadata:", err)
		}
	} else {
		// If no existing metadata, initialize a new DHTMetadata
		currentMetadata = models.DHTMetadata{
			Name:              currentInfo.Name,
			Type:              currentInfo.Type,
			Size:              currentInfo.Size,
			Description:       currentInfo.Description,
			CreatedAt:         currentInfo.CreatedAt,
			NameWithExtension: currentInfo.NameWithExtension,
			Rating:            0,
			Hash:              currentInfo.Hash,
		}
	}

	currentMetadata.Providers = make(map[string]models.Provider)
	// Add the new provider to the list of current providers
	provider := models.Provider{
		PeerAddr: DHT.Host().Addrs()[0].String(),
		IsActive: currentInfo.IsPublished,
		Fee:      currentInfo.Fee,
		// leave Rating empty -> no vote
		// this will be updated when handling upvoting/downvoting
	}

	// Check if the provider already exists in the metadata by PeerID
	if existingProvider, exists := currentMetadata.Providers[PeerID]; exists {
		// Update the IsActive field
		existingProvider.IsActive = provider.IsActive
		// Update the provider in the map
		currentMetadata.Providers[PeerID] = existingProvider
	} else {
		// If provider does not exist, add the new provider
		currentMetadata.Providers[PeerID] = provider
	}

	// Marshal the updated metadata
	dhtMetadataBytes, err := json.Marshal(currentMetadata)
	if err != nil {
		fmt.Errorf("failed to marshal updated DHTMetadata: %w\n", err)
	}

	// Begin providing ourselves as a provider for that file
	err = ProvideKey(GlobalCtx, DHT, currentInfo.Hash)
	if err != nil {
		fmt.Errorf("failed to register updated file to dht: %w\n", err)
	}

	// Store the updated metadata in the DHT
	err = DHT.PutValue(GlobalCtx, "/orcanet/"+currentInfo.Hash, dhtMetadataBytes)
	if err != nil {
		fmt.Errorf("failed to updated file in dht: %w\n", err)
	}
	fmt.Println("successfully updated file to dht with new provider", currentInfo.Hash)

	return nil
}

//////////////////////////////////////////
// for testing - TA changed bootstrap stuff
//////////////////////////////////////////

// var (
// 	globalCtx      context.Context
// 	relay_addr     = "/ip4/130.245.173.221/tcp/4001/p2p/12D3KooWDpJ7As7BWAwRMfu1VU2WCqNjvq387JEYKDBj4kx6nXTN"
// 	bootstrap_seed = ":("
// )

// func createNode() (host.Host, *dht.IpfsDHT, error) {
// 	ctx := context.Background()
// 	globalCtx = ctx

// 	customAddr, err := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/61000")
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("failed to parse multiaddr: %w", err)
// 	}

// 	seed := []byte(bootstrap_seed)
// 	privKey, err := generatePrivateKeyFromSeed(seed)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	node, err := libp2p.New(
// 		libp2p.ListenAddrs(customAddr),
// 		libp2p.Identity(privKey),
// 		libp2p.NATPortMap(),
// 		libp2p.EnableNATService(),
// 		// libp2p.EnableAutoRelay(),
// 		// libp2p.StaticRelays(staticRelays),
// 	)

// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	_, err = relay.New(node)
// 	if err != nil {
// 		log.Printf("Failed to instantiate the relay: %v", err)
// 	}

// 	dhtRouting, err := dht.New(ctx, node, dht.Mode(dht.ModeServer))
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	err = dhtRouting.Bootstrap(ctx)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	namespacedValidator := record.NamespacedValidator{
// 		"orcanet": &CustomValidator{}, // Add a custom validator for the "orcanet" namespace
// 	}
// 	// Configure the DHT to use the custom validator
// 	dhtRouting.Validator = namespacedValidator

// 	// Set up notifications for new connections
// 	node.Network().Notify(&network.NotifyBundle{
// 		ConnectedF: func(n network.Network, conn network.Conn) {
// 			go exchangePeers(node, conn.RemotePeer())
// 			// fmt.Printf("New peer connected: %s\n", conn.RemotePeer().String())
// 			// fmt.Println("peers in network", node.Network().Peers())
// 		},
// 	})

// 	return node, dhtRouting, nil
// }

// func exchangePeers(node host.Host, newPeer peer.ID) {
// 	knownPeers := node.Network().Peers()
// 	var peerInfos []string
// 	data := map[string]interface{}{
// 		"known_peers": []map[string]string{},
// 	}
// 	var temp map[string]string
// 	relay_info, _ := peer.AddrInfoFromString(relay_addr)
// 	for _, peer := range knownPeers {
// 		if peer != newPeer && peer != node.ID() && peer != relay_info.ID {
// 			temp = make(map[string]string)
// 			temp["peer_id"] = peer.String()
// 			peerInfos = append(peerInfos, peer.String())
// 			data["known_peers"] = append(data["known_peers"].([]map[string]string), temp)
// 		}
// 	}

// 	s, err := node.NewStream(network.WithAllowLimitedConn(globalCtx, "/orcanet/p2p"), newPeer, "/orcanet/p2p")
// 	if err != nil {
// 		//log.Printf("Failed to open stream to %s: %s", newPeer, err)
// 		return
// 	}
// 	defer s.Close()
// 	jsonData, err := json.Marshal(data)
// 	if err != nil {
// 		log.Fatalf("Error marshaling map to JSON: %s", err)
// 	}
// 	s.Write([]byte(jsonData))

// 	fmt.Printf("Shared %d peers with %s\n", len(peerInfos), newPeer.String())
// }
