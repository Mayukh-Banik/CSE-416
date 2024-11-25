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
	r.HandleFunc("/files/refresh", getAdjacentNodeFilesMetadata).Methods("GET")
	r.HandleFunc("/files/getTransactions", getTransactions).Methods("GET")
	return r
}
