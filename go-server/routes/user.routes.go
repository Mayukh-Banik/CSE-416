package routes

import (
    // "go-server/controllers"
    "github.com/gorilla/mux"
)

// UserRoutes는 /users 관련 경로들을 설정합니다.
func UserRoutes() *mux.Router {
    router := mux.NewRouter()

    // // /users/info 경로
    // router.HandleFunc("/api/auth/info", controllers.GetUserInfo).Methods("GET")

    // // /users/{id} 경로
    // router.HandleFunc("/api/auth/{id}", controllers.GetUserByID).Methods("GET")

    return router
}
