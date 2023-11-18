package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/mkabdelrahman/hotel-reservation/auth"
	"github.com/mkabdelrahman/hotel-reservation/business"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		token, err := auth.ParseToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if !auth.IsTokenNotExpired(token) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Set("userID", claims["id"])
		c.Next()
	}
}

func AdminOnlyMiddleware(manager *business.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user ID from the authentication token or session

		userID, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized HERE"})
			c.Abort()
			return
		}

		// Retrieve user by ID
		user, err := manager.UserStore.GetUserByID(c, userID.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}

		// Check if the user is an admin
		if user == nil || !user.IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied. Admins only"})
			c.Abort()
			return
		}

		// If the user is an admin, continue with the next middleware/handler
		c.Next()
	}
}
