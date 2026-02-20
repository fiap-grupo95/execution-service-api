package routes

import (
	"github.com/gin-gonic/gin"
)

type executionHandler interface {
	CreateExecution(c *gin.Context)
	FinishExecution(c *gin.Context)
	CancelExecution(c *gin.Context)
}

func addExecutionRoutes(rg *gin.Engine, executionHandler executionHandler) {
	rg.POST("/v1/execution", executionHandler.CreateExecution)
	rg.POST("/v1/execution/finish/:id", executionHandler.FinishExecution)
	rg.POST("/v1/execution/cancel/:id", executionHandler.CancelExecution)
}
