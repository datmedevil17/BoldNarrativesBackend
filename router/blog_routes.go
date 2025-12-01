package router

import (
	"github.com/datmedevil17/backend/controller"
	"github.com/gorilla/mux"
)

func BlogRoutes(router *mux.Router) {
	// Basic blog CRUD operations
	router.HandleFunc("/", controller.CreateBlog).Methods("POST") // Create blog
	router.HandleFunc("/blog/{id}", controller.GetBlogByID).Methods("GET")
	router.HandleFunc("/update/{id}", controller.UpdateBlog).Methods("PUT")
	router.HandleFunc("/delete/{id}", controller.DeleteBlog).Methods("DELETE")

	// Blog listing and sorting
	router.HandleFunc("/blogs/time", controller.ListBlogsByTime).Methods("GET")
	router.HandleFunc("/blogs/views", controller.ListBlogsByViews).Methods("GET")
	router.HandleFunc("/blogs/advanced", controller.ListBlogsAdvanced).Methods("GET")
	router.HandleFunc("/blogs/count", controller.GetBlogCount).Methods("GET")

	// Blog sorting with parameters
	router.HandleFunc("/sort/time/{id}", controller.SortBlogsByTime).Methods("POST")
	router.HandleFunc("/sort/time", controller.SortBlogsByTime).Methods("POST") // Without ID
	router.HandleFunc("/sort/views", controller.SortBlogsByViews).Methods("POST")
	router.HandleFunc("/sort/trending", controller.GetTrendingBlogs).Methods("GET")

	// Blog statistics
	router.HandleFunc("/total", controller.GetTotalBlogCount).Methods("POST")

	// View count
	router.HandleFunc("/view", controller.IncrementViewCount).Methods("PUT")
}
