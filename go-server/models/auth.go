package models

import "time"

// ChallengeResponse contains the challenge and the signed response
type ChallengeResponse struct {
	ID        string `json:"user_id"`    // 퍼블릭 키 또는 그 ID
	PublicKey string `json:"public_key"` // 클라이언트의 퍼블릭 키
	Signature string `json:"signature"`  // 클라이언트가 생성한 서명
	CreatedAt time.Time
	ExpiresAt time.Time
}

// ChallengeRequest contains the challenge request by the user
type ChallengeRequest struct {
	PublicKey string `json:"public_key"`
}

type ChallengeData struct {
	Challenge string
	CreatedAt time.Time
}
