package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// SecurityHeadersMiddleware adds necessary security headers
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Permissive CSP untuk development
		c.Header("Content-Security-Policy",
			"default-src 'self' http://localhost:5173* 'unsafe-inline' 'unsafe-eval'; "+
				"script-src 'self' 'unsafe-inline' 'unsafe-eval' http://localhost:5173*; "+
				"object-src 'none'; base-uri 'self'; frame-ancestors 'none';")

		c.Next()
	}
}