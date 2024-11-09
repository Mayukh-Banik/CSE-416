package files

// local file - on user's machine not the dht
type FileMetadata struct {
	Name        string `json:"Name"`
	Type        string `json:"Type"`
	Size        int64  `json:"Size"`
	Description string `json:"Description"`
	Hash        string `json:"Hash"`
	IsPublished bool   `json:"IsPublished"`
	Fee         int64  `json:"Fee"`
}

type DHTMetadata struct {
	Name        string
	Type        string
	Size        int64
	Description string
	Providers   []Provider
}

type Provider struct {
	PeerID   string
	IsActive bool
	Fee      int64
}
