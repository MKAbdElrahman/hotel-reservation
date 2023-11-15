package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Logger(c *gin.Context) {
	// Log information about the incoming request
	method := c.Request.Method
	path := c.Request.URL.Path

	fmt.Printf("[%s] %s\n", method, path)

	// Pass control to the next middleware or route handler
	c.Next()

	// Log information about the outgoing response
	status := c.Writer.Status()
	fmt.Printf("[%s] %s - Status: %d\n", method, path, status)
}
