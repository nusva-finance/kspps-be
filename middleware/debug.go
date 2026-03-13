package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware logs all incoming requests with duration
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		log.Printf("[REQUEST] %s | %s | %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP())

		c.Next()

		duration := time.Since(startTime)
		log.Printf("[REQUEST] %s | %s | %s | %dms",
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			duration.Milliseconds())
	}
}

// RequestLoggerMiddleware logs response status after request is processed
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		log.Printf("[RESPONSE] Status: %d | Path: %s | Client: %s",
			c.Writer.Status(),
			c.Request.URL.Path,
			c.ClientIP())
	}
}