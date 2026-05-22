package middlewares

import (
	"go_emmie/internal/config"
	"go_emmie/internal/types"
	"go_emmie/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequireAuth(cfg *config.Config) gin.HandlerFunc {

	return func(c *gin.Context) {
		// This is where you would implement JWT token validation logic.
		// For example, you could extract the token from the Authorization header,
		// validate it, and then set the user information in the context for downstream handlers to use.

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		// Check and split the "Bearer <token>" formatting rule cleanly
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must use Bearer schema"})
			return
		}

		tokenString := parts[1]
		claims, err := utils.VerifyToken(tokenString, cfg.JWT.Secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Session expired or invalid token signature"})
			return
		}

		// set Context values for downstream handlers to use (e.g. user ID, email, role)
		c.Set(types.UserIDContextKey, claims.UserID)
		c.Set(types.EmailContextKey, claims.Email)
		c.Set(types.RoleContextKey, claims.Role)

		c.Next()
	}

}



func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get(types.RoleContextKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "user role not found in context"})
			return
		}

		userRole, ok := roleValue.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid role type in context"})
			return
		}

		if userRole != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		c.Next()
	}
}