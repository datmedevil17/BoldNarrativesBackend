package router

import (
	"github.com/datmedevil17/backend/controller"
	"github.com/gorilla/mux"
)

func UserRoutes(router *mux.Router) {
	// Public user routes (no authentication required)
	router.HandleFunc("/getuser/{id}", controller.GetUserByID).Methods("GET")
	router.HandleFunc("/signup", controller.UserSignUp).Methods("POST")
	router.HandleFunc("/view/{id}", controller.ViewUserProfile).Methods("GET")
	router.HandleFunc("/profile", controller.GetUserProfileByQuery).Methods("GET")
	router.HandleFunc("/signin", controller.UserSignIn).Methods("POST")

	// Protected user routes (authentication required)
	protectedRouter := router.NewRoute().Subrouter()
	protectedRouter.Use(JWTAuthMiddleware)
	protectedRouter.HandleFunc("/getid", controller.GetAuthenticatedUserID).Methods("GET")
	protectedRouter.HandleFunc("/protected", controller.ProtectedRoute).Methods("GET")
}
