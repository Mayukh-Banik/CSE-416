// application-layer/routes/routes.go
package routes

import (
	"application-layer/controllers"

	"github.com/gorilla/mux"
)

// RegisterRoutes는 모든 주요 라우트를 등록합니다.
func RegisterRoutes(router *mux.Router, btcController *controllers.BtcController) {
	RegisterBtcRoutes(router, btcController)  // /api/btc 라우트 등록
	RegisterAuthRoutes(router, btcController) // /api/auth 라우트 등록

	// 다른 라우트도 여기서 등록 가능
}
