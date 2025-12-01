package router

import (
	"github.com/datmedevil17/backend/controller"
	"github.com/gorilla/mux"
)

func BlogRoutes(router *mux.Router) {
	// Public blog routes (no authentication required)
	router.HandleFunc("/blog/{id}", controller.GetBlogByID).Methods("GET")
	router.HandleFunc("/blogs/time", controller.ListBlogsByTime).Methods("GET")
	router.HandleFunc("/blogs/views", controller.ListBlogsByViews).Methods("GET")
	router.HandleFunc("/blogs/advanced", controller.ListBlogsAdvanced).Methods("GET")
	router.HandleFunc("/blogs/count", controller.GetBlogCount).Methods("GET")
	router.HandleFunc("/sort/trending", controller.GetTrendingBlogs).Methods("GET")
	router.HandleFunc("/total", controller.GetTotalBlogCount).Methods("POST")
	router.HandleFunc("/view/{id}", controller.IncrementViewCount).Methods("POST")

	// Protected blog routes (authentication required)
	protectedRouter := router.NewRoute().Subrouter()
	protectedRouter.Use(JWTAuthMiddleware)
	protectedRouter.HandleFunc("/", controller.CreateBlog).Methods("POST") // Create blog
	protectedRouter.HandleFunc("/update/{id}", controller.UpdateBlog).Methods("PUT")
	protectedRouter.HandleFunc("/delete/{id}", controller.DeleteBlog).Methods("DELETE")
	protectedRouter.HandleFunc("/sort/time/{id}", controller.SortBlogsByTime).Methods("POST")
	protectedRouter.HandleFunc("/sort/time", controller.SortBlogsByTime).Methods("POST")
	protectedRouter.HandleFunc("/sort/views", controller.SortBlogsByViews).Methods("POST")
}
