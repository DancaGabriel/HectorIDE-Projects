package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a Gin middleware function that checks for a valid JWT token in the Authorization header.
// It expects the token to be in the format "Bearer <token>".
// If the token is missing or invalid, it aborts the request with a 401 Unauthorized error.
// NOTE: This is a placeholder implementation. In a real application, this would verify the JWT token
//       using a library like "github.com/dgrijalva/jwt-go" and extract user information.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format. Expected 'Bearer <token>'"})
			return
		}

		tokenString := authHeader[len("Bearer "):] // Extract the token

		// Placeholder: Validate the token (replace with actual JWT validation logic)
		// For example:
		// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 	// Validate signing method (e.g., HMAC)
		// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		// 		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		// 	}
		// 	return []byte("your-secret-key"), nil // Replace with your secret key
		// })
		//
		// if err != nil {
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		// 	return
		// }

		// if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 	// Access claims (e.g., user ID, roles)
		// 	userID := claims["user_id"].(string)
		// 	c.Set("userID", userID)
		// 	c.Next()
		// } else {
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		// 	return
		// }

		// For now, just log the token and pass the request through.
		fmt.Println("Token received:", tokenString)
		// In a real application, you would extract user information from the token
		// and set it in the context for subsequent handlers to use.
		c.Next()
	}
}