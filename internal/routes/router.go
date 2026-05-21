package routes

import (
	"go_emmie/internal/config"
	"go_emmie/internal/database"
	"go_emmie/internal/handlers"
	"go_emmie/internal/middlewares"
	"go_emmie/internal/repositories"
	"go_emmie/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func New(cfg *config.Config, db *database.DB) *gin.Engine {
	router := gin.New()

	router.Use(middlewares.Logger())

	// initialize Repo
	userRepo := repositories.NewUserRepository(db)

	// initialize services
	authService := services.NewAuthService(userRepo)

	// init handlers
	authHandler := handlers.NewAuthHandler(authService)

	v1 := router.Group("/api/v1")
	{
		// Add health check endpoint that verifies database connectivity
		v1.GET("/health", healthCheck(db))

		// group
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			// authGroup.POST("/login", authHandler.) // Implement login handler similarly
		}

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
