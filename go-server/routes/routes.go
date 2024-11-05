package routes

import (
    "github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
    router := mux.NewRouter()

    AuthRoutes(router)
    UserRoutes(router)
    //FileRoutes(router)
    
    return router
}
