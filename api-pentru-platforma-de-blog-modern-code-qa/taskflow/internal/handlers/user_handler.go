package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Mock User data (Replace with actual database interaction)
var (
	users = []User{
		{ID: 1, Username: "johndoe", Email: "john.doe@example.com", FirstName: "John", LastName: "Doe"},
		{ID: 2, Username: "janedoe", Email: "jane.doe@example.com", FirstName: "Jane", LastName: "Doe"},
	}
	nextUserID = 3
)

// User struct
type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"` // Omit password from JSON responses (should be hashed in reality)
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

// Helper function to send JSON responses
func respondWithJSON(c *gin.Context, code int, payload interface{}) {
	c.Header("Content-Type", "application/json")
	c.JSON(code, payload)
}

// Helper function to send error responses
func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}

// RegisterUserHandler handles user registration.
func RegisterUserHandler(c *gin.Context) {
	var user User

	if err := c.BindJSON(&user); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Basic validation (replace with more robust validation)
	if user.Username == "" || user.Email == "" || user.Password == "" {
		respondWithError(c, http.StatusBadRequest, "Missing required fields")
		return
	}

	// In a real application, password would be hashed here.
	user.ID = nextUserID
	nextUserID++
	users = append(users, user)
	respondWithJSON(c, http.StatusCreated, user)
}

// LoginUserHandler handles user login.
func LoginUserHandler(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Find user (replace with database lookup)
	for _, u := range users {
		if u.Email == user.Email && u.Password == user.Password { // In a real application, password would be compared against the hashed value.
			// Generate JWT token (placeholder)
			token := "fake_jwt_token" // Replace with actual JWT generation
			respondWithJSON(c, http.StatusOK, gin.H{"token": token})
			return
		}
	}

	respondWithError(c, http.StatusUnauthorized, "Invalid credentials")
}

// GetUserHandler retrieves a user by ID.
func GetUserHandler(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	for _, u := range users {
		if u.ID == userID {
			respondWithJSON(c, http.StatusOK, u)
			return
		}
	}

	respondWithError(c, http.StatusNotFound, "User not found")
}

// UpdateUserHandler updates user information.
func UpdateUserHandler(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var updatedUser User
	if err := c.BindJSON(&updatedUser); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	for i, u := range users {
		if u.ID == userID {
			// Simple update (replace with more robust update logic)
			updatedUser.ID = userID
			users[i] = updatedUser
			respondWithJSON(c, http.StatusOK, updatedUser)
			return
		}
	}

	respondWithError(c, http.StatusNotFound, "User not found")
}

// DeleteUserHandler deletes a user.
func DeleteUserHandler(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	for i, u := range users {
		if u.ID == userID {
			// Delete user
			users = append(users[:i], users[i+1:]...)
			c.Status(http.StatusNoContent)
			return
		}
	}

	respondWithError(c, http.StatusNotFound, "User not found")
}

// ListUsersHandler retrieves all users.  Paginated.
func ListUsersHandler(c *gin.Context) {
	// Pagination (example)
	pageStr := c.DefaultQuery("page", "1")   // Default to page 1
	limitStr := c.DefaultQuery("limit", "10") // Default to 10 items per page

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		respondWithError(c, http.StatusBadRequest, "Invalid page number")
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		respondWithError(c, http.StatusBadRequest, "Invalid limit number")
		return
	}

	start := (page - 1) * limit
	end := start + limit

	if start > len(users) {
		respondWithJSON(c, http.StatusOK, []User{}) // Empty slice if page is out of range
		return
	}

	if end > len(users) {
		end = len(users)
	}

	respondWithJSON(c, http.StatusOK, users[start:end])
}

// GetUserByEmailHandler retrieves a user by email.
func GetUserByEmailHandler(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		respondWithError(c, http.StatusBadRequest, "Email is required")
		return
	}

	for _, u := range users {
		if u.Email == email {
			respondWithJSON(c, http.StatusOK, u)
			return
		}
	}

	respondWithError(c, http.StatusNotFound, "User not found with that email")
}

// Mock Authentication (Replace with JWT or other secure auth)
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Example: Check for an API key in the header
		apiKey := c.GetHeader("X-API-Key")
		if apiKey != "supersecretapikey" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized.  Missing or invalid API Key."})
			return
		}
		c.Next()
	}
}

// Admin Middleware (Placeholder - check user roles)
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// In a real application, check the user's role (e.g., from JWT claim).
		// For this example, just allow access.
		// Example:  claims := jwt.ExtractClaims(c) ... check claims["role"] == "admin"

		// Placeholder: Assume the user is an admin
		isAdmin := true //Replace with proper authentication/authorization logic.

		if !isAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden. Admin access required."})
			return
		}

		c.Next()
	}
}

// SetupUserRoutes defines the user-related API endpoints and middleware.
func SetupUserRoutes(router *gin.Engine) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("/register", RegisterUserHandler)
		userGroup.POST("/login", LoginUserHandler)

		// Authentication Required for the following routes
		userGroup.Use(AuthMiddleware())
		{
			userGroup.GET("/", ListUsersHandler)
			userGroup.GET("/email", GetUserByEmailHandler)

			userGroup.GET("/:user_id", GetUserHandler)
			userGroup.PUT("/:user_id", UpdateUserHandler)
			userGroup.DELETE("/:user_id", AdminMiddleware(), DeleteUserHandler) //Admin only route for deletion.
		}
	}

	fmt.Println("User routes configured.")
}