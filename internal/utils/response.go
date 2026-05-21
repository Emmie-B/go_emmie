package utils

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}


func SuccessResponse(c *gin.Context, message string, data interface{})  {
	c.JSON(200, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, message string, err error) {
	c.JSON(400, APIResponse{
		Success: false,
		Message: message,
		Error:   err.Error(),
	})
	 
}	
	 