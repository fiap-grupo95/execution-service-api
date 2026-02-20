package handlers

import (
	"net/http"

	"github.com/fiap-grupo95/execution-service-api/internal/adapter/http/dto/request"
	"github.com/fiap-grupo95/execution-service-api/internal/usecase/interfaces"
	"github.com/gin-gonic/gin"
)

type ExecutionHandler struct {
	usecase interfaces.IExecutionUsecase
}

func NewExecutionHandler(usecase interfaces.IExecutionUsecase) *ExecutionHandler {
	return &ExecutionHandler{usecase: usecase}
}

func (h *ExecutionHandler) CreateExecution(c *gin.Context) {
	var req request.ExecutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	execution, err := h.usecase.CreateExecution(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, execution)
}

func (h *ExecutionHandler) FinishExecution(c *gin.Context) {
	executionID := c.Param("id")

	var req request.ExecutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = executionID
	execution, err := h.usecase.FinishExecution(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, execution)
}

func (h *ExecutionHandler) CancelExecution(c *gin.Context) {
	executionID := c.Param("id")
	execution, err := h.usecase.CancelExecution(c.Request.Context(), executionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, execution)
}
