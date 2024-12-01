package proxyService

import (
	"fmt"

	"github.com/gorilla/mux"
)

func InitProxyRoutes() *mux.Router {
	r := mux.NewRouter()
	fmt.Print("Inside handleProxyData")
	r.HandleFunc("/proxy-data/", handleProxyData).Methods("GET")
	r.HandleFunc("/proxy-data/", handleProxyData).Methods("POST")
	return r
}
