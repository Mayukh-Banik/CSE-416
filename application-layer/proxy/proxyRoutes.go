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

	r.HandleFunc("/disconnect-from-proxy/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Recieved request for /disconnect-proxy/")
		handleDisconnectFromProxy(w, r)
	}).Methods("GET")

	r.HandleFunc("/stop-hosting/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Recieved request for /disconnect-proxy/")
		stopHosting(w, r)
	}).Methods("GET")

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
