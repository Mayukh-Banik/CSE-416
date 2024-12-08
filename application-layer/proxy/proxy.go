package proxyService

import (
	dht_kad "application-layer/dht"
	"application-layer/models"
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/client"
	ma "github.com/multiformats/go-multiaddr"
)

var (
	node_id         = ""
	peer_id         = ""
	relayNode       = "12D3KooWDpJ7As7BWAwRMfu1VU2WCqNjvq387JEYKDBj4kx6nXTN"
	bootstrapNode   = "12D3KooWE1xpVccUXZJWZLVWPxXzUJQ7kMqN8UQ2WLn9uQVytmdA"
	relay_node_addr = "/ip4/130.245.173.221/tcp/4001/p2p/12D3KooWDpJ7As7BWAwRMfu1VU2WCqNjvq387JEYKDBj4kx6nXTN"
	globalCtx       context.Context
	Peer_Addresses  []ma.Multiaddr
	isHost          = true
	fileMutex       sync.Mutex
)

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
	fmt.Println("Known peers:", peers)

	// Add the current node (itself) to the list of peers
	currentNodeID := dht_kad.DHT.Host().ID()
	peers = append(peers, currentNodeID)
	fmt.Println("Including current node:", currentNodeID)

	// Iterate through all peers, including the current node
	for _, peerID := range peers {
		key := prefix + peerID.String()
		fmt.Println("Checking key:", key)

		// Check if the key exists in the DHT
		value, err := dht_kad.DHT.GetValue(context.Background(), key)
		if err == nil {
			keys = append(keys, key)
			// Optionally, log the value associated with the key
			fmt.Println("Found proxy for key:", key, "with value:", string(value))
		} else {
			fmt.Println("Error retrieving key:", key, "Error:", err)
		}
	}

	return keys
}
func isEmptyProxy(p models.Proxy) bool {
	return p.Name == "" && p.Location == "" && p.PeerID == "" && p.Address == ""
}

func getAllProxiesFromDHT(dht *dht.IpfsDHT, localPeerID peer.ID, localProxy models.Proxy) ([]models.Proxy, error) {
	log.Printf("Debug: Starting getAllProxiesFromDHT function")
	ctx := context.Background()
	var proxies []models.Proxy
	done := make(chan struct{})

	proxyKeys := getKnownProxyKeys()
	log.Printf("Debug: Found %d known proxy keys", len(proxyKeys))

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

			mu.Lock()
			proxies = append(proxies, proxy)
			mu.Unlock()
			log.Printf("Debug: Successfully added proxy for key: %s", k)
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
		} else {
			log.Printf("Debug: Skipped empty local proxy")
		}
	}()

	go func() {
		wg.Wait()
		close(done)
	}()

	<-done
	log.Printf("Debug: getAllProxiesFromDHT function completed. Found %d proxies", len(proxies))
	return proxies, nil
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
							connectToPeer(node, peerID)
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
func getAdjacentNodeProxiesMetadata(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Trying to get adjacent node proxies in backend")

	relayNode := "12D3KooWDpJ7As7BWAwRMfu1VU2WCqNjvq387JEYKDBj4kx6nXTN"
	bootstrapNode := "12D3KooWQd1K1k8XA9xVEzSAu7HUCodC7LJB6uW5Kw4VwkRdstPE"

	// Retrieve connected peers
	adjacentNodes := dht_kad.Host.Network().Peers()
	fmt.Println("Connected peers:", adjacentNodes)

	var sendWG sync.WaitGroup
	var responseWG sync.WaitGroup

	// Iterate over peers and request proxy metadata
	for _, peer := range adjacentNodes {
		peerID := peer.String()
		if peerID != relayNode && peerID != bootstrapNode && peerID != dht_kad.PeerID && nodeSupportRefreshStreams(peer) {
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

	fmt.Println("getAdjacentNodeProxiesMetadata: received everyone's proxy metadata: ", dht_kad.ProxyResponse)

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode the collected proxy metadata as JSON
	if err := json.NewEncoder(w).Encode(dht_kad.ProxyResponse); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}

	fmt.Println("getAdjacentNodeProxiesMetadata: response sent to frontend")
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

// Retrieveing proxies data, and adding yourself as host
func handleProxyData(w http.ResponseWriter, r *http.Request) {
	log.Printf("Debug: Handling proxy data request. Method: %s", r.Method)
	node := dht_kad.Host

	globalCtx = context.Background()

	if r.Method == "POST" {
		log.Printf("Debug: Processing POST request")
		isHost = true
		var newProxy models.Proxy
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&newProxy)
		if err != nil {
			log.Printf("Debug: Failed to decode proxy data: %v", err)
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

	log.Printf("Debug: Connecting to relay node")
	connectToPeer(node, relay_node_addr)
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
