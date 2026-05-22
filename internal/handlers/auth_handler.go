package handlers

import (
	"errors"
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
		status := http.StatusInternalServerError
		message := "internal server error"

		switch {
		case errors.Is(err, services.ErrEmailAlreadyExists):
			status = http.StatusConflict
			message = "email already exists"
		case errors.Is(err, services.ErrValidation):
			status = http.StatusBadRequest
			message = err.Error()
		}

		c.JSON(status, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusCreated, response)
}


func (h *AuthHandler) Login(c *gin.Context) {
	// Implement login handler similarly to Register, but call h.authService.LoginUser and return the token on success.
	var payload types.LoginRequestDTO

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	response, err := h.authService.LoginUser(c.Request.Context(), payload)
	if err != nil {
		status := http.StatusInternalServerError
		message := "internal server error"

		switch {
		case errors.Is(err, services.ErrValidation):
			status = http.StatusBadRequest
			message = err.Error()
		default:
			message = err.Error() // Return actual error message for unexpected errors (consider security implications)
		}		

		c.JSON(status, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusOK, response)
}


func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Pull fields from context safely
	userID := c.GetString(types.UserIDContextKey)
	email := c.GetString(types.EmailContextKey)

	// Query data directly using the context value
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully unpacked profile secure token details",
		"user_id": userID,
		"email":   email,
	})
}