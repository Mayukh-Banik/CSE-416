package controllers

// import (
// 	"application-layer/wallet"
// 	"encoding/json"
// 	"log"
// 	"net/http"
// )

// // WalletController는 지갑 관련 요청을 처리합니다.
// type WalletController struct {
// 	WalletService wallet.WalletServiceInterface
// }

// // NewWalletController는 WalletController를 초기화합니다.
// func NewWalletController(ws wallet.WalletServiceInterface) *WalletController {
// 	return &WalletController{
// 		WalletService: ws,
// 	}
// }

// // walletController.go
// func (wc *WalletController) HandleGenerateWallet(w http.ResponseWriter, r *http.Request) {
// 	log.Println("HandleGenerateWallet 요청 수신됨") // 로그 추가

// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	var req struct {
// 		Passphrase string `json:"passphrase"`
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		log.Println("요청 본문 디코딩 오류:", err) // 로그 추가
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	log.Println("Passphrase:", req.Passphrase) // 로그 추가

// 	address, pubKey, privateKey, balance, err := wc.WalletService.GenerateNewAddressWithPubKeyAndPrivKey(req.Passphrase)
// 	if err != nil {
// 		log.Printf("지갑 생성 오류: %v", err) // 로그 추가
// 		http.Error(w, "Failed to generate wallet", http.StatusInternalServerError)
// 		return
// 	}

// 	resp := map[string]interface{}{
// 		"address":     address,
// 		"public_key":  pubKey,
// 		"private_key": privateKey,
// 		"balance":     balance,
// 	}

// 	log.Println("지갑 생성 성공:", resp) // 로그 추가

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(resp)
// }
