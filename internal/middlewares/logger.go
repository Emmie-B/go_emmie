package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)


func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Placeholder for logging middleware implementation
		// Generate request ID
		requestID := c.GetHeader("X-Request-Id")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("requestId", requestID)
		c.Writer.Header().Set("X-Request-Id", requestID)

		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)
		statusCode := c.Writer.Status()

		// Log request
		logger := log.With().
			Str("request_id", requestID).
			Str("method", c.Request.Method).
			Str("path", path).
			Str("query", query).
			Int("status", statusCode).
			Dur("latency", latency).
			Str("client_ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Logger()

		// Get user ID if available
		if userID, exists := c.Get("userId"); exists {
			logger = logger.With().Str("user_id", userID.(string)).Logger()
		}

		switch {
		case statusCode >= 500:
			logger.Error().Msg("Server error")
		case statusCode >= 400:
			logger.Warn().Msg("Client error")
		case statusCode >= 300:
			logger.Info().Msg("Redirect")
		default:
			logger.Info().Msg("Request completed")
		}
	
	}
}