package usecase

import (
	"context"
	"errors"

	"github.com/fiap-grupo95/execution-service-api/internal/adapter/http/dto/request"
	"github.com/fiap-grupo95/execution-service-api/internal/adapter/http/dto/response"
	"github.com/fiap-grupo95/execution-service-api/internal/usecase/interfaces"
)

type ExecutionUsecase struct {
	executionGateway interfaces.IExecutionGateway
}

func NewExecutionUsecase(executionGateway interfaces.IExecutionGateway) *ExecutionUsecase {
	return &ExecutionUsecase{executionGateway: executionGateway}
}

func (u *ExecutionUsecase) CreateExecution(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
	if execution == nil {
		return nil, errors.New("no execution request provided")
	}
	return u.executionGateway.CreateExecution(ctx, execution)
}

func (u *ExecutionUsecase) FinishExecution(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
	if execution == nil {
		return nil, errors.New("no execution request provided")
	}
	return u.executionGateway.FinishExecution(ctx, execution)
}

func (u *ExecutionUsecase) CancelExecution(ctx context.Context, executionID string) (*response.ExecutionResponse, error) {
	if executionID == "" {
		return nil, errors.New("no execution ID provided")
	}
	return u.executionGateway.CancelExecution(ctx, executionID)
}
