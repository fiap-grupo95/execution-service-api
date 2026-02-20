package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fiap-grupo95/execution-service-api/internal/adapter/http/dto/request"
	"github.com/fiap-grupo95/execution-service-api/internal/adapter/http/dto/response"
	"github.com/fiap-grupo95/execution-service-api/internal/usecase/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type stubExecutionUsecase struct {
	createFn func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error)
	finishFn func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error)
	cancelFn func(ctx context.Context, executionID string) (*response.ExecutionResponse, error)
}

func (s *stubExecutionUsecase) CreateExecution(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
	return s.createFn(ctx, execution)
}
func (s *stubExecutionUsecase) FinishExecution(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
	return s.finishFn(ctx, execution)
}
func (s *stubExecutionUsecase) CancelExecution(ctx context.Context, executionID string) (*response.ExecutionResponse, error) {
	return s.cancelFn(ctx, executionID)
}

var _ interfaces.IExecutionUsecase = (*stubExecutionUsecase)(nil)

func TestExecutionHandler_CreateExecution_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewExecutionHandler(&stubExecutionUsecase{
		createFn: func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
			return nil, nil
		},
		finishFn: func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
			return nil, nil
		},
		cancelFn: func(ctx context.Context, executionID string) (*response.ExecutionResponse, error) {
			return nil, nil
		},
	})

	r := gin.New()
	r.POST("/v1/execution", h.CreateExecution)

	req := httptest.NewRequest(http.MethodPost, "/v1/execution", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestExecutionHandler_CreateExecution_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewExecutionHandler(&stubExecutionUsecase{
		createFn: func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
			return nil, errors.New("boom")
		},
		finishFn: func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
			return nil, nil
		},
		cancelFn: func(ctx context.Context, executionID string) (*response.ExecutionResponse, error) {
			return nil, nil
		},
	})

	r := gin.New()
	r.POST("/v1/execution", h.CreateExecution)

	body, _ := json.Marshal(request.ExecutionRequest{ServiceOrderID: "so_1"})
	req := httptest.NewRequest(http.MethodPost, "/v1/execution", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestExecutionHandler_CreateExecution_Created(t *testing.T) {
	gin.SetMode(gin.TestMode)

	now := time.Now()
	h := NewExecutionHandler(&stubExecutionUsecase{
		createFn: func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
			require.NotNil(t, execution)
			require.Equal(t, "so_1", execution.ServiceOrderID)
			return &response.ExecutionResponse{ID: "exe_1", ServiceOrderID: execution.ServiceOrderID, Status: "started", StartedAt: &now}, nil
		},
		finishFn: func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
			return nil, nil
		},
		cancelFn: func(ctx context.Context, executionID string) (*response.ExecutionResponse, error) {
			return nil, nil
		},
	})

	r := gin.New()
	r.POST("/v1/execution", h.CreateExecution)

	body, _ := json.Marshal(request.ExecutionRequest{ServiceOrderID: "so_1"})
	req := httptest.NewRequest(http.MethodPost, "/v1/execution", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestExecutionHandler_FinishExecution_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewExecutionHandler(&stubExecutionUsecase{
		createFn: func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) { return nil, nil },
		finishFn: func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) { return nil, nil },
		cancelFn: func(ctx context.Context, executionID string) (*response.ExecutionResponse, error) { return nil, nil },
	})

	r := gin.New()
	r.POST("/v1/execution/finish/:id", h.FinishExecution)

	req := httptest.NewRequest(http.MethodPost, "/v1/execution/finish/exe_1", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestExecutionHandler_FinishExecution_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	now := time.Now()
	h := NewExecutionHandler(&stubExecutionUsecase{
		createFn: func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) { return nil, nil },
		finishFn: func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
			require.NotNil(t, execution)
			require.Equal(t, "exe_1", execution.ID)
			require.Equal(t, "so_1", execution.ServiceOrderID)
			return &response.ExecutionResponse{ID: execution.ID, ServiceOrderID: execution.ServiceOrderID, Status: "finished", FinishedAt: &now}, nil
		},
		cancelFn: func(ctx context.Context, executionID string) (*response.ExecutionResponse, error) { return nil, nil },
	})

	r := gin.New()
	r.POST("/v1/execution/finish/:id", h.FinishExecution)

	body, _ := json.Marshal(request.ExecutionRequest{ServiceOrderID: "so_1"})
	req := httptest.NewRequest(http.MethodPost, "/v1/execution/finish/exe_1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestExecutionHandler_CancelExecution_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewExecutionHandler(&stubExecutionUsecase{
		createFn: func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) { return nil, nil },
		finishFn: func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) { return nil, nil },
		cancelFn: func(ctx context.Context, executionID string) (*response.ExecutionResponse, error) {
			return nil, errors.New("boom")
		},
	})

	r := gin.New()
	r.POST("/v1/execution/cancel/:id", h.CancelExecution)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/execution/cancel/exe_1", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestExecutionHandler_CancelExecution_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	now := time.Now()
	h := NewExecutionHandler(&stubExecutionUsecase{
		createFn: func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) { return nil, nil },
		finishFn: func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) { return nil, nil },
		cancelFn: func(ctx context.Context, executionID string) (*response.ExecutionResponse, error) {
			require.Equal(t, "exe_1", executionID)
			return &response.ExecutionResponse{ID: executionID, Status: "canceled", FinishedAt: &now}, nil
		},
	})

	r := gin.New()
	r.POST("/v1/execution/cancel/:id", h.CancelExecution)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/execution/cancel/exe_1", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
