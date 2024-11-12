package market

import "application-layer/files"

type MarketFileRequest struct{}

type MarketFileResponse struct {
	MarketFiles []files.FileMetadata
}
