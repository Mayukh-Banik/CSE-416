package routes

import (
    "net/http"
    "github.com/gorilla/mux"
)

func UserRoutes(router *mux.Router) {
    router.HandleFunc("/api/user/profile", profileHandler).Methods("GET")
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("User profile endpoint"))
}
