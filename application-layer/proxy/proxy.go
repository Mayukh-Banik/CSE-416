package proxyService

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

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/client"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
	ma "github.com/multiformats/go-multiaddr"
)

var (
	node_id         = ""
	peer_id         = ""
	relay_node_addr = "/ip4/130.245.173.221/tcp/4001/p2p/12D3KooWDpJ7As7BWAwRMfu1VU2WCqNjvq387JEYKDBj4kx6nXTN"
	globalCtx       context.Context
	Peer_Addresses  []ma.Multiaddr
	isHost          = true
)

type Proxy struct {
	Name       string   `json:"name"`
	Location   string   `json:"location"`
	Logs       []string `json:"logs"`
	Statistics struct {
		Uptime string `json:"uptime"`
	}
	Bandwidth string `json:"bandwidth"`
	Address   string `json:"address"`
	PeerID    string `json:"peer_id"`
	IsEnabled bool   `json:"isEnabled"`

	Price string `json:"price"`
}

const proxyFilePath = "../utils/proxy_data.json"

func saveProxyToFile(proxy Proxy) error {
	dir := filepath.Dir(proxyFilePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("could not create directory: %v", err)
		}
	}

	// Open file for reading and writing
	file, err := os.OpenFile(proxyFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	// Read existing proxies from the file
	var proxies []Proxy
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&proxies); err != nil && err.Error() != "EOF" {
		return fmt.Errorf("could not decode existing proxy data: %v", err)
	}

	if proxies == nil {
		proxies = []Proxy{} // Initialize an empty slice if no proxies exist
	}

	// Check if the proxy already exists based on PeerID (or another unique identifier)
	for _, existingProxy := range proxies {
		if existingProxy.PeerID == proxy.PeerID {
			fmt.Println("Proxy already exists in file. Skipping append.")
			return nil // Skip if the proxy already exists
		}
	}

	// Append the new proxy
	proxies = append(proxies, proxy)

	// Rewind the file to the beginning and overwrite it with the updated data
	file.Seek(0, 0)
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Optional: for better readability
	if err := encoder.Encode(proxies); err != nil {
		return fmt.Errorf("could not encode proxy data: %v", err)
	}

	fmt.Println("Proxy saved successfully.")
	return nil
}

func loadProxyFromFile() (Proxy, error) {
	var proxy Proxy
	file, err := os.Open(proxyFilePath)
	if err != nil {
		return proxy, fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&proxy); err != nil {
		return proxy, fmt.Errorf("could not decode data: %v", err)
	}
	return proxy, nil
}
func addProxy(address string, peerID string) {
	proxy := Proxy{Address: address, PeerID: peerID}
	fmt.Print("attempting to add proxy", proxy)
	err := saveProxyToFile(proxy)

	if err != nil {
		fmt.Println("Error saving proxy:", err)
	} else {
		fmt.Println("Proxy saved successfully.")
	}
}

/*
Makes private key from a seed, right now the program uses command line args
to make the seed, use file path later on
*/
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

/*
Creates Basic Node for HTTP Proxy, have to set stream handlers and such
*/
func createNode() (host.Host, error) {
	seed := []byte(node_id)
	customAddr, err := ma.NewMultiaddr("/ip4/0.0.0.0/tcp/0")
	if err != nil {
		return nil, fmt.Errorf("failed to parse multiaddr: %w", err)
	}
	privKey, err := generatePrivateKeyFromSeed(seed)
	if err != nil {
		log.Fatal(err)
	}
	relayAddr, err := ma.NewMultiaddr(relay_node_addr)
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
		return nil, err
	}
	_, err = relay.New(node)
	if err != nil {
		log.Printf("Failed to instantiate the relay: %v", err)
	}

	return node, nil
}

func connectToPeer(node host.Host, peerAddr string) {
	addr, err := ma.NewMultiaddr(peerAddr)
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
	relayAddr, err := ma.NewMultiaddr(relay_node_addr)
	if err != nil {
		log.Printf("Failed to create relay multiaddr: %v", err)
	}
	peerMultiaddr := relayAddr.Encapsulate(ma.StringCast("/p2p-circuit/p2p/" + targetPeerID))

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

		// Split the received data into substrings using space as delimiter
		cleanData := strings.ReplaceAll(string(data), "\n", "")
		parts := strings.Split(string(cleanData), "\"")
		var partTemp []string

		for _, part := range parts {
			if len(part) == 0 {
			} else {
				partTemp = append(partTemp, part)
			}
		}

		if isHost {
			for _, part := range partTemp {
				temp, err := ma.NewMultiaddr(part)
				if err != nil {
					fmt.Println("Error in making multiaddress from sent data")
					continue
				}
				Peer_Addresses = append(Peer_Addresses, temp)
			}
			addr, _ := peer.AddrInfoFromP2pAddr(Peer_Addresses[0])
			sendDataToPeer(node, addr.ID.String())
			peer_id = addr.ID.String()

			for _, b := range Peer_Addresses {
				// connectToPeerDirect(node, b.String())
				fmt.Println(b.String())
			}
			connectToPeer(node, Peer_Addresses[0].String())
		} else {
			for _, part := range partTemp {
				temp, err := ma.NewMultiaddr(part)
				if err != nil {
					fmt.Println("Error in making multiaddress from sent data")
					continue
				}
				Peer_Addresses = append(Peer_Addresses, temp)
			}
			for _, b := range Peer_Addresses {
				// connectToPeerDirect(node, b.String())
				fmt.Println(b.String())
			}
			connectToPeer(node, Peer_Addresses[0].String())
		}
	})
}

func sendDataToPeer(node host.Host, targetpeerid string) {
	var ctx = context.Background()
	targetPeerID := strings.TrimSpace(targetpeerid)
	relayAddr, err := ma.NewMultiaddr(relay_node_addr)
	if err != nil {
		log.Printf("Failed to create relay multiaddr: %v", err)
	}
	peerMultiaddr := relayAddr.Encapsulate(ma.StringCast("/p2p-circuit/p2p/" + targetPeerID))

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

	var addrString string
	// Create the string to send with node ID followed by multiaddresses
	for _, a := range node.Addrs() {
		addrString += fmt.Sprintf("\"%s/p2p/%s\"", a.String(), node.ID().String())
	}

	addrString += "\n"

	// Write the address string to the stream
	_, err = s.Write([]byte(addrString))
	if err != nil {
		log.Printf("Failed to write to stream: %s", err)
		return
	}

	fmt.Printf("Sent multiaddresses to peer %s: %s\n", targetPeerID, addrString)
	defer s.Close()
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

func pollPeerAddresses(node host.Host) {
	if isHost {
		fmt.Println("In host part")
		httpHostToClient(node)
	}
	for {
		if len(Peer_Addresses) > 0 {
			fmt.Println("Peer addresses are not empty, setting up proxy.")
			fmt.Println(Peer_Addresses)

			// Hosting
			if !(isHost) {
				fmt.Println("In Client part")

				// Start an HTTP server on port 9900
				http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
					// Send the received HTTP request to the peer
					clientHTTPRequestToHost(node, peer_id, r, w)
				})

				// Listen on port 9900
				serverAddr := ":8081"
				fmt.Printf("HTTP server listening on %s\n", serverAddr)
				if err := http.ListenAndServe(serverAddr, nil); err != nil {
					log.Fatalf("Failed to start HTTP server: %s", err)
				}
			}
		}
		time.Sleep(5 * time.Second)
	}
}

// Retrieveing proxies data, and adding yourself as host
func handleProxyData(w http.ResponseWriter, r *http.Request) {
	// Create a new node
	node, err := createNode()
	if err != nil {
		log.Fatalf("Failed to create node: %s", err)
	}
	fmt.Println("This node's addresses", node.Addrs())

	globalCtx = context.Background()

	// Log the node's details
	fmt.Println("Node multiaddresses:", node.Addrs())
	fmt.Println("Node Peer ID:", node.ID())

	if r.Method == "POST" {
		isHost = true
		// Save the proxy information (address and peer ID) to file
		var newProxy Proxy
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&newProxy)
		newProxy.Address = node.Addrs()[0].String() // Use the first address
		newProxy.PeerID = node.ID().String()
		err = saveProxyToFile(newProxy)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to save proxy: %v", err), http.StatusInternalServerError)
			return
		}
		fmt.Println("Proxy saved successfully.")
	}

	// Poll for peer addresses (optional)
	go pollPeerAddresses(node)
	if r.Method == "GET" {
		// Open the file that stores the proxy data
		file, err := os.Open(proxyFilePath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not open file: %v", err), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Initialize a slice to hold the proxy data
		var proxies []Proxy

		// Decode the JSON data from the file into the proxies slice
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&proxies)
		if err != nil && err.Error() != "EOF" {
			http.Error(w, fmt.Sprintf("Could not decode proxy data: %v", err), http.StatusInternalServerError)
			return
		}

		// Set the content type to JSON
		w.Header().Set("Content-Type", "application/json")

		// Marshal the proxy data into JSON and send it as the response
		jsonResponse, err := json.Marshal(proxies)
		if err != nil {
			http.Error(w, "Failed to generate response", http.StatusInternalServerError)
			return
		}

		// Write the JSON response
		w.Write(jsonResponse)
	}
}

func main() {

	// switch len(os.Args) {
	// case 1:
	// 	fmt.Println("Error: Missing required arguments.")
	// 	os.Exit(1)
	// case 2:
	// 	node_id = os.Args[1]
	// case 3:
	// 	node_id = os.Args[1]
	// 	peer_id = os.Args[2]
	// 	isHost = false
	// default:
	// 	fmt.Println("Error: Too many arguments provided.")
	// 	os.Exit(1)
	// }

	node, err := createNode()
	if err != nil {
		log.Fatalf("Failed to create node: %s", err)
	}
	fmt.Println("This node's addresses", node.Addrs())

	globalCtx = context.Background()

	fmt.Println("Node multiaddresses:", node.Addrs())
	fmt.Println("Node Peer ID:", node.ID())

	connectToPeer(node, relay_node_addr) // connect to relay node
	makeReservation(node)                // make reservation on realy node
	go refreshReservation(node, 10*time.Minute)

	go handlePeerExchange(node)

	receiveDataFromPeer(node)
	if len(os.Args) == 3 {
		sendDataToPeer(node, peer_id)
	}

	defer node.Close()

	pollPeerAddresses(node)

	select {}
}

/*
Makes a reservation on TA's relay node, don't edit
*/
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

/*
Refreshes reservation on TA's node, don't edit, run as separate thread
*/
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

func clientHTTPRequestToHost(node host.Host, targetpeerid string, req *http.Request, w http.ResponseWriter) {
	fmt.Println("In client http request to host")
	var ctx = context.Background()
	targetPeerID := strings.TrimSpace(targetpeerid)
	relayAddr, err := ma.NewMultiaddr(relay_node_addr)
	if err != nil {
		log.Printf("Failed to create relay multiaddr: %v", err)
		return
	}
	peerMultiaddr := relayAddr.Encapsulate(ma.StringCast("/p2p-circuit/p2p/" + targetPeerID))

	peerinfo, err := peer.AddrInfoFromP2pAddr(peerMultiaddr)
	if err != nil {
		log.Fatalf("Failed to parse peer address: %s", err)
		return
	}
	if err := node.Connect(ctx, *peerinfo); err != nil {
		log.Printf("Failed to connect to peer %s via relay: %v", peerinfo.ID, err)
		return
	}

	s, err := node.NewStream(network.WithAllowLimitedConn(ctx, "/http-temp-protocol"), peerinfo.ID, "/http-temp-protocol")
	if err != nil {
		log.Printf("Failed to open stream to %s: %s", peerinfo.ID, err)
		return
	}
	defer s.Close()

	// Serialize the HTTP request
	var buf bytes.Buffer
	if err := req.Write(&buf); err != nil {
		log.Printf("Failed to serialize HTTP request: %v", err)
		return
	}
	httpData := buf.Bytes()

	// Write serialized HTTP request to the stream
	_, err = s.Write(httpData)
	if err != nil {
		log.Printf("Failed to write to stream: %s", err)
		return
	}
	s.CloseWrite() // Close the write side to signal EOF
	fmt.Printf("Sent HTTP request to peer %s: %s\n", targetPeerID, req.URL.String())

	// Wait for a response
	responseBuf := new(bytes.Buffer)
	_, err = responseBuf.ReadFrom(s)
	if err != nil {
		log.Printf("Failed to read response from stream: %v", err)
		return
	}

	// Parse the HTTP response
	responseReader := bufio.NewReader(responseBuf) // Wrap the buffer in a bufio.Reader
	resp, err := http.ReadResponse(responseReader, req)
	if err != nil {
		log.Printf("Failed to parse HTTP response: %v", err)
		return
	}

	// Write the response to the ResponseWriter
	for k, v := range resp.Header {
		for _, h := range v {
			w.Header().Add(k, h)
		}
	}
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Failed to relay response body: %v", err)
	}
}

func httpHostToClient(node host.Host) {
	node.SetStreamHandler("/http-temp-protocol", func(s network.Stream) {
		fmt.Println("In host to client")
		defer s.Close()

		buf := bufio.NewReader(s)
		// Read the HTTP request from the stream
		req, err := http.ReadRequest(buf)
		if err != nil {
			if err == io.EOF {
				log.Println("End of request stream.")
			} else {
				log.Println("Failed to read HTTP request:", err)
				s.Reset()
			}
			return
		}
		defer req.Body.Close()

		// Modify the request as needed
		req.URL.Scheme = "http"
		hp := strings.Split(req.Host, ":")
		if len(hp) > 1 && hp[1] == "443" {
			req.URL.Scheme = "https"
		} else {
			req.URL.Scheme = "http"
		}
		req.URL.Host = req.Host

		outreq := new(http.Request)
		*outreq = *req

		// Make the request
		fmt.Printf("Making request to %s\n", req.URL)
		resp, err := http.DefaultTransport.RoundTrip(outreq)
		if err != nil {
			log.Println("Failed to make request:", err)
			s.Reset()
			return
		}
		defer resp.Body.Close()

		// Write the response back to the stream
		err = resp.Write(s)
		if err != nil {
			log.Println("Failed to write response to stream:", err)
			s.Reset()
			return
		}
		log.Println("Response successfully written to stream.")
	})
}
