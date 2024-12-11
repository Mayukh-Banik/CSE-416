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
	"net"
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
	node_id          = ""
	peer_id          = ""
	globalCtx        context.Context
	Peer_Addresses   []ma.Multiaddr
	isHost           = true
	connectedPeers   sync.Map
	proxyUpdateMutex sync.Mutex
	proxyHistory     []models.ProxyHistoryEntry
	historyMutex     sync.Mutex
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

			if proxy.PeerID == localPeerID.String() {
				proxy.IsHost = true

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

	}()

	go func() {
		wg.Wait()
		close(done)
	}()

	<-done
	return proxies, nil
}

func pollPeerAddresses(ProxyIsHost bool, ip string) {
	node := dht_kad.Host
	if ProxyIsHost {
		fmt.Println("IN HOST", ip)
		httpHostToClient(node)
	} else {
		fmt.Println("IN CLIENT")
		fmt.Println("IP", ip)
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

	// go dht_kad.ConnectToPeer(node, Cloud_node_addr)
	globalCtx = context.Background()
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
		newProxy.IsHost = true
		log.Print("Debug: New proxy  info", newProxy)

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
		ip, _ := getPrivateIP()
		fmt.Println("BEFORE POLLING", ip)
		go pollPeerAddresses(true, ip)

		if err := json.NewEncoder(w).Encode(responseData); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding proxy data: %v", err), http.StatusInternalServerError)
			return
		}

	}

}

func updateProxyConnections(hostPeerID string, clientPeerID string) {
	proxyUpdateMutex.Lock()
	defer proxyUpdateMutex.Unlock()

	// Retrieve the current proxy information for the host
	proxyInfo, err := getProxyFromDHT(dht_kad.DHT, peer.ID(hostPeerID))
	if err != nil {
		log.Printf("Error retrieving proxy info: %v", err)
		return
	}

	var proxy models.Proxy
	err = json.Unmarshal([]byte(proxyInfo), &proxy)
	if err != nil {
		log.Printf("Error unmarshalling proxy info: %v", err)
		return
	}

	// Add the new client to the connected peers list if not already present
	if !contains(proxy.ConnectedPeers, clientPeerID) {
		proxy.ConnectedPeers = append(proxy.ConnectedPeers, clientPeerID)
	}

	// Save the updated proxy information back to the DHT
	updatedProxyJSON, err := json.Marshal(proxy)
	if err != nil {
		log.Printf("Error marshalling updated proxy info: %v", err)
		return
	}

	err = dht_kad.DHT.PutValue(context.Background(), "/orcanet/proxy/"+hostPeerID, updatedProxyJSON)
	if err != nil {
		log.Printf("Error saving updated proxy info to DHT: %v", err)
		return
	}

	log.Printf("Updated host proxy info with new connected peer: %s", clientPeerID)
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func updateHostProxyInfo(clientPeerID string) {
	proxyInfo, err := getProxyFromDHT(dht_kad.DHT, dht_kad.Host.ID())
	if err != nil {
		log.Printf("Error retrieving proxy info: %v", err)
		return
	}

	var proxy models.Proxy
	err = json.Unmarshal([]byte(proxyInfo), &proxy)
	if err != nil {
		log.Printf("Error unmarshalling proxy info: %v", err)
		return
	}

	// Add the new client to the connected peers list
	proxy.ConnectedPeers = append(proxy.ConnectedPeers, clientPeerID)

	// Save the updated proxy information back to the DHT
	updatedProxyJSON, err := json.Marshal(proxy)
	if err != nil {
		log.Printf("Error marshalling updated proxy info: %v", err)
		return
	}

	err = dht_kad.DHT.PutValue(context.Background(), "/orcanet/proxy/"+proxy.PeerID, updatedProxyJSON)
	if err != nil {
		log.Printf("Error saving updated proxy info to DHT: %v", err)
		return
	}

	log.Printf("Updated host proxy info with new connected peer: %s", clientPeerID)
}

// Function to add a new history entry
func addProxyHistoryEntry(hostPeerID, proxyIP string) {
	historyMutex.Lock()
	defer historyMutex.Unlock()

	newEntry := models.ProxyHistoryEntry{
		HostPeerID: hostPeerID,
		ProxyIP:    proxyIP,
		Timestamp:  time.Now(),
	}

	proxyHistory = append(proxyHistory, newEntry)
}

// Function to send the history to the host
func sendHistoryToHost(hostAddress string) error {
	historyMutex.Lock()
	defer historyMutex.Unlock()

	historyData, err := json.Marshal(proxyHistory)
	if err != nil {
		return fmt.Errorf("failed to marshal history data: %v", err)
	}

	req, err := http.NewRequest("POST", hostAddress+"/update-history", bytes.NewReader(historyData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Sending the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send history to host: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update history on host, status: %d", resp.StatusCode)
	}

	return nil
}
func handleUpdateHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var history []models.ProxyHistoryEntry
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&history)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode history data: %v", err), http.StatusBadRequest)
		return
	}

	// Process the history data (e.g., store in a database, etc.)
	fmt.Printf("Received proxy history: %v\n", history)

	w.WriteHeader(http.StatusOK)
}

