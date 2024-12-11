package proxyService

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func InitProxyRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/proxy-data/", handleProxyData).Methods("GET")
	r.HandleFunc("/proxy-data/", handleProxyData).Methods("POST")

	r.HandleFunc("/disconnect-from-proxy/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("\n\n\nRecieved request for /disconnect-proxy/")
		handleDisconnectFromProxy(w, r)
	}).Methods("GET")
	r.HandleFunc("/disconnect-from-proxy/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("\n\n\nRecieved request for /disconnect-proxy/")
		handleDisconnectFromProxy(w, r)
	}).Methods("POST")

	r.HandleFunc("/stop-hosting/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("\n\n\nRecieved request for /disconnect-proxy/")
		stopHosting(w, r)
	}).Methods("GET")
	r.HandleFunc("/stop-hosting/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("\n\n\nRecieved request for /disconnect-proxy/")
		stopHosting(w, r)
	}).Methods("POST")

	//Fetchiing history
	r.HandleFunc("/proxy-history/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request for /proxy-history/")
		handleGetProxyHistory(w, r)
	}).Methods("POST", "GET")

	r.HandleFunc("/connect-proxy/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request for /connect-proxy/")
		handleConnectMethod(w, r)
	}).Methods("POST")
	r.HandleFunc("/connect-proxy/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request for /connect-proxy/")
		handleConnectMethod(w, r)
	}).Methods("GET")

	r.HandleFunc("/check-balance/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Check balance")
		handleCheckBalance(w, r)
	}).Methods("POST", "GET")
	return r
}
