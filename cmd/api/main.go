package main

import (
	"fmt"
	"log"

	"github.com/datmedevil17/BoldNarrativesBackend/internal/config"
	"github.com/datmedevil17/BoldNarrativesBackend/internal/database"
	"github.com/datmedevil17/BoldNarrativesBackend/internal/handlers/blog"
	"github.com/datmedevil17/BoldNarrativesBackend/internal/handlers/user"
	"github.com/datmedevil17/BoldNarrativesBackend/internal/middleware"
	blogService "github.com/datmedevil17/BoldNarrativesBackend/internal/services/blog"
	userService "github.com/datmedevil17/BoldNarrativesBackend/internal/services/user"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid config: %v", err)
	}
	if err := database.Connect(cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	db := database.GetDB()
	userSvc := userService.NewService(db)
	userHandler := user.NewHandler(userSvc, cfg.JWTSecret)
	blogSvc := blogService.NewService(db)
	blogHandler := blog.NewHandler(blogSvc, cfg.JWTSecret)

	SetUpRoutes(router, userHandler, blogHandler, cfg.JWTSecret)
	addr:=fmt.Sprintf("%s",cfg.Port)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}


}

func SetUpRoutes(router *gin.Engine, userHandler *user.Handler, blogHandler *blog.Handler, jwtSecret string) {
	 api := router.Group("/api")
    
    // User routes
    userRoutes := api.Group("/user")
    {
        // Public routes
        userRoutes.POST("/signup", userHandler.SignUp)
        userRoutes.POST("/signin", userHandler.SignIn)
        userRoutes.GET("/getuser/:id", userHandler.GetUserById)
        
        // Protected routes
        protected := userRoutes.Group("")
        protected.Use(middleware.AuthMiddleware(jwtSecret))
        {
            protected.GET("/getid", userHandler.GetCurrentUserId)
            protected.GET("/view/:id", userHandler.ViewProfile)
            protected.GET("/profile", userHandler.GetProfile)
            protected.POST("/follow/check", userHandler.CheckFollowStatus)
            protected.POST("/follow", userHandler.FollowUser)
            protected.POST("/unfollow", userHandler.UnFollowUser)
            protected.GET("/followers", userHandler.GetFollowers)
            protected.GET("/following", userHandler.GetFollowing)
        }
    }
    
    // Blog routes (all protected)
    blogRoutes := api.Group("/blog")
    blogRoutes.Use(middleware.AuthMiddleware(jwtSecret))
    {
        blogRoutes.POST("", blogHandler.CreateBlog)
        blogRoutes.GET("/blog/:id", blogHandler.GetBlogById)
        blogRoutes.PUT("/update/:id", blogHandler.UpdateBlog)
        blogRoutes.DELETE("/delete/:id", blogHandler.DeleteBlog)
        
        blogRoutes.POST("total", blogHandler.GetTotalCount)
        blogRoutes.POST("/sort/time/:id", blogHandler.SortByTime)
        blogRoutes.POST("/sort/views", blogHandler.SortByViews)
        blogRoutes.GET("/sort/trending", blogHandler.GetTrending)
        
        blogRoutes.PUT("/view", blogHandler.IncrementViews)
        blogRoutes.POST("/vote/check", blogHandler.CheckVote)
        blogRoutes.POST("/vote", blogHandler.ToggleVote)
        
        blogRoutes.POST("/comment", blogHandler.CreateComment)
        blogRoutes.GET("/comment/:id", blogHandler.GetCommentsByBlogId)
        blogRoutes.DELETE("/comment/:id", blogHandler.DeleteComment)
    }
    
    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
	
}
