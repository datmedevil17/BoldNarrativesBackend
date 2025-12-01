package router

import (
	"github.com/datmedevil17/backend/controller"
	"github.com/gorilla/mux"
)

func VoteRoutes(router *mux.Router) {
	// Public vote routes (no authentication required)
	router.HandleFunc("/vote/check", controller.CheckVoteStatus).Methods("POST")

	// Protected vote routes (authentication required)
	protectedRouter := router.NewRoute().Subrouter()
	protectedRouter.Use(JWTAuthMiddleware)
	protectedRouter.HandleFunc("/vote", controller.ToggleVoteBlog).Methods("POST")
}
