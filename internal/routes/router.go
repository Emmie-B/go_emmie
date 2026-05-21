package routes

import (
	"go_emmie/internal/config"
	"go_emmie/internal/database"
	"go_emmie/internal/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func New(cfg *config.Config, db *database.DB) *gin.Engine{
	router := gin.New()


	 

	router.Use(middlewares.Logger())

	v1 := router.Group("/api/v1")	

	{
		// Placeholder for future API routes
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		// Add health check endpoint that verifies database connectivity
		v1.GET("/health", healthCheck(db))

		// STATUS CHECK
		v1.GET("/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
			})
		})
	}

	return router
}

// healthCheck returns a health check handler.
func healthCheck(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := "ok"
		dbStatus := "connected"

		if err := db.HealthCheck(c.Request.Context()); err != nil {
			status = "degraded"
			dbStatus = "disconnected"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   status,
			"database": dbStatus,
			"version":  "1.0.0",
		})
	}
}
