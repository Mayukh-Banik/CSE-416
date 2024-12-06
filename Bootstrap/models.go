package main 

type FileMetadata struct {
	Name             string `json:"Name"`
	Type             string `json:"Type"`
	Size             int64  `json:"Size"`
	Description      string `json:"Description"`
	Hash             string `json:"Hash"`
	IsPublished      bool   `json:"IsPublished"`
	Fee              int64  `json:"Fee"`
	CreatedAt        string `json:"CreatedAt"`
	Reputation       int64  `json:"Reputation"`
	OriginalUploader bool   `json:"OriginalUploader"`
	Extension        string `json:"Extension"`
}