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
	OriginalUploader  bool   `json:"OriginalUploader"`
	Rating            string  `json:"Rating"` // either "", upvote, or downvote
	HasVoted          bool   `json:"HasVoted"`
}

type DHTMetadata struct {
	Name              string
	NameWithExtension string
	Type              string
	Size              int64
	Description       string
	CreatedAt         string
	Rating            int64               //
	Providers         map[string]Provider // use PeerID as key
	NumRaters         int64
	Upvote            int64
	Downvote          int64
	Hash              string
}

type Provider struct {
	PeerAddr string
	IsActive bool
	Fee      int64
	Rating   string // upvote, downvote, no vote
}

type Transaction struct {
	Type          string `json:"Type"`        // "request" or "response"
	FileHash      string `json:"FileHash"`    // Unique identifier for the file
	RequesterID   string `json:"RequesterID"` // ID of the requesting node
	RequesterAddr string `json:"RequesterAddr"`
	TargetID      string `json:"TargetID"` // ID of the target node
	TargetAddr    string `json:"TargetAddr"`
	Status        string `json:"Status"`  // "pending", "accepted", "declined"
	Message       string `json:"Message"` // Additional info
	CreatedAt     string `json:"CreatedAt"`
	FileName      string `json:"FileName"`
	TransactionID string `json:"TransactionID"`
}

type RefreshRequest struct {
	Message     string `json:"message"`
	RequesterID string `json:"requesterID"`
	TargetID    string `json:"targetID"`
}
