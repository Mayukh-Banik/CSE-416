package rating

import (
	"github.com/gorilla/mux"
)

func InitRatingRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/rating/upvote", handleUpvote).Methods("POST")
	r.HandleFunc("/rating/downvote", handleDownvote).Methods("POST")

	return r
}
