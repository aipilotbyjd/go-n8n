package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jaydeep/go-n8n/pkg/logger"
)

// Logger returns a gin middleware for logging requests
func Logger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log request details
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		log.WithFields(map[string]interface{}{
			"status":     statusCode,
			"latency":    latency,
			"client_ip":  clientIP,
			"method":     method,
			"path":       path,
			"request_id": c.GetString("RequestID"),
		}).Info("Request processed")

		// Log errors if any
		if len(c.Errors) > 0 {
			log.WithFields(map[string]interface{}{
				"errors": c.Errors.String(),
			}).Error("Request failed")
		}
	}
}
