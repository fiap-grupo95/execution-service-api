package execution_service

import (
	"context"
	"errors"
	"time"

	"github.com/fiap-grupo95/execution-service-api/internal/adapter/http/dto/request"
	"github.com/fiap-grupo95/execution-service-api/internal/adapter/http/dto/response"
	"github.com/fiap-grupo95/execution-service-api/internal/infrastructure/database"
	"github.com/fiap-grupo95/execution-service-api/internal/infrastructure/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type singleResult interface {
	Err() error
	Decode(v interface{}) error
}

type executionCollection interface {
	InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error)
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) singleResult
}

type mongoExecutionCollection struct {
	collection *mongo.Collection
}

func (c *mongoExecutionCollection) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	return c.collection.InsertOne(ctx, document)
}

func (c *mongoExecutionCollection) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) singleResult {
	return c.collection.FindOneAndUpdate(ctx, filter, update, opts...)
}

type ExecutionServiceRepository struct {
	collection executionCollection
}

type executionDocument struct {
	ID                 string     `bson:"_id"`
	ServiceOrderID     string     `bson:"os_id"`
	Status             string     `bson:"status"`
	StartedExecutionAt *time.Time `bson:"started_execution_date"`
	FinalExecutionAt   *time.Time `bson:"final_execution_date"`
}

func NewExecutionServiceRepository(db *database.MongoDBConfig) *ExecutionServiceRepository {
	return &ExecutionServiceRepository{
		collection: &mongoExecutionCollection{collection: db.GetCollection("executions")},
	}
}

func (r *ExecutionServiceRepository) CreateExecution(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
	if execution == nil {
		return nil, errors.New("no execution request provided")
	}

	logger := logs.Logger()
	now := time.Now()

	doc := executionDocument{
		ID:                 execution.ID,
		ServiceOrderID:     execution.ServiceOrderID,
		Status:             "started",
		StartedExecutionAt: &now,
		FinalExecutionAt:   nil,
	}

	if doc.ID == "" {
		doc.ID = "exe_" + now.Format("20060102150405.000000000")
	}

	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create execution")
		return nil, err
	}

	return &response.ExecutionResponse{
		ID:             doc.ID,
		ServiceOrderID: doc.ServiceOrderID,
		Status:         doc.Status,
		StartedAt:      doc.StartedExecutionAt,
		FinishedAt:     doc.FinalExecutionAt,
	}, nil
}

func (r *ExecutionServiceRepository) FinishExecution(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
	if execution == nil || execution.ID == "" {
		return nil, errors.New("no execution request or execution ID provided")
	}

	logger := logs.Logger()
	now := time.Now()

	filter := bson.M{"_id": execution.ID}
	update := bson.M{"$set": bson.M{"status": "finished", "final_execution_date": &now}}

	res := r.collection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if err := res.Err(); err != nil {
		logger.Error().Err(err).Msg("failed to finish execution")
		return nil, err
	}

	var doc executionDocument
	if err := res.Decode(&doc); err != nil {
		return nil, err
	}

	return &response.ExecutionResponse{
		ID:             doc.ID,
		ServiceOrderID: doc.ServiceOrderID,
		Status:         doc.Status,
		StartedAt:      doc.StartedExecutionAt,
		FinishedAt:     doc.FinalExecutionAt,
	}, nil
}

func (r *ExecutionServiceRepository) CancelExecution(ctx context.Context, executionID string) (*response.ExecutionResponse, error) {
	if executionID == "" {
		return nil, errors.New("no execution ID provided")
	}

	logger := logs.Logger()
	now := time.Now()
	filter := bson.M{"_id": executionID}
	update := bson.M{"$set": bson.M{"status": "canceled", "final_execution_date": &now}}

	res := r.collection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if err := res.Err(); err != nil {
		logger.Error().Err(err).Msg("failed to cancel execution")
		return nil, err
	}

	var doc executionDocument
	if err := res.Decode(&doc); err != nil {
		return nil, err
	}

	return &response.ExecutionResponse{
		ID:             doc.ID,
		ServiceOrderID: doc.ServiceOrderID,
		Status:         doc.Status,
		StartedAt:      doc.StartedExecutionAt,
		FinishedAt:     doc.FinalExecutionAt,
	}, nil
}
