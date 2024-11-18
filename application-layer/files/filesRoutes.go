package files

import (
	"github.com/gorilla/mux"
)

func InitFileRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/files/upload", uploadFileHandler).Methods("POST")
	r.HandleFunc("/files/fetchAll", getUploadedFiles).Methods("GET")
	r.HandleFunc("/files/getFile", handleGetFileByHash).Methods("GET")
	r.HandleFunc("/files/delete", deleteFile).Methods("DELETE")
	return r
}
