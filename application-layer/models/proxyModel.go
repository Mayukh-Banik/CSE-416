package models

import "time"

type Proxy struct {
	Name       string   `json:"name"`
	Location   string   `json:"location"`
	Logs       []string `json:"logs"`
	Statistics struct {
		Uptime string `json:"uptime"`
	}
	Bandwidth         string              `json:"bandwidth"`
	Address           string              `json:"address"`
	PeerID            string              `json:"peer_id"`
	IsEnabled         bool                `json:"isEnabled"`
	IsHost            bool                `json:"isHost"`
	Price             string              `json:"price"`
	ConnectedTimed    time.Time           `json:"connected_time"`
	ConnectedPeers    []string            `json:"connected_peers"` // Add this field
	ConnectionHistory []ProxyHistoryEntry `json:"history"`
}

type ProxyHistoryEntry struct {
	ClientPeerID string    `json:"client_peer_id"`
	Timestamp    time.Time `json:"timestamp"`
}
