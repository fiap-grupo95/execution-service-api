package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap-grupo95/execution-service-api/internal/adapter/http/dto/request"
	"github.com/fiap-grupo95/execution-service-api/internal/adapter/http/dto/response"
	"github.com/fiap-grupo95/execution-service-api/internal/usecase/interfaces"
	"github.com/stretchr/testify/require"
)

type stubExecutionGateway struct {
	createFn func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error)
	finishFn func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error)
	cancelFn func(ctx context.Context, executionID string) (*response.ExecutionResponse, error)
}

func (s *stubExecutionGateway) CreateExecution(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
	return s.createFn(ctx, execution)
}
func (s *stubExecutionGateway) FinishExecution(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
	return s.finishFn(ctx, execution)
}
func (s *stubExecutionGateway) CancelExecution(ctx context.Context, executionID string) (*response.ExecutionResponse, error) {
	return s.cancelFn(ctx, executionID)
}

var _ interfaces.IExecutionGateway = (*stubExecutionGateway)(nil)

func TestExecutionUsecase_CreateExecution_NilRequest(t *testing.T) {
	u := NewExecutionUsecase(&stubExecutionGateway{})
	resp, err := u.CreateExecution(context.Background(), nil)
	require.Nil(t, resp)
	require.Error(t, err)
}

func TestExecutionUsecase_FinishExecution_NilRequest(t *testing.T) {
	u := NewExecutionUsecase(&stubExecutionGateway{})
	resp, err := u.FinishExecution(context.Background(), nil)
	require.Nil(t, resp)
	require.Error(t, err)
}

func TestExecutionUsecase_CancelExecution_EmptyID(t *testing.T) {
	u := NewExecutionUsecase(&stubExecutionGateway{})
	resp, err := u.CancelExecution(context.Background(), "")
	require.Nil(t, resp)
	require.Error(t, err)
}

func TestExecutionUsecase_HappyPath_DelegatesToGateway(t *testing.T) {
	gw := &stubExecutionGateway{
		createFn: func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
			require.Equal(t, "so_1", execution.ServiceOrderID)
			return &response.ExecutionResponse{ID: "exe_1"}, nil
		},
		finishFn: func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
			require.Equal(t, "exe_1", execution.ID)
			return &response.ExecutionResponse{ID: "exe_1", Status: "finished"}, nil
		},
		cancelFn: func(ctx context.Context, executionID string) (*response.ExecutionResponse, error) {
			require.Equal(t, "exe_1", executionID)
			return &response.ExecutionResponse{ID: "exe_1", Status: "canceled"}, nil
		},
	}

	u := NewExecutionUsecase(gw)

	_, err := u.CreateExecution(context.Background(), &request.ExecutionRequest{ServiceOrderID: "so_1"})
	require.NoError(t, err)

	_, err = u.FinishExecution(context.Background(), &request.ExecutionRequest{ID: "exe_1", ServiceOrderID: "so_1"})
	require.NoError(t, err)

	_, err = u.CancelExecution(context.Background(), "exe_1")
	require.NoError(t, err)
}

func TestExecutionUsecase_GatewayErrorsPropagate(t *testing.T) {
	gw := &stubExecutionGateway{
		createFn: func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
			return nil, errors.New("create failed")
		},
		finishFn: func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
			return nil, errors.New("finish failed")
		},
		cancelFn: func(ctx context.Context, executionID string) (*response.ExecutionResponse, error) {
			return nil, errors.New("cancel failed")
		},
	}

	u := NewExecutionUsecase(gw)

	resp, err := u.CreateExecution(context.Background(), &request.ExecutionRequest{ServiceOrderID: "so_1"})
	require.Nil(t, resp)
	require.Error(t, err)

	resp, err = u.FinishExecution(context.Background(), &request.ExecutionRequest{ID: "exe_1", ServiceOrderID: "so_1"})
	require.Nil(t, resp)
	require.Error(t, err)

	resp, err = u.CancelExecution(context.Background(), "exe_1")
	require.Nil(t, resp)
	require.Error(t, err)
}
