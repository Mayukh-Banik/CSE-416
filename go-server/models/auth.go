package models

// ChallengeResponse contains the challenge and the signed response
type ChallengeResponse struct {
    UserID     string `json:"user_id"`
    Challenge  string `json:"challenge"`
    Signature  string `json:"signature"`
}

// ChallengeRequest contains the challenge request by the user
type ChallengeRequest struct {
    UserID string `json:"user_id"`
}
