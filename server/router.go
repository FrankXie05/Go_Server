package server

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(TraceIDMiddleware())
	r.POST("/script", RunTasks)
	return r
}
