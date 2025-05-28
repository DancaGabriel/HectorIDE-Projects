package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"taskflow/internal/handlers" // Adjust import path as needed
	"taskflow/internal/middleware"
)

// SetupRouter defines the API routes and their corresponding handlers.
// It takes a gin.Engine instance and injects dependencies into the handlers.
func SetupRouter(router *gin.Engine, userHandler *handlers.UserHandler, postHandler *handlers.PostHandler, commentHandler *handlers.CommentHandler, categoryHandler *handlers.CategoryHandler) {

	// Health Check Endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// User Routes
	router.POST("/users/register", userHandler.RegisterUser)
	router.POST("/users/login", userHandler.LoginUser)
	userGroup := router.Group("/users")
	userGroup.Use(middleware.AuthMiddleware()) // Example auth middleware
	{
		userGroup.GET("/:user_id", userHandler.GetUser)
		userGroup.PUT("/:user_id", userHandler.UpdateUser)
		userGroup.DELETE("/:user_id", userHandler.DeleteUser)
	}

	// Post Routes
	router.POST("/posts", postHandler.CreatePost)
	router.GET("/posts", postHandler.GetPosts)
	router.GET("/posts/:post_id", postHandler.GetPost)
	postGroup := router.Group("/posts")
	postGroup.Use(middleware.AuthMiddleware())
	{
		postGroup.PUT("/:post_id", postHandler.UpdatePost)
		postGroup.DELETE("/:post_id", postHandler.DeletePost)
	}

	// Comment Routes
	router.POST("/posts/:post_id/comments", commentHandler.CreateComment)
	router.GET("/posts/:post_id/comments", commentHandler.GetComments)
	commentGroup := router.Group("/comments")
	commentGroup.Use(middleware.AuthMiddleware())
	{
		commentGroup.PUT("/:comment_id", commentHandler.UpdateComment)
		commentGroup.DELETE("/:comment_id", commentHandler.DeleteComment)
	}

	// Category Routes
	categoryGroup := router.Group("/categories")
	categoryGroup.Use(middleware.AuthMiddleware()) // Example auth middleware for categories (admin only)
	{
		categoryGroup.POST("", categoryHandler.CreateCategory)
		categoryGroup.GET("", categoryHandler.GetCategories)
		categoryGroup.GET("/:category_id", categoryHandler.GetCategory)
		categoryGroup.PUT("/:category_id", categoryHandler.UpdateCategory)
		categoryGroup.DELETE("/:category_id", categoryHandler.DeleteCategory)
	}
}