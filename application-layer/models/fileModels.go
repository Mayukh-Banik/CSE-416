package models

// local file - on user's machine not the dht
type FileMetadata struct {
	Name        string `json:"Name"`
	Type        string `json:"Type"`
	Size        int64  `json:"Size"`
	Description string `json:"Description"`
	Hash        string `json:"Hash"`
	IsPublished bool   `json:"IsPublished"`
	Fee         int64  `json:"Fee"`
	CreatedAt   string `json:"CreatedAt"`
	Reputation  int64  `json:"Reputation"`
}

type DHTMetadata struct {
	Name        string
	Type        string
	Size        int64
	Description string
	CreatedAt   string
	Reputation  int64
	Providers   []Provider
}

type Provider struct {
	PeerID   string
	PeerAddr string
	IsActive bool
	Fee      int64
}

type Transaction struct {
	Type        string `json:"type"`        // "request" or "response"
	FileHash    string `json:"fileHash"`    // Unique identifier for the file
	RequesterID string `json:"requesterID"` // ID of the requesting node
	TargetID    string `json:"targetID"`    // ID of the target node
	Status      string `json:"status"`      // "pending", "accepted", "declined"
	Message     string `json:"message"`     // Additional info
	CreatedAt   string `json:"CreatedAt"`
	FileName    string `json:"fileName"`
}
