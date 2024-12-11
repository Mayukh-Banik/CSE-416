package models

// local file - on user's machine not the dht
type FileMetadata struct {
	Name              string `json:"Name"`
	NameWithExtension string `json:"NameWithExtension"`
	Type              string `json:"Type"`
	Size              int64  `json:"Size"`
	Description       string `json:"Description"`
	Hash              string `json:"Hash"`
	IsPublished       bool   `json:"IsPublished"`
	Fee               int64  `json:"Fee"`
	CreatedAt         string `json:"CreatedAt"`
	Reputation        int64  `json:"Reputation"`
	OriginalUploader  bool   `json:"OriginalUploader"`
}

type DHTMetadata struct {
	Name              string
	NameWithExtension string
	Type              string
	Size              int64
	Description       string
	CreatedAt         string
	Reputation        int64
	Providers         []Provider
}

type Provider struct {
	PeerID   string
	PeerAddr string
	IsActive bool
	Fee      int64
}

type Transaction struct {
	Type          string `json:"Type"`        // "request" or "response"
	FileHash      string `json:"FileHash"`    // Unique identifier for the file
	RequesterID   string `json:"RequesterID"` // ID of the requesting node
	TargetID      string `json:"TargetID"`    // ID of the target node
	Status        string `json:"Status"`      // "pending", "accepted", "declined"
	Message       string `json:"Message"`     // Additional info
	CreatedAt     string `json:"CreatedAt"`
	FileName      string `json:"FileName"`
	TransactionID string `json:"TransactionID"`
}

type RefreshRequest struct {
	Message     string `json:"message"`
	RequesterID string `json:"requesterID"`
	TargetID    string `json:"targetID"`
}
