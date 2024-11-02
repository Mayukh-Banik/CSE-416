package controllers

import (
	"encoding/json"
	"fmt"
	"go-server/services"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWT generation function
func generateJWT(publicKey string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT_SECRET 환경 변수가 설정되지 않았습니다.")
	}

	claims := jwt.MapClaims{
		"publicKey": publicKey,
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // Expiration time set to 24 hours later
	}
	log.Printf("[generateJWT]: %s", jwtSecret) // 개발 중에만 사용. 실제 서비스에서는 제거하세요.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

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
	log.Printf("RequestChallenge@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	// Generate a new challenge
	challenge, err := services.GenerateChallenge()
	if err != nil {
		http.Error(w, "Failed to generate challenge", http.StatusInternalServerError)
		return
	}

	log.Printf("[challenge]: %s", challenge)

	// The challenge can be stored in a session or cookie. Here, it's simply stored in memory.
	services.StoreChallenge(challenge)

	// Send the challenge to the client
	json.NewEncoder(w).Encode(map[string]string{"challenge": challenge})
}

// VerifyChallenge handles the login process by verifying the signature
func VerifyChallenge(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PublicKey string `json:"public_key"`
		Signature string `json:"signature"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Printf("[Received PublicKey]: \n %s", req.PublicKey)
	log.Printf("[Received Signature]: \n %s", req.Signature)

	// Signature verification
	verified, err := services.VerifySignature(req.PublicKey, req.Signature)
	if err != nil || !verified {
		log.Printf("Error verifying signature: %v", err)
		http.Error(w, "Invalid signature. Login failed.", http.StatusUnauthorized)
		return
	}

	// On successful signature verification, generate a JWT token
	token, err := generateJWT(req.PublicKey)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Set the token as an HttpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	})

	// Success response
	json.NewEncoder(w).Encode(map[string]string{"status": "login successful"})
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

// AuthStatusResponse is the structure for the status check response
type AuthStatusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func CheckAuthStatus(w http.ResponseWriter, r *http.Request) {
	log.Printf("CheckAuthStatus 호출됨")

	// 클라이언트 쿠키에서 토큰 추출
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		log.Println("토큰이 없습니다.")
		return
	}

	// 토큰 문자열 가져오기
	tokenStr := cookie.Value
	log.Printf("받은 토큰 값: %s", tokenStr) // 토큰 값 출력

	// 환경 변수에서 JWT_SECRET 가져오기
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		http.Error(w, "Server configuration error", http.StatusInternalServerError)
		log.Println("JWT_SECRET 환경 변수가 설정되지 않았습니다.")
		return
	}
	log.Printf("JWT Secret Key Used for Verification쮸발: %s", jwtSecret) // 개발 중에만 사용. 실제 서비스에서는 제거하세요.

	// 토큰 검증
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// 서명 방법 확인
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	// 검증 실패 시 Unauthorized 응답
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// 디코딩된 토큰의 클레임 정보 출력
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		log.Printf("Decoded Claims: %+v", claims)
	} else {
		log.Println("Failed to decode claims or token invalid")
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	// 토큰이 유효한 경우 JSON 응답으로 유효성 반환
	response := AuthStatusResponse{
		Status:  "valid",
		Message: "Token is valid",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	log.Println("토큰이 유효합니다.")
}
