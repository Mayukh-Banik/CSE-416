package controllers

import (
	"encoding/json"
	"go-server/models"
	"go-server/services"
	"io"
	"log"
	"net/http"
)

// Signup handles the signup request
func Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Println("Signup request received")

	// Call the wallet service to generate a key pair
	walletResponse, err := services.GenerateWallet()
	if err != nil {
		http.Error(w, "Error generating wallet", http.StatusInternalServerError)
		return
	}

	// Log success message
	log.Printf("User %s created successfully", walletResponse.PublicKey)

	// Return the wallet response to the frontend
	json.NewEncoder(w).Encode(walletResponse)
}

// Login handles the login process (to be implemented)
func Login(w http.ResponseWriter, r *http.Request) {
	// Logic for login (to be implemented)
}

// RequestChallenge handles the request for a login challenge
func RequestChallenge(w http.ResponseWriter, r *http.Request) {

	var challengeReq models.ChallengeRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &challengeReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Printf("Received PublicKey: %s", challengeReq.PublicKey)

		// Check if the user exists for the given public key
	user, err := services.FindUserByPublicKey(challengeReq.PublicKey)
	if err != nil {
		// User not found, return 401 Unauthorized
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// 퍼블릭 키가 존재하든 존재하지 않든 난수를 생성하여 반환
	challengeData, exists := services.GetChallenge(user.PublicKey)
	if exists {
		log.Printf("Returning existing challenge for publicKey: %s", challengeReq.PublicKey)
		response := map[string]string{"challenge": challengeData.Challenge}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
		return
	}

	// 새로운 난수 생성
	challenge, err := services.GenerateChallenge()
	if err != nil {
		http.Error(w, "Error generating challenge", http.StatusInternalServerError)
		return
	}

	// 생성된 난수 로그 출력
	log.Printf("Generated new challenge: %s", challenge)

	// 생성된 난수를 메모리에 저장 (퍼블릭 키와 함께)
	if err := services.StoreChallenge(challengeReq.PublicKey, challenge); err != nil {
		http.Error(w, "Error storing challenge", http.StatusInternalServerError)
		return
	}

	// 생성된 난수를 클라이언트로 반환
	response := map[string]string{"challenge": challenge}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// VerifyChallenge handles the login process by verifying the signature
func VerifyChallenge(w http.ResponseWriter, r *http.Request) {
	var challengeResponse models.ChallengeResponse

	// 요청 본문을 파싱
	if err := json.NewDecoder(r.Body).Decode(&challengeResponse); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Error decoding request payload: %v", err) // 로그 추가
		return
	}

	// 로그 추가: 수신한 챌린지 응답 출력
	log.Printf("Received Challenge Response: %+v", challengeResponse)

	// 퍼블릭 키로 저장된 챌린지와 서명을 검증
	verified, err := services.VerifySignature(challengeResponse.ID, challengeResponse.Signature)
	if err != nil {
		log.Printf("Error verifying signature: %v", err) // 서명 검증 에러 로그
		http.Error(w, `{"error": "Invalid signature. Login failed."}`, http.StatusUnauthorized)
		return
	}
	if !verified {
		log.Printf("Invalid signature for user ID: %s", challengeResponse.ID) // 유효하지 않은 서명 로그
		http.Error(w, `{"error": "Invalid signature. Login failed."}`, http.StatusUnauthorized)
		return
	}

	// 로그 추가: 서명이 성공적으로 검증되었음을 알림
	log.Printf("Signature verified successfully for user ID: %s", challengeResponse.ID)

	// 데이터베이스에서 퍼블릭 키로 사용자 존재 여부 확인
	user, err := services.FindUserByPublicKey(challengeResponse.ID)
	if err != nil {
		log.Printf("User not found for public key: %s", challengeResponse.ID) // 사용자 미발견 로그
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// 사용자 인증이 성공적으로 이루어졌음
	log.Printf("User authenticated successfully: %s", user.ID.Hex()) // 인증 성공 로그
	json.NewEncoder(w).Encode(map[string]string{"status": "login successful", "userId": user.ID.Hex()})
}

// LoginWithWalletID handles the login using only the walletId (public key)
func LoginWithWalletID(w http.ResponseWriter, r *http.Request) {
	var req struct {
		WalletID string `json:"walletId"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Call the service to login with walletId (public key)
	success, err := services.LoginWithWalletID(req.WalletID)
	if err != nil {
		http.Error(w, "Login failed: "+err.Error(), http.StatusUnauthorized)
		return
	}

	if success {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
	} else {
		http.Error(w, "Login failed", http.StatusUnauthorized)
	}
}

// AuthStatusResponse는 상태 확인 응답 구조체입니다.
type AuthStatusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// // AuthStatus는 JWT 토큰 상태를 검증하는 핸들러입니다.
// func AuthStatus(w http.ResponseWriter, r *http.Request) {
// 	// Authorization 헤더에서 토큰 추출
// 	token := r.Header.Get("Authorization")
// 	if token == "" {
// 		http.Error(w, "Missing token", http.StatusUnauthorized)
// 		return
// 	}

// 	// 토큰 검증
// 	valid, err := services.VerifyToken(token)
// 	if err != nil || !valid {
// 		http.Error(w, "Invalid token", http.StatusUnauthorized)
// 		return
// 	}

// 	// 토큰이 유효한 경우
// 	response := AuthStatusResponse{
// 		Status:  "valid",
// 		Message: "Token is valid",
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }
