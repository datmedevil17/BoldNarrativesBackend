package router

import (
	"github.com/datmedevil17/backend/controller"
	"github.com/gorilla/mux"
)

func CommentRoutes(router *mux.Router) {
	// Public comment routes (no authentication required)
	router.HandleFunc("/comment/{id}", controller.GetComments).Methods("GET")

	// Protected comment routes (authentication required)
	protectedRouter := router.NewRoute().Subrouter()
	protectedRouter.Use(JWTAuthMiddleware)
	protectedRouter.HandleFunc("/comment", controller.CreateComment).Methods("POST")
	protectedRouter.HandleFunc("/comment/{id}", controller.DeleteComment).Methods("DELETE")
}
