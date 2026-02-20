package interfaces

import (
	"context"

	"github.com/fiap-grupo95/execution-service-api/internal/adapter/http/dto/request"
	"github.com/fiap-grupo95/execution-service-api/internal/adapter/http/dto/response"
)

type IExecutionGateway interface {
	CreateExecution(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error)
	FinishExecution(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error)
	CancelExecution(ctx context.Context, executionID string) (*response.ExecutionResponse, error)
}

type IExecutionUsecase interface {
	CreateExecution(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error)
	FinishExecution(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error)
	CancelExecution(ctx context.Context, executionID string) (*response.ExecutionResponse, error)
}

type IExecutionRepository interface {
	CreateExecution(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error)
	FinishExecution(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error)
	CancelExecution(ctx context.Context, executionID string) (*response.ExecutionResponse, error)
}
