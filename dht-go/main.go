// taken from TA - use as inspo - have to modify

package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	record "github.com/libp2p/go-libp2p-record"

	// new imports

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"

	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/client"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
	"github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multihash"
	"github.com/rs/cors"
)

var (
	node_id             = "114271046" // give your SBU ID
	relay_node_addr     = "/ip4/130.245.173.221/tcp/4001/p2p/12D3KooWDpJ7As7BWAwRMfu1VU2WCqNjvq387JEYKDBj4kx6nXTN"
	bootstrap_node_addr = "/ip4/130.245.173.222/tcp/61000/p2p/12D3KooWQd1K1k8XA9xVEzSAu7HUCodC7LJB6uW5Kw4VwkRdstPE"
	globalCtx           context.Context
	node                host.Host
)

var dirPath = filepath.Join("..", "utils")
var filePath = filepath.Join(dirPath, "files.json")

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
	seed := []byte(node_id)
	customAddr, err := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/0")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse multiaddr: %w", err)
	}
	privKey, err := generatePrivateKeyFromSeed(seed)
	if err != nil {
		log.Fatal(err)
	}
	relayAddr, err := multiaddr.NewMultiaddr(relay_node_addr)
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

func connectToPeer(node host.Host, peerAddr string) {
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
	err = node.Connect(context.Background(), *info)
	if err != nil {
		log.Printf("Failed to connect to peer: %s", err)
		return
	}

	fmt.Println("Connected to:", info.ID)
}

func connectToPeerUsingRelay(node host.Host, targetPeerID string) {
	ctx := globalCtx
	targetPeerID = strings.TrimSpace(targetPeerID)
	relayAddr, err := multiaddr.NewMultiaddr(relay_node_addr)
	if err != nil {
		log.Printf("Failed to create relay multiaddr: %v", err)
	}
	peerMultiaddr := relayAddr.Encapsulate(multiaddr.StringCast("/p2p-circuit/p2p/" + targetPeerID))

	relayedAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMultiaddr)
	if err != nil {
		log.Println("Failed to get relayed AddrInfo: %w", err)
		return
	}
	// Connect to the peer through the relay
	err = node.Connect(ctx, *relayedAddrInfo)
	if err != nil {
		log.Println("Failed to connect to peer through relay: %w", err)
		return
	}

	fmt.Printf("Connected to peer via relay: %s\n", targetPeerID)
}

func receiveDataFromPeer(node host.Host) {
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

func sendDataToPeer(node host.Host, targetpeerid string) {
	var ctx = context.Background()
	targetPeerID := strings.TrimSpace(targetpeerid)
	relayAddr, err := multiaddr.NewMultiaddr(relay_node_addr)
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
	relayInfo, _ := peer.AddrInfoFromString(relay_node_addr)
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
							connectToPeerUsingRelay(node, peerID)
						}
					}
				}
			}
		}
	})
}

func main() {

	node, dht, err := createNode()
	if err != nil {
		log.Fatalf("Failed to create node: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	globalCtx = ctx

	fmt.Println("Node multiaddresses:", node.Addrs())
	fmt.Println("Node Peer ID:", node.ID())
	setupDHT(ctx, dht.Host())
	connectToPeer(node, relay_node_addr) // connect to relay node
	makeReservation(node)                // make reservation on relay node
	go refreshReservation(node, 10*time.Minute)
	connectToPeer(node, bootstrap_node_addr) // connect to bootstrap node
	go handlePeerExchange(node)
	go handleInput(ctx, dht)

	// Set up HTTP handlers
	router := mux.NewRouter()

	//file routes
	router.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		uploadFileHandler(ctx, w, r, dht) // Pass the `dht` instance here
	}).Methods("POST")
	router.HandleFunc("/fetch", func(w http.ResponseWriter, r *http.Request) {
		handleGetProvidersByFileHash(w, r, dht)
	}).Methods("POST")
	router.HandleFunc("/file/", func(w http.ResponseWriter, r *http.Request) {
		getFileHandler(w, r, dht) // Pass the `dht` instance here
	}).Methods("GET")
	router.HandleFunc("/getUploadedFiles", getUploadedFiles)
	router.HandleFunc("/updateFile", func(w http.ResponseWriter, r *http.Request) {
		handleUpdateFile(w, r, dht)
	}).Methods("PUT")

	// Configure CORS
	corsOptions := cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Adjust this to your frontend origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Hash"},
		AllowCredentials: true,
	}

	// CORS handler
	handler := cors.New(corsOptions).Handler(router)

	// Create server with CORS handler
	server := &http.Server{
		Addr:    ":8081",
		Handler: handler,
	}

	// Start the server in a goroutine
	go func() {
		log.Println("Starting server on :8081")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %s\n", err.Error())
		}
	}()

	// Graceful shutdown
	defer func() {
		log.Println("Shutting down server...")
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
		log.Println("Server exited")
	}()

	// Block until a signal is received
	select {}
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
			// provideKey(ctx, dht, key)
			fmt.Println("Record stored successfully")

		case "PUT_PROVIDER":
			if len(args) < 2 {
				fmt.Println("Expected key")
				continue
			}
			key := args[1]
			provideKey(ctx, dht, key)
		default:
			fmt.Println("Expected GET, GET_PROVIDERS, PUT or PUT_PROVIDER")
		}
	}
}

func provideKey(ctx context.Context, dht *dht.IpfsDHT, key string) error {
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
	ctx := globalCtx
	relayInfo, err := peer.AddrInfoFromString(relay_node_addr)
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
		case <-globalCtx.Done():
			fmt.Println("Context done, stopping reservation refresh.")
			return
		}
	}
}

// Handler for retrieving a file
func getFileHandler(w http.ResponseWriter, r *http.Request, dht *dht.IpfsDHT) {
	fmt.Println("in get file handler")
	key := r.URL.Path[len("/file/"):] // Extract the key from the URL
	fmt.Println("after extrating key", key)

	ctx := globalCtx

	res, err := dht.GetValue(ctx, "/orcanet/"+key)
	fmt.Println("result of getValue", res)

	if err != nil {
		http.Error(w, "Failed to retrieve file: "+err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res) // Return the file content
}

// FileMetadata represents the metadata for an uploaded file
type FileMetadata struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Size        int64  `json:"size"`
	Description string `json:"description"`
	Hash        string `json:"hash"`
	IsPublished bool   `json:"isPublished"`
	Fee         int64  `json:"fee"`
}

//var files = make(map[string]FileMetadata) // Store uploaded files metadata by hash

func uploadFileHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, dht *dht.IpfsDHT) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		Size        int64  `json:"size"`
		Description string `json:"description"`
		Hash        string `json:"hash"`
		IsPublished bool   `json:"isPublished"`
		Fee         int64  `json:"fee"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	//ctx := globalCtx

	jsonData, err := json.Marshal(requestBody)
	fmt.Println("Json data is ", string(jsonData))
	if err != nil {
		http.Error(w, "Failed to convert request to JSON", http.StatusBadRequest)
		return
	}

	fmt.Println("after parsing paths, paths are ", dirPath, filePath)

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.Mkdir(dirPath, os.ModePerm); err != nil {
			http.Error(w, "Failed to create utils diretory", http.StatusInternalServerError)
		}
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := os.WriteFile(filePath, []byte("[]"), 0644); err != nil {
			http.Error(w, "Failed to create files.json", http.StatusInternalServerError)
		}
	}

	var fileEntries []map[string]interface{}
	fileData, err := os.ReadFile(filePath)

	if err != nil {
		http.Error(w, "failed to read jiles.json", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(fileData, &fileEntries); err != nil {
		http.Error(w, "failed to parse files.json", http.StatusInternalServerError)
		return
	}

	dup := false
	for _, entry := range fileEntries {
		if entry["hash"] == requestBody.Hash {
			dup = true
			break
		}
	}

	if !dup {
		var newEntry map[string]interface{}
		if err := json.Unmarshal(jsonData, &newEntry); err != nil {
			http.Error(w, "failed to add new entry to files.json", http.StatusInternalServerError)
			return
		}
		fileEntries = append(fileEntries, newEntry)

		updatedData, err := json.MarshalIndent(fileEntries, "", "  ")
		if err != nil {
			http.Error(w, "failed to marshal updated data to json", http.StatusInternalServerError)
			return
		}

		if err := os.WriteFile(filePath, updatedData, 0644); err != nil {
			http.Error(w, "failed to write updated data to files.json", http.StatusInternalServerError)
			return
		}
		fmt.Println("New entry added to files.json")
	} else {
		fmt.Println("duplicate found in files.json")
	}

	err = dht.PutValue(ctx, "/orcanet/"+requestBody.Hash, jsonData)
	if err != nil {
		http.Error(w, "Failed to put inside dht: "+err.Error(), http.StatusInternalServerError)
	}
	fmt.Println("key is ", requestBody.Hash)
	if err != nil {
		http.Error(w, "Failed to publish file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Begin providing ourselves as a provider for that file
	provideKey(ctx, dht, requestBody.Hash)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File published successfully:"))
}

// Create file metadata from the request

// Respond with a success message

// func listKeys(ctx context.Context, dht *dht.IpfsDHT) {
// 	q := query.Query{Prefix: "/orcanet/"} // Adjust prefix to match the namespace used
// 	results, err := dht.Datastore().Query(ctx, q)
// 	if err != nil {
// 		fmt.Printf("Failed to query datastore: %v\n", err)
// 		return
// 	}
// 	defer results.Close()

//		fmt.Println("Available keys:")
//		for result := range results.Next() {
//			if result.Error != nil {
//				fmt.Printf("Error retrieving key: %v\n", result.Error)
//				continue
//			}
//			fmt.Println(result.Key)
//		}
//	}
//
// Helper function to try reverse lookup of peers
func handleGetProvidersByFileHash(w http.ResponseWriter, r *http.Request, dht *dht.IpfsDHT) {
	var requestBody struct {
		Val string `json:"val"` // file hash sent in request
	}

	// Decode request JSON
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// ctx := context.Background()
	key := requestBody.Val // file hash
	data := []byte(key)
	hash := sha256.Sum256(data)
	mh, err := multihash.EncodeName(hash[:], "sha2-256")

	if err != nil {
		fmt.Printf("Error encoding multihash: %v\n", err)
		http.Error(w, "Error encoding hash", http.StatusInternalServerError)
		return
	}

	c := cid.NewCidV1(cid.Raw, mh)
	providers := dht.FindProvidersAsync(globalCtx, c, 20) // asynchronous find

	var providerList []map[string]string
	for p := range providers {
		for _, addr := range p.Addrs {
			providerInfo := map[string]string{
				"peerID":  p.ID.String(),
				"address": addr.String(),
			}
			providerList = append(providerList, providerInfo)
		}
	}

	fmt.Printf("ahhhh provider list: %v\n", providerList)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(providerList)
}

func getUploadedFiles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("tring to fetch user's uploaded files")
	file, err := os.ReadFile("../utils/files.json")

	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	var files []FileMetadata
	if err := json.Unmarshal(file, &files); err != nil {
		http.Error(w, "Failed to parse files data", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(files); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
func updateFileInfo(hash string, newFileData FileMetadata) error {
	fmt.Println("trying to update file info in json", filePath)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %v", err)
	}

	var files []FileMetadata
	err = json.Unmarshal(data, &files)
	if err != nil {
		return fmt.Errorf("error parsing JSON %v", err)
	}
	fmt.Println("data from json", files)
	fmt.Println("new file metadata", newFileData)

	// Update file metadata based on hash
	for i := range files {
		if files[i].Hash == hash {
			fmt.Printf("replacing %v with %v", files[i], newFileData)
			files[i] = newFileData
			break
		}
	}

	// Marshal the updated data back to JSON
	updatedData, err := json.MarshalIndent(files, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// Write the updated data back to the file
	err = os.WriteFile(filePath, updatedData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write updated file: %v", err)
	}

	return nil
}

func handleUpdateFile(w http.ResponseWriter, r *http.Request, dht *dht.IpfsDHT) {
	// Check the Content-Type header to ensure it is application/json
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	// Get the file hash from the headers
	hash := r.Header.Get("Hash") // Make sure it's sent as 'File-Hash'
	if hash == "" {
		http.Error(w, "file hash missing", http.StatusBadRequest)
		return
	}

	// Parse the request body to get the updated file metadata
	var fileData FileMetadata
	if err := json.NewDecoder(r.Body).Decode(&fileData); err != nil {
		http.Error(w, "error decoding metadata", http.StatusBadRequest)
		return
	}

	// Update the file metadata in the JSON file
	err := updateFileInfo(hash, fileData)
	if err != nil {
		http.Error(w, fmt.Sprintf("error updating file info: %v", err), http.StatusInternalServerError)
		return
	}

	// Optionally update the metadata in the DHT (if necessary)
	// dht.UpdateMetadata(hash, fileData) // Example of how you might update metadata in DHT

	// Send a success response back
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File metadata updated successfully"))
}
