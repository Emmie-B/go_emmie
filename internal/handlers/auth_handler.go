package handlers

import (
	"go_emmie/internal/services"
	"go_emmie/internal/types"
	"net/http"

	"github.com/gin-gonic/gin"
)


type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	// This is where you would handle the HTTP request, parse the JSON body into a RegisterRequestDTO,
	// call h.authService.RegisterUser, and return the appropriate HTTP response.

	var payload types.RegisterRequestDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	response, err := h.authService.RegisterUser(c.Request.Context(), payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}