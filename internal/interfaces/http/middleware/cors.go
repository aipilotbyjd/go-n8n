package middleware

import (
	"time"
	
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jaydeep/go-n8n/configs"
)

// CORS returns a gin middleware for CORS
func CORS(cfg configs.CORSConfig) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     cfg.AllowedMethods,
		AllowHeaders:     cfg.AllowedHeaders,
		ExposeHeaders:    cfg.ExposedHeaders,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           time.Duration(cfg.MaxAge) * time.Second,
	})
}
