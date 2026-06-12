package utils

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Code    string      `json:"code,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

func SuccessResponse(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, status int, message string, code string, details interface{}) {
	c.JSON(status, Response{
		Status:  status,
		Message: message,
		Code:    code,
		Details: details,
	})
}
