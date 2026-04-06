package utils

import "github.com/gin-gonic/gin"

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func SuccessResponse(context *gin.Context, statusCode int, message string, data interface{}) {
  context.JSON(statusCode, Response {
    Success: true,
    Message: message,
    Data: data,
  })
}

func ErrorResponse(context *gin.Context, statusCode int, message string) {
  context.JSON(statusCode, Response {
    Success: false,
    Error: message,
  })
}