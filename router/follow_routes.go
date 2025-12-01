package router

import (
	"github.com/datmedevil17/backend/controller"
	"github.com/gorilla/mux"
)

func FollowRoutes(router *mux.Router) {
	// Public follow routes (no authentication required)
	router.HandleFunc("/follow/check", controller.CheckFollow).Methods("POST")
	router.HandleFunc("/follow/check/{follower_id}/{following_id}", controller.CheckFollow).Methods("GET")
	router.HandleFunc("/followers/{user_id}", controller.ShowFollowers).Methods("GET")
	router.HandleFunc("/following/{user_id}", controller.ShowFollowing).Methods("GET")

	// Protected follow routes (authentication required)
	protectedRouter := router.NewRoute().Subrouter()
	protectedRouter.Use(JWTAuthMiddleware)
	protectedRouter.HandleFunc("/follow", controller.Follow).Methods("POST")
	protectedRouter.HandleFunc("/unfollow", controller.Unfollow).Methods("POST")
}
