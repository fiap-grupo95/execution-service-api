package mocks

import (
	"context"

	"github.com/fiap-grupo95/execution-service-api/internal/adapter/http/dto/request"
	"github.com/fiap-grupo95/execution-service-api/internal/adapter/http/dto/response"
	"github.com/stretchr/testify/mock"
)

type MockExecutionGateway struct {
	mock.Mock
}

func (m *MockExecutionGateway) CreateExecution(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
	args := m.Called(ctx, execution)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ExecutionResponse), args.Error(1)
}

func (m *MockExecutionGateway) FinishExecution(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
	args := m.Called(ctx, execution)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ExecutionResponse), args.Error(1)
}

func (m *MockExecutionGateway) CancelExecution(ctx context.Context, executionID string) (*response.ExecutionResponse, error) {
	args := m.Called(ctx, executionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ExecutionResponse), args.Error(1)
}
