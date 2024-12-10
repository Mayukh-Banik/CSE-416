package proxyService

import (
	dht_kad "application-layer/dht"
	"application-layer/models"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

var (
	node_id        = ""
	peer_id        = ""
	globalCtx      context.Context
	Peer_Addresses []ma.Multiaddr
	isHost         = true
	fileMutex      sync.Mutex
	cancel         context.CancelFunc
	stopFlag       bool
)

const (
	bootstrapNode = "/ip4/35.222.31.85/tcp/61000/p2p/12D3KooWAZv5dC3xtzos2KiJm2wDqiLGJ5y4gwC7WSKU5DvmCLEL"

	// bootstrapNode = "/ip4/130.245.173.221/tcp/6001/p2p/12D3KooWE1xpVccUXZJWZLVWPxXzUJQ7kMqN8UQ2WLn9uQVytmdA"
	// bootstrapNode   = "/ip4/130.245.173.222/tcp/61020/p2p/12D3KooWM8uovScE5NPihSCKhXe8sbgdJAi88i2aXT2MmwjGWoSX"
	proxyKeyPrefix  = "/orcanet/proxy/"
	Cloud_node_addr = "/ip4/35.222.31.85/tcp/61000/p2p/12D3KooWAZv5dC3xtzos2KiJm2wDqiLGJ5y4gwC7WSKU5DvmCLEL"
	Cloud_node_id   = "12D3KooWAZv5dC3xtzos2KiJm2wDqiLGJ5y4gwC7WSKU5DvmCLEL"
)

type ProxyService struct {
	dht  *dht.IpfsDHT
	host host.Host
}

func NewProxyService(dht *dht.IpfsDHT, host host.Host) *ProxyService {
	return &ProxyService{
		dht:  dht,
		host: host,
	}
}

func getProxyFromDHT(dht *dht.IpfsDHT, peerID peer.ID) (string, error) {
	ctx := context.Background()
	key := []byte("/orcanet/proxy/" + peerID.String())
	value, err := dht.GetValue(ctx, string(key))
	if err != nil {
		return "", fmt.Errorf("failed to retrieve proxy info from DHT: %v", err)
	}
	return string(value), nil
}

func getKnownProxyKeys() []string {
	var keys []string
	prefix := "/orcanet/proxy/"

	// Get the known peers from the DHT
	peers := dht_kad.DHT.Host().Peerstore().Peers()

	// Add the current node (itself) to the list of peers
	currentNodeID := dht_kad.DHT.Host().ID()
	peers = append(peers, currentNodeID)

	// Iterate through all peers, including the current node
	for _, peerID := range peers {
		key := prefix + peerID.String()

		// Check if the key exists in the DHT
		value, err := dht_kad.DHT.GetValue(context.Background(), key)
		if err == nil {
			keys = append(keys, key)
			// Optionally, log the value associated with the key
			fmt.Println("Found proxy for key:", key, "with value:", string(value))
		}
	}

	return keys
}

func isEmptyProxy(p models.Proxy) bool {
	return p.Name == "" && p.Location == "" && p.PeerID == "" && p.Address == ""
}

func getAllProxiesFromDHT(dht *dht.IpfsDHT, localPeerID peer.ID, localProxy models.Proxy) ([]models.Proxy, error) {
	ctx := context.Background()
	var proxies []models.Proxy
	seenProxies := make(map[string]struct{}) // Track seen PeerIDs to avoid duplicates
	done := make(chan struct{})

	proxyKeys := getKnownProxyKeys()

	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(len(proxyKeys))

	for _, key := range proxyKeys {
		go func(k string) {
			defer wg.Done()
			log.Printf("Debug: Retrieving proxy info for key: %s", k)
			value, err := dht.GetValue(ctx, k)
			if err != nil {
				log.Printf("Debug: Error retrieving proxy info for key %s: %v", k, err)
				return
			}

			var proxy models.Proxy
			err = json.Unmarshal(value, &proxy)
			if err != nil {
				log.Printf("Debug: Error unmarshalling proxy data for key %s: %v", k, err)
				return
			}

			// Avoid duplicates by checking the PeerID
			if _, seen := seenProxies[proxy.PeerID]; !seen {
				mu.Lock()
				proxies = append(proxies, proxy)
				seenProxies[proxy.PeerID] = struct{}{} // Mark this PeerID as seen
				mu.Unlock()
			}
		}(key)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if !isEmptyProxy(localProxy) {
			mu.Lock()
			fmt.Println("Local proxy", localProxy)
			proxies = append(proxies, localProxy)
			mu.Unlock()
			log.Printf("Debug: Added local proxy information")
		}
	}()

	go func() {
		wg.Wait()
		close(done)
	}()

	<-done
	return proxies, nil
}

func pollPeerAddresses(node host.Host) {
	if isHost {
		httpHostToClient(node)
	} else {
		fmt.Println("RIGHT BEFORE ACCESSING PEER ADDRESSES1")
		for len(Peer_Addresses) == 0 {
			time.Sleep(3 * time.Second)
		}
		var ip string
		var err error
		fmt.Println("RIGHT BEFORE ACCESSING PEER ADDRESSES2")
		for _, val := range Peer_Addresses {
			ip, err = val.ValueForProtocol(ma.P_IP4)
			if err != nil {
				continue
			}
			break
		}
		var script string
		var args []string
		script = "proxy/client.py"
		args = []string{"--remote-host", ip}

		// Function to run the command
		runCommand := func(pythonCmd string) error {
			cmd := exec.Command(pythonCmd, append([]string{script}, args...)...)
			cmd.Stdout = os.Stderr // Redirect standard output to stderr
			cmd.Stderr = os.Stderr // Redirect standard error to stderr
			return cmd.Run()
		}

		// Try running with `python`
		if err := runCommand("python"); err != nil {
			fmt.Println("`python` not found or failed, trying `python3`...")
			// If `python` fails, try `python3`
			if err := runCommand("python3"); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to run %s with both `python` and `python3`: %v\n", script, err)
			}
		}
	}
}

func getAdjacentNodeProxiesMetadata(w http.ResponseWriter, r *http.Request) {
	// for _, node := range dht_kad.RoutingTable.NearestPeers(kbucket.ID(peer_id), 5) {
	// 	fmt.Println("node: ", node)
	// }

	// Retrieve connected peers
	adjacentNodes := dht_kad.Host.Network().Peers()
	fmt.Println("Connected peers:", adjacentNodes)

	var sendWG sync.WaitGroup
	var responseWG sync.WaitGroup

	// Iterate over peers and request proxy metadata
	for _, peer := range adjacentNodes {
		peerID := peer.String()
		if peerID != dht_kad.Bootstrap_node_addr && peerID != dht_kad.PeerID && nodeSupportRefreshStreams(peer) {
			sendWG.Add(1)
			responseWG.Add(1)
			go func(peerID string) {
				defer responseWG.Done()
				go dht_kad.SendProxyRequest(peerID, &sendWG) // Adjust to match your request handler for proxies
			}(peerID)
		}
	}

	// Wait for all requests to complete
	sendWG.Wait()
	responseWG.Wait()

	// Introduce a short delay if necessary for processing
	<-time.After(3 * time.Second)

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode the collected proxy metadata as JSON
	if err := json.NewEncoder(w).Encode(dht_kad.ProxyResponse); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func nodeSupportRefreshStreams(peerID peer.ID) bool {
	supportSendRefreshRequest := false
	supportSendRefreshResponse := false

	protocols, _ := dht_kad.Host.Peerstore().GetProtocols(peerID)
	fmt.Printf("protocols supported by peer %v: %v\n", peerID, protocols)

	for _, protocol := range protocols {
		if protocol == "/sendRefreshRequest/p2p" {
			supportSendRefreshRequest = true
		} else if protocol == "/sendRefreshResponse/p2p" {
			supportSendRefreshResponse = true
		}
	}
	return supportSendRefreshRequest && supportSendRefreshResponse
}

func handlePeerExchange(node host.Host) {
	bootstrap_node_info, _ := peer.AddrInfoFromString(dht_kad.Bootstrap_node_addr)
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
						if string(peerID) != string(bootstrap_node_info.ID) {
							dht_kad.ConnectToPeer(node, peerID)
						}
					}
				}
			}
		}
	})
}

// Retrieveing proxies data, and adding yourself as host
func handleProxyData(w http.ResponseWriter, r *http.Request) {
	node := dht_kad.Host
	// go dht_kad.ConnectToPeer(node, dht_kad.Bootstrap_node_addr)
	// go dht_kad.ConnectToPeer(node, Cloud_node_addr)
	globalCtx, cancel = context.WithCancel(context.Background())
	stopFlag = false
	if r.Method == "POST" {
		isHost = true
		var newProxy models.Proxy
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&newProxy)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode proxy data: %v", err), http.StatusBadRequest)
			return
		}

		newProxy.Address = node.Addrs()[0].String()
		newProxy.PeerID = node.ID().String()
		log.Printf("Debug: New proxy created with PeerID: %s", newProxy.PeerID)

		if err := saveProxyToDHT(newProxy); err != nil {
			log.Printf("Debug: Failed to save proxy to DHT: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("Debug: Proxy saved to DHT successfully")

		proxyInfo, err := getAllProxiesFromDHT(dht_kad.DHT, node.ID(), newProxy)
		if err != nil {
			log.Printf("Debug: Error retrieving proxies from DHT: %v", err)
		} else {
			log.Printf("Debug: Retrieved %d proxies from DHT", len(proxyInfo))
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(proxyInfo); err != nil {
			log.Printf("Debug: Error encoding proxy data: %v", err)
			http.Error(w, fmt.Sprintf("Error encoding proxy data: %v", err), http.StatusInternalServerError)
		}
		return
	}

	log.Printf("Debug: Connecting to bootstrap node")
	getAdjacentNodeProxiesMetadata(w, r)

	go pollPeerAddresses(node)

	if r.Method == "GET" {
		// clearAllProxies()

		proxyInfo, err := getAllProxiesFromDHT(dht_kad.DHT, node.ID(), models.Proxy{})
		if err != nil {
			http.Error(w, fmt.Sprintf("Error retrieving proxies: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		// Ensure proxyInfo is wrapped in an array if it's not already
		var responseData []models.Proxy
		if len(proxyInfo) == 0 {
			responseData = []models.Proxy{}
		} else {
			responseData = proxyInfo
		}

		if err := json.NewEncoder(w).Encode(responseData); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding proxy data: %v", err), http.StatusInternalServerError)
			return
		}

	}

}

// handleConnectMethod handles the CONNECT HTTP method for tunneling.
func handleConnectMethod(w http.ResponseWriter, r *http.Request) {
	// Log the incoming request method and URL
	fmt.Print("INSIDE THE CONNECT METHOD")
	host_peerid := r.URL.Query().Get("val")
	fmt.Print("HOST PEER ID", host_peerid)
	// Check if the request method is POST
	if r.Method == "GET" {
		log.Println("Processing POST request")

		// Parse the destination from the CONNECT request
		node := dht_kad.Host
		destination := r.URL.Host

		// Check if the destination is empty and log the error if so
		if destination == "" {
			log.Println("Destination not specified")
			http.Error(w, "Destination not specified", http.StatusBadRequest)
			return
		}
		log.Printf("Destination: %s", destination)

		// Ensure the destination is reachable via the DHT or peer network
		peerAddr := "/p2p-circuit/p2p/" + destination
		log.Printf("Constructed peer address: %s", peerAddr)

		// Try to create a Multiaddr for the peer
		maddr, err := ma.NewMultiaddr(peerAddr)
		if err != nil {
			log.Printf("Invalid peer address: %v", err)
			http.Error(w, fmt.Sprintf("Invalid peer address: %v", err), http.StatusBadRequest)
			return
		}
		log.Printf("Multiaddr created: %s", maddr)

		// Attempt to get peer information from the Multiaddr
		peerInfo, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			log.Printf("Failed to get peer info: %v", err)
			http.Error(w, fmt.Sprintf("Failed to get peer info: %v", err), http.StatusInternalServerError)
			return
		}
		log.Printf("Peer info: %+v", peerInfo)

		// Attempt to connect to the destination peer
		log.Println("Attempting to connect to the peer...")
		err = node.Connect(globalCtx, *peerInfo)
		if err != nil {
			log.Printf("Failed to connect to peer: %v", err)
			http.Error(w, fmt.Sprintf("Failed to connect to peer: %v", err), http.StatusInternalServerError)
			return
		}
		log.Println("Successfully connected to the peer.")

		// Respond with 200 OK to indicate the connection has been established
		w.WriteHeader(http.StatusOK)

		// Now relay data between client and peer asynchronously
		log.Println("Relaying data between client and peer...")
		go tunnelDataBetweenClientAndPeer(node, peerInfo.ID, r.Body, w)
	} else {
		// If the method is not POST, log it and return method not allowed
		log.Printf("Unsupported request method: %s", r.Method)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// tunnelDataBetweenClientAndPeer relays data between the client and the destination peer.
func tunnelDataBetweenClientAndPeer(node host.Host, peerID peer.ID, clientReader io.Reader, w http.ResponseWriter) {
	// Open a stream to the connected peer (destination)
	stream, err := node.NewStream(globalCtx, peerID, "/http-temp-protocol")
	if err != nil {
		log.Printf("Failed to open stream to peer %s: %v", peerID, err)
		return
	}
	defer stream.Close()

	// Relay data from client to peer
	go func() {
		_, err := io.Copy(stream, clientReader) // Read data from client (r.Body) and write to the peer stream
		if err != nil {
			log.Printf("Failed to forward data from client to peer: %v", err)
		}
	}()

	// Relay data from peer to client
	_, err = io.Copy(w, stream) // Read data from peer stream and write to the client (http.ResponseWriter)
	if err != nil {
		log.Printf("Failed to forward data from peer to client: %v", err)
	}
}

func saveProxyToDHT(proxy models.Proxy) error {
	ctx := context.Background()
	key := "/orcanet/proxy/" + proxy.PeerID

	// Check if the proxy already exists
	existingValue, err := dht_kad.DHT.GetValue(ctx, key)
	if err == nil {
		// Proxy exists, update it
		var existingProxy models.Proxy
		if err := json.Unmarshal(existingValue, &existingProxy); err != nil {
			return fmt.Errorf("failed to unmarshal existing proxy data: %v", err)
		}

		// Check if the new proxy's PeerID matches the existing one
		if existingProxy.PeerID == proxy.PeerID {
			// If they are the same, either update or reject
			existingProxy.Name = proxy.Name
			existingProxy.Location = proxy.Location
			existingProxy.Price = proxy.Price
			existingProxy.Statistics = proxy.Statistics
			existingProxy.Bandwidth = proxy.Bandwidth
			existingProxy.IsEnabled = proxy.IsEnabled

			// Serialize and update the proxy as needed
			updatedProxyJSON, err := json.Marshal(existingProxy)
			if err != nil {
				return fmt.Errorf("failed to serialize updated proxy data: %v", err)
			}

			err = dht_kad.DHT.PutValue(ctx, key, updatedProxyJSON)
			if err != nil {
				return fmt.Errorf("failed to update proxy in DHT: %v", err)
			}

			fmt.Printf("Proxy updated successfully in DHT for PeerID: %s\n", proxy.PeerID)
		}
	} else {
		// Proxy doesn't exist, add it as a new entry
		proxyJSON, err := json.Marshal(proxy)

		if err != nil {
			return fmt.Errorf("failed to serialize new proxy data: %v", err)
		}

		err = dht_kad.DHT.PutValue(ctx, key, proxyJSON)
		if err != nil {
			return fmt.Errorf("failed to store new proxy in DHT: %v", err)
		}

		fmt.Printf("New proxy added successfully to DHT for PeerID: %s\n", proxy.PeerID)
	}
	return nil
}

func clientHTTPRequestToHost(node host.Host, targetpeerid string, req *http.Request, w http.ResponseWriter) {
	fmt.Println("In client http request to host")
	var ctx = context.Background()
	targetPeerID := strings.TrimSpace(targetpeerid)
	bootstrapAddr, err := ma.NewMultiaddr(dht_kad.Bootstrap_node_addr)
	if err != nil {
		log.Printf("Failed to create bootstrapAddr multiaddr: %v", err)
		return
	}
	peerMultiaddr := bootstrapAddr.Encapsulate(ma.StringCast("/p2p-circuit/p2p/" + targetPeerID))

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
	var script string
	var args []string
	script = "proxy/server.py"
	args = []string{}

	// Function to run the command
	runCommand := func(pythonCmd string) error {
		cmd := exec.Command(pythonCmd, append([]string{script}, args...)...)
		cmd.Stdout = os.Stderr // Redirect standard output to stderr
		cmd.Stderr = os.Stderr // Redirect standard error to stderr
		return cmd.Run()
	}

	// Try running with `python`
	if err := runCommand("python"); err != nil {
		fmt.Println("`python` not found or failed, trying `python3`...")
		// If `python` fails, try `python3`
		if err := runCommand("python3"); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to run %s with both `python` and `python3`: %v\n", script, err)
		}
	}
}

func clearAllProxies() {
	ctx := context.Background()

	// Get all known proxy keys
	proxyKeys := getKnownProxyKeys()

	for _, key := range proxyKeys {
		emptyProxy := models.Proxy{}
		emptyProxyJSON, err := json.Marshal(emptyProxy)
		if err != nil {
			log.Printf("Failed to marshal empty proxy: %v", err)
			continue
		}

		err = dht_kad.DHT.PutValue(ctx, key, emptyProxyJSON)
		if err != nil {
			log.Printf("Failed to clear proxy for key %s: %v", key, err)
		} else {
			log.Printf("Proxy for key %s cleared", key)
		}
	}
}
