package routes

import (
    "github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
    router := mux.NewRouter()

    // router.HandleFunc("/signup", controllers.Signup).Methods("POST")

    // router.HandleFunc("/signup", controllers.Signup).Methods("POST")

    AuthRoutes(router)
    UserRoutes(router)

    return router
}
