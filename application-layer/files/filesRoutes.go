package files

import (
	"github.com/gorilla/mux"
)

func InitFileRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/files/upload", uploadFileHandler).Methods("POST")
	r.HandleFunc("/files/fetch", getFiles).Methods("GET")
	r.HandleFunc("/files/getFile", handleGetFileByHash).Methods("GET")
	r.HandleFunc("/files/delete", deleteFile).Methods("DELETE")
	r.HandleFunc("/files/refresh", getMarketplaceFiles).Methods("GET")
	r.HandleFunc("/files/getTransactions", getTransactions).Methods("GET")
	r.HandleFunc("/files/vote", handleVote).Methods("POST")
	r.HandleFunc("/files/getRating", handleGetRating).Methods("GET")
	// r.HandleFunc("/files/searchByName", handleGetFilesByName).Methods("GET")
	return r
}
