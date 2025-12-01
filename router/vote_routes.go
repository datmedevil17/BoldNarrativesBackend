package router

import (
	"github.com/datmedevil17/backend/controller"
	"github.com/gorilla/mux"
)

func VoteRoutes(router *mux.Router) {
	// Vote operations
	router.HandleFunc("/vote/check", controller.CheckVoteStatus).Methods("POST")
	router.HandleFunc("/vote", controller.ToggleVoteBlog).Methods("POST")
}
