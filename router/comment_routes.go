package router

import (
	"github.com/datmedevil17/backend/controller"
	"github.com/gorilla/mux"
)

func CommentRoutes(router *mux.Router) {
	// Comment operations
	router.HandleFunc("/comment", controller.CreateComment).Methods("POST")
	router.HandleFunc("/comment/{id}", controller.GetComments).Methods("GET")
	router.HandleFunc("/comment/{id}", controller.DeleteComment).Methods("DELETE")
}
