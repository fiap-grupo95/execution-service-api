package gateway

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap-grupo95/execution-service-api/internal/adapter/http/dto/request"
	"github.com/fiap-grupo95/execution-service-api/internal/adapter/http/dto/response"
	"github.com/fiap-grupo95/execution-service-api/internal/usecase/interfaces"
	"github.com/stretchr/testify/require"
)

type stubExecutionRepo struct {
	createFn func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error)
	finishFn func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error)
	cancelFn func(ctx context.Context, executionID string) (*response.ExecutionResponse, error)
}

func (s *stubExecutionRepo) CreateExecution(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
	return s.createFn(ctx, execution)
}
func (s *stubExecutionRepo) FinishExecution(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
	return s.finishFn(ctx, execution)
}
func (s *stubExecutionRepo) CancelExecution(ctx context.Context, executionID string) (*response.ExecutionResponse, error) {
	return s.cancelFn(ctx, executionID)
}

var _ interfaces.IExecutionRepository = (*stubExecutionRepo)(nil)

func TestExecutionServiceGateway_CreateExecution_Nil(t *testing.T) {
	g := NewExecutionServiceGateway(&stubExecutionRepo{})
	resp, err := g.CreateExecution(context.Background(), nil)
	require.Nil(t, resp)
	require.Error(t, err)
}

func TestExecutionServiceGateway_CreateExecution_RepoError(t *testing.T) {
	repo := &stubExecutionRepo{
		createFn: func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
			return nil, errors.New("fail")
		},
		finishFn: func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) { return nil, nil },
		cancelFn: func(ctx context.Context, executionID string) (*response.ExecutionResponse, error) { return nil, nil },
	}
	g := NewExecutionServiceGateway(repo)
	resp, err := g.CreateExecution(context.Background(), &request.ExecutionRequest{ServiceOrderID: "so"})
	require.Nil(t, resp)
	require.Error(t, err)
}

func TestExecutionServiceGateway_CreateExecution_Success(t *testing.T) {
	repo := &stubExecutionRepo{
		createFn: func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
			require.Equal(t, "so", execution.ServiceOrderID)
			return &response.ExecutionResponse{ID: "exe"}, nil
		},
		finishFn: func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) { return nil, nil },
		cancelFn: func(ctx context.Context, executionID string) (*response.ExecutionResponse, error) { return nil, nil },
	}
	g := NewExecutionServiceGateway(repo)
	resp, err := g.CreateExecution(context.Background(), &request.ExecutionRequest{ServiceOrderID: "so"})
	require.NoError(t, err)
	require.Equal(t, "exe", resp.ID)
}

func TestExecutionServiceGateway_FinishExecution_Nil(t *testing.T) {
	g := NewExecutionServiceGateway(&stubExecutionRepo{})
	resp, err := g.FinishExecution(context.Background(), nil)
	require.Nil(t, resp)
	require.Error(t, err)
}

func TestExecutionServiceGateway_FinishExecution_Success(t *testing.T) {
	repo := &stubExecutionRepo{
		createFn: func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) { return nil, nil },
		finishFn: func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
			require.Equal(t, "exe", execution.ID)
			return &response.ExecutionResponse{ID: "exe", Status: "finished"}, nil
		},
		cancelFn: func(ctx context.Context, executionID string) (*response.ExecutionResponse, error) { return nil, nil },
	}
	g := NewExecutionServiceGateway(repo)
	resp, err := g.FinishExecution(context.Background(), &request.ExecutionRequest{ID: "exe"})
	require.NoError(t, err)
	require.Equal(t, "finished", resp.Status)
}

func TestExecutionServiceGateway_CancelExecution_EmptyID(t *testing.T) {
	g := NewExecutionServiceGateway(&stubExecutionRepo{})
	resp, err := g.CancelExecution(context.Background(), "")
	require.Nil(t, resp)
	require.Error(t, err)
}

func TestExecutionServiceGateway_CancelExecution_Success(t *testing.T) {
	repo := &stubExecutionRepo{
		createFn: func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) { return nil, nil },
		finishFn: func(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) { return nil, nil },
		cancelFn: func(ctx context.Context, executionID string) (*response.ExecutionResponse, error) {
			require.Equal(t, "exe", executionID)
			return &response.ExecutionResponse{ID: executionID, Status: "canceled"}, nil
		},
	}
	g := NewExecutionServiceGateway(repo)
	resp, err := g.CancelExecution(context.Background(), "exe")
	require.NoError(t, err)
	require.Equal(t, "canceled", resp.Status)
}
