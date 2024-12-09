package proxyService

import (
	"github.com/gorilla/mux"
)

func InitProxyRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/proxy-data/", handleProxyData).Methods("GET")
	r.HandleFunc("/proxy-data/", handleProxyData).Methods("POST")

	r.HandleFunc("/connect-proxy/", handleConnectMethod).Methods("POST")
	return r
}
