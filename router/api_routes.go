package router

import (
	"github.com/gorilla/mux"
)

// SetupAPIRoutes creates versioned API routes
func SetupAPIRoutes() *mux.Router {
	router := mux.NewRouter()

	// Apply global middleware
	router.Use(CORSMiddleware)
	router.Use(LoggingMiddleware)

	// API v1 routes
	api := router.PathPrefix("/api").Subrouter()
	v1 := api.PathPrefix("/v1").Subrouter()

	// Register versioned route groups
	UserRoutes(v1)
	FollowRoutes(v1)
	BlogRoutes(v1)
	VoteRoutes(v1)
	CommentRoutes(v1)

	// Health check endpoint
	router.HandleFunc("/health", HealthCheck).Methods("GET")

	return router
}
