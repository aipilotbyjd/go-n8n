package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jaydeep/go-n8n/configs"
	"golang.org/x/time/rate"
)

// RateLimit returns a gin middleware for rate limiting
func RateLimit(cfg configs.RateLimitConfig) gin.HandlerFunc {
	// Create a new rate limiter
	limiter := rate.NewLimiter(
		rate.Every(cfg.Duration/time.Duration(cfg.Requests)),
		cfg.Burst,
	)

	return func(c *gin.Context) {
		// Check if request is allowed
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}
