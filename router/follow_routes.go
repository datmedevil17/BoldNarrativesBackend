package router

import (
	"github.com/datmedevil17/backend/controller"
	"github.com/gorilla/mux"
)

func FollowRoutes(router *mux.Router) {
	// Follow/unfollow routes
	router.HandleFunc("/follow/check", controller.CheckFollow).Methods("POST")
	router.HandleFunc("/follow/check/{follower_id}/{following_id}", controller.CheckFollow).Methods("GET")
	router.HandleFunc("/follow", controller.Follow).Methods("POST")
	router.HandleFunc("/unfollow", controller.Unfollow).Methods("POST")

	// Followers and Following routes
	router.HandleFunc("/followers/{user_id}", controller.ShowFollowers).Methods("GET")
	router.HandleFunc("/following/{user_id}", controller.ShowFollowing).Methods("GET")
}
