package execution_service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fiap-grupo95/execution-service-api/internal/adapter/http/dto/request"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type fakeSingleResult struct {
	err      error
	decodeFn func(v interface{}) error
}

func (f *fakeSingleResult) Err() error { return f.err }
func (f *fakeSingleResult) Decode(v interface{}) error {
	if f.decodeFn == nil {
		return nil
	}
	return f.decodeFn(v)
}

type mockExecutionCollection struct {
	insertOneFn        func(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error)
	findOneAndUpdateFn func(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) singleResult
}

func (c *mockExecutionCollection) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	if c.insertOneFn == nil {
		return nil, errors.New("insertOne not stubbed")
	}
	return c.insertOneFn(ctx, document)
}

func (c *mockExecutionCollection) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) singleResult {
	if c.findOneAndUpdateFn == nil {
		return &fakeSingleResult{err: errors.New("findOneAndUpdate not stubbed")}
	}
	return c.findOneAndUpdateFn(ctx, filter, update, opts...)
}

func TestExecutionServiceRepository_CreateExecution_NilRequest(t *testing.T) {
	r := &ExecutionServiceRepository{}
	resp, err := r.CreateExecution(context.Background(), nil)
	require.Nil(t, resp)
	require.Error(t, err)
}

func TestExecutionServiceRepository_FinishExecution_InvalidRequest(t *testing.T) {
	r := &ExecutionServiceRepository{}
	resp, err := r.FinishExecution(context.Background(), &request.ExecutionRequest{})
	require.Nil(t, resp)
	require.Error(t, err)
}

func TestExecutionServiceRepository_CancelExecution_EmptyID(t *testing.T) {
	r := &ExecutionServiceRepository{}
	resp, err := r.CancelExecution(context.Background(), "")
	require.Nil(t, resp)
	require.Error(t, err)
}

func TestExecutionServiceRepository_CreateExecution_Success(t *testing.T) {
	now := time.Now()
	r := &ExecutionServiceRepository{collection: &mockExecutionCollection{insertOneFn: func(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
		doc, ok := document.(executionDocument)
		require.True(t, ok)
		require.Equal(t, "so_1", doc.ServiceOrderID)
		require.Equal(t, "started", doc.Status)
		require.NotNil(t, doc.StartedExecutionAt)
		_ = now
		return &mongo.InsertOneResult{InsertedID: doc.ID}, nil
	}}}

	resp, err := r.CreateExecution(context.Background(), &request.ExecutionRequest{ServiceOrderID: "so_1"})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "so_1", resp.ServiceOrderID)
	require.Equal(t, "started", resp.Status)
	require.NotNil(t, resp.StartedAt)
}

func TestExecutionServiceRepository_CreateExecution_InsertError(t *testing.T) {
	r := &ExecutionServiceRepository{collection: &mockExecutionCollection{
		insertOneFn: func(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
			return nil, errors.New("insert failed")
		},
	}}
	resp, err := r.CreateExecution(context.Background(), &request.ExecutionRequest{ServiceOrderID: "so_1"})
	require.Nil(t, resp)
	require.Error(t, err)
}

func TestExecutionServiceRepository_FinishExecution_Success(t *testing.T) {
	now := time.Now()
	r := &ExecutionServiceRepository{collection: &mockExecutionCollection{
		findOneAndUpdateFn: func(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) singleResult {
			return &fakeSingleResult{decodeFn: func(v interface{}) error {
				doc := v.(*executionDocument)
				doc.ID = "exe_1"
				doc.ServiceOrderID = "so_1"
				doc.Status = "finished"
				doc.StartedExecutionAt = &now
				doc.FinalExecutionAt = &now
				return nil
			}}
		},
	}}

	resp, err := r.FinishExecution(context.Background(), &request.ExecutionRequest{ID: "exe_1", ServiceOrderID: "so_1"})
	require.NoError(t, err)
	require.Equal(t, "exe_1", resp.ID)
	require.Equal(t, "finished", resp.Status)
	require.NotNil(t, resp.FinishedAt)
}

func TestExecutionServiceRepository_FinishExecution_FindError(t *testing.T) {
	r := &ExecutionServiceRepository{collection: &mockExecutionCollection{
		findOneAndUpdateFn: func(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) singleResult {
			return &fakeSingleResult{err: errors.New("db error")}
		},
	}}

	resp, err := r.FinishExecution(context.Background(), &request.ExecutionRequest{ID: "exe_1", ServiceOrderID: "so_1"})
	require.Nil(t, resp)
	require.Error(t, err)
}

func TestExecutionServiceRepository_CancelExecution_Success(t *testing.T) {
	now := time.Now()
	r := &ExecutionServiceRepository{collection: &mockExecutionCollection{
		findOneAndUpdateFn: func(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) singleResult {
			return &fakeSingleResult{decodeFn: func(v interface{}) error {
				doc := v.(*executionDocument)
				doc.ID = "exe_1"
				doc.ServiceOrderID = "so_1"
				doc.Status = "canceled"
				doc.FinalExecutionAt = &now
				return nil
			}}
		},
	}}

	resp, err := r.CancelExecution(context.Background(), "exe_1")
	require.NoError(t, err)
	require.Equal(t, "exe_1", resp.ID)
	require.Equal(t, "canceled", resp.Status)
	require.NotNil(t, resp.FinishedAt)
}

func TestExecutionServiceRepository_CancelExecution_FindError(t *testing.T) {
	r := &ExecutionServiceRepository{collection: &mockExecutionCollection{
		findOneAndUpdateFn: func(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) singleResult {
			return &fakeSingleResult{err: errors.New("db error")}
		},
	}}

	resp, err := r.CancelExecution(context.Background(), "exe_1")
	require.Nil(t, resp)
	require.Error(t, err)
}
