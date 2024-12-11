package proxyService

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func InitProxyRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/proxy-data/", handleProxyData).Methods("GET")
	r.HandleFunc("/proxy-data/", handleProxyData).Methods("POST")
	r.HandleFunc("/proxy-history/", handleGetProxyHistory).Methods("GET")
	r.HandleFunc("/connect-proxy/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request for /connect-proxy/")
		handleConnectMethod(w, r)
	}).Methods("POST")
	r.HandleFunc("/connect-proxy/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request for /connect-proxy/")
		handleConnectMethod(w, r)
	}).Methods("GET")
	return r
}
