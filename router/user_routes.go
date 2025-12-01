package router

import (
	"github.com/datmedevil17/backend/controller"
	"github.com/gorilla/mux"
)

func UserRoutes(router *mux.Router) {
	// User authentication and profile routes
	router.HandleFunc("/getuser/{id}", controller.GetUserByID).Methods("GET")
	router.HandleFunc("/signup", controller.UserSignUp).Methods("POST")
	router.HandleFunc("/view/{id}", controller.ViewUserProfile).Methods("GET")
	router.HandleFunc("/getid", controller.GetAuthenticatedUserID).Methods("GET")
	router.HandleFunc("/profile", controller.GetUserProfileByQuery).Methods("GET")
	router.HandleFunc("/signin", controller.UserSignIn).Methods("POST")
	router.HandleFunc("/protected", controller.ProtectedRoute).Methods("GET")
}
