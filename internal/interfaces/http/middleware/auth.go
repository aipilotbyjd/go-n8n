package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jaydeep/go-n8n/configs"
)

// Auth returns a gin middleware for JWT authentication
func Auth(cfg configs.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		// Check Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.Secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Set user context
			if userID, ok := claims["user_id"].(string); ok {
				c.Set("UserID", userID)
			}
			if email, ok := claims["email"].(string); ok {
				c.Set("Email", email)
			}
			if role, ok := claims["role"].(string); ok {
				c.Set("Role", role)
			}
		}

		c.Next()
	}
}

// RequireRole returns a middleware that checks if user has required role
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("Role")
		if !exists || userRole != role {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}
		c.Next()
	}
}
