package gateway

import (
	"context"
	"errors"

	"github.com/fiap-grupo95/execution-service-api/internal/adapter/http/dto/request"
	"github.com/fiap-grupo95/execution-service-api/internal/adapter/http/dto/response"
	"github.com/fiap-grupo95/execution-service-api/internal/infrastructure/logs"
	"github.com/fiap-grupo95/execution-service-api/internal/usecase/interfaces"
)

type ExecutionServiceGateway struct {
	executionRepository interfaces.IExecutionRepository
}

func NewExecutionServiceGateway(executionRepository interfaces.IExecutionRepository) *ExecutionServiceGateway {
	return &ExecutionServiceGateway{
		executionRepository: executionRepository,
	}
}

func (g *ExecutionServiceGateway) CreateExecution(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
	logger := logs.Logger()
	if execution == nil {
		return nil, errors.New("no execution request provided")
	}

	ExecutionResponse, err := g.executionRepository.CreateExecution(ctx, execution)
	if err != nil {
		logger.Error().Err(err).Msg("error creating execution")
		return nil, err
	}

	return ExecutionResponse, nil
}

func (g *ExecutionServiceGateway) CancelExecution(ctx context.Context, executionID string) (*response.ExecutionResponse, error) {
	logger := logs.Logger()
	if executionID == "" {
		return nil, errors.New("no execution ID provided")
	}

	ExecutionResponse, err := g.executionRepository.CancelExecution(ctx, executionID)
	if err != nil {
		logger.Error().Err(err).Msg("error canceling execution")
		return nil, err
	}

	return ExecutionResponse, nil
}

func (g *ExecutionServiceGateway) FinishExecution(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
	logger := logs.Logger()
	if execution == nil {
		return nil, errors.New("no execution request provided")
	}

	ExecutionResponse, err := g.executionRepository.FinishExecution(ctx, execution)
	if err != nil {
		logger.Error().Err(err).Msg("error finishing execution")
		return nil, err
	}

	return ExecutionResponse, nil
}
