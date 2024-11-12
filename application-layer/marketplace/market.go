package market

import (
	dht_kad "application-layer/dht"
	"application-layer/files"
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

func setupMetadataHandler(host libp2p.Host, cache []files.FileMetadata) {
	host.SetStreamHandler("/market/getAllFiles", func(s network.Stream) {
		defer s.Close()

		// Send cached metadata as response
		response := MarketFileResponse{MarketFiles: cache}
		json.NewEncoder(s).Encode(response)
	})
}

func broadcastMetadata(host libp2p.Host, peers []peer.AddrInfo, metadata files.DHTMetadata) {
	for _, peer := range peers {
		stream, err := host.NewStream(dht_kad.GlobalCtx, peer.ID, "/market/getAllFiles")
		if err != nil {
			fmt.Println("Error opening stream to peer:", err)
			continue
		}
		defer stream.Close()

		// Send metadata to peer
		json.NewEncoder(stream).Encode(metadata)
	}
}

func gatherMetadataFromPeers(host libp2p.Host, peers []peer.AddrInfo) ([]files.DHTMetadata, error) {
	var allMetadata []files.DHTMetadata

	for _, peer := range peers {
		stream, err := host.NewStream(dht_kad.GlobalCtx, peer.ID, "/market/getAllFiles")
		if err != nil {
			fmt.Println("Error opening stream to peer:", err)
			continue
		}
		defer stream.Close()

		var response MarketFileResponse
		if err := json.NewDecoder(stream).Decode(&response); err != nil {
			fmt.Println("Error decoding metadata response:", err)
			continue
		}
		// allMetadata = append(allMetadata, response.MarketFiles...)
	}
	return allMetadata, nil
}
