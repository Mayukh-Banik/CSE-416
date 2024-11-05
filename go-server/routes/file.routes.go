package routes

import (
    "github.com/gorilla/mux"
)



/**func UploadFilehandler(w http.ResponseWriter, r *http.Request)
{
	var file_metadata models.File

	if r.Method != http.MethodPost{
		http.Error(w, "Invalid request", http.StatusMethodNotAllowed)
		return 
	}

	if err := json.NewDecoder(r.Body).decode(&file_metadata); err!=nil{
		http.Error(w, "Invalid requests data: %v",err)
		return
	}

	if err :=services.StoreFileMetaData(file_metadata); err!=nil{
		http.Error(w, "Failed to save file metadata", http.StatusInternalServerError)
		log.Printf("Failed to save file metadata: %v",err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File metadata sucessfully uploaded"))
}


func FileRoutes(router *mux.Router) {
	router.HandleFunc("/uploadfile", ).Methods("POST")
}*/