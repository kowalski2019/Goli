package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger returns a Gin middleware that logs HTTP requests with detailed information
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate request duration
		latency := time.Since(start)

		// Get client IP
		clientIP := c.ClientIP()

		// Get method, status code, and error message
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// Build log message
		if raw != "" {
			path = path + "?" + raw
		}

		if errorMessage != "" {
			log.Printf("[%s] %s %s | %d | %v | %s | %s",
				clientIP,
				method,
				path,
				statusCode,
				latency,
				c.Request.UserAgent(),
				errorMessage,
			)
		} else {
			log.Printf("[%s] %s %s | %d | %v | %s",
				clientIP,
				method,
				path,
				statusCode,
				latency,
				c.Request.UserAgent(),
			)
		}
	}
}
