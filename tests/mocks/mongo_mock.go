package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

// MockMongoSession is a mock implementation of mongo.Session
type MockMongoSession struct {
	mock.Mock
}

func (m *MockMongoSession) StartTransaction(...*mongo.TransactionOptions) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockMongoSession) AbortTransaction(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMongoSession) CommitTransaction(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMongoSession) EndSession(ctx context.Context) {
	m.Called(ctx)
}

// MockMongoClient is a mock implementation of mongo.Client
type MockMongoClient struct {
	mock.Mock
}

func (m *MockMongoClient) StartSession(opts ...*mongo.SessionOptions) (mongo.Session, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(mongo.Session), args.Error(1)
}

// MockMongoDatabase is a mock implementation of *mongo.Database
type MockMongoDatabase struct {
	mock.Mock
}

func (m *MockMongoDatabase) Client() *mongo.Client {
	args := m.Called()
	return args.Get(0).(*mongo.Client)
}

func (m *MockMongoDatabase) Collection(name string, opts ...*mongo.CollectionOptions) *mongo.Collection {
	args := m.Called(name)
	return args.Get(0).(*mongo.Collection)
}
