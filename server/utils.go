package server

import (
	"github.com/gin-gonic/gin"
)

func RespondJSON(c *gin.Context, httpCode int, status, message string, data gin.H) {
	response := gin.H{
		"status":  status,
		"message": message,
		"code":    httpCode,
	}
	for k, v := range data {
		response[k] = v
	}
	c.JSON(httpCode, response)
}