func handleConnectMethod(w http.ResponseWriter, r *http.Request) {
	// Log the incoming request method and URL
	fmt.Print("INSIDE THE CONNECT METHOD")
	host_peerid := r.URL.Query().Get("val")
	proxyIP := r.URL.Query().Get("ip")
	fmt.Print("HOST PEER IP", proxyIP)
	fmt.Print(dht_kad.Host.ID().String())
	if host_peerid == dht_kad.Host.ID().String() {
		log.Println("The peer ID matches the current node ID.")
		http.Error(w, "Cannot connect to self.", http.StatusBadRequest)
		return
	}
	if r.Method == "GET" {
		log.Println("Relaying data between client and peer...")
		pollPeerAddresses(false, proxyIP)

		addProxyHistoryEntry(host_peerid, proxyIP)
		err := sendHistoryToHost("http://host-address") // Use actual host address here
		if err != nil {
			log.Printf("Error sending history to host: %v", err)
			http.Error(w, "Failed to send history to host.", http.StatusInternalServerError)
			return
		}

		log.Println("Successfully connected to the peer.")

		w.WriteHeader(http.StatusOK)

	} else {
		log.Printf("Unsupported request method: %s", r.Method)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
func handleGetProxyHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	proxyInfo, err := getProxyFromDHT(dht_kad.DHT, dht_kad.Host.ID())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving proxy info: %v", err), http.StatusInternalServerError)
		return
	}

	var proxy models.Proxy
	err = json.Unmarshal([]byte(proxyInfo), &proxy)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error unmarshalling proxy info: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(proxy.ConnectedPeers)
}

func getPrivateIP() (string, error) {
	// Get all network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("error retrieving network interfaces: %v", err)
	}

	// Iterate over interfaces to find a non-loopback IP address
	for _, iface := range interfaces {
		addresses, err := iface.Addrs()
		if err != nil {
			return "", fmt.Errorf("error getting addresses for interface %v: %v", iface.Name, err)
		}

		for _, addr := range addresses {
			// Ignore loopback IPs
			ip, ok := addr.(*net.IPNet)
			if !ok || ip.IP.IsLoopback() {
				continue
			}

			// IPv4 check and return the first non-loopback IP found
			if ip := ip.IP.To4(); ip != nil {
				return ip.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no private IP found")
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
			existingProxy.Address, _ = getPrivateIP()
			fmt.Println("PRXOYS PRIVATE IP:", existingProxy.Address)
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
		proxy.IsHost = isHost
		proxy.Address, _ = getPrivateIP()
		fmt.Println("PRXOYS PRIVATE IP:", proxy.Address)
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
