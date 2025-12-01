package router

import (
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	// Apply middleware
	router.Use(CORSMiddleware)
	router.Use(LoggingMiddleware)

	// Register route groups
	UserRoutes(router)
	FollowRoutes(router)
	BlogRoutes(router)
	VoteRoutes(router)
	CommentRoutes(router)

	return router
}
