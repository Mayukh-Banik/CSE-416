package dht_kad

import (
	"context"
	"fmt"
	"log"
	"time"
)

func StartDHTService() {
	getNodeId()
	node, dht, err := createNode()
	PeerID = node.ID().String()
	if err != nil {
		log.Fatalf("Failed to create node: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	GlobalCtx = ctx

	DHT = setupDHT(ctx, dht.Host())
	ProviderStore = DHT.ProviderStore()
	ConnectToPeer(node, Relay_node_addr) // connect to relay node
	makeReservation(node)                // make reservation on relay node
	go refreshReservation(node, 5*time.Minute)
	ConnectToPeer(node, Bootstrap_node_addr) // connect to bootstrap node
	go handlePeerExchange(node)

	ReceiveDataFromPeer(node) //listen on stream /senddata/p2p
	setupStreams(node)

	fmt.Println("My Node MULTIADDRESS:", node.Addrs())
	fmt.Println("MY NODE PEER ID:", PeerID)
	fmt.Println("Supported protocols:", node.Mux().Protocols())

	go handleInput(ctx, dht)
	Host = node
	// block until a signal is received
	select {}
}
