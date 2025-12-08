package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request already has an ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			// Generate new request ID
			requestID = uuid.New().String()
		}

		// Set request ID in context
		c.Set("RequestID", requestID)
		
		// Add to response headers
		c.Writer.Header().Set("X-Request-ID", requestID)
		
		c.Next()
	}
}
