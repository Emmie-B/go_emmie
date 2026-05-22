package routes

import (
	"fmt"
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
	authService, err := services.NewAuthService(userRepo, cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize auth service: %v", err))
	}

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
			authGroup.POST("/login", authHandler.Login) // Implement login handler similarly
		}

		// Example of a protected route that requires authentication
		protectedGroup := v1.Group("/")
		protectedGroup.Use(middlewares.RequireAuth(cfg))
		{
			protectedGroup.GET("/profile", authHandler.GetProfile)

			// Example of a route that requires a specific role (e.g. admin)
			protectedGroup.GET("/admin/dashboard", middlewares.RequireRole("ADMIN"), (func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Welcome to the admin dashboard!"})
			}))
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
