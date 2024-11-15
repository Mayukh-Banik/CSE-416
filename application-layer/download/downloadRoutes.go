package download

import (
	"github.com/gorilla/mux"
)

func InitDownloadRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/download/request", handleDownloadRequest).Methods("POST")
	r.HandleFunc("/download/respond", handleDownloadRequestOrResponse).Methods("POST")
	r.HandleFunc("/download/getRequests", handleGetPendingRequests).Methods("GET")
	return r
}
