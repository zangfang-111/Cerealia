// Package daltest contains dal test stub functions
package daltest

import (
	"context"

	driver "github.com/arangodb/go-driver"
)

// DBStub is a noop stub for arabgodb
type DBStub struct {
}

// Name implements DB interface
func (f DBStub) Name() string {
	return ""
}

// Info implements DB interface
func (f DBStub) Info(ctx context.Context) (driver.DatabaseInfo, error) {
	return driver.DatabaseInfo{}, nil
}

// EngineInfo implements DB interfaace
func (f DBStub) EngineInfo(ctx context.Context) (driver.EngineInfo, error) {
	return driver.EngineInfo{}, nil
}

// Remove implements DB interfaace
func (f DBStub) Remove(ctx context.Context) error {
	return nil
}

// Query implements DB interfaace
func (f DBStub) Query(ctx context.Context, query string, bindVars map[string]interface{}) (driver.Cursor, error) {
	return nil, nil
}

// ValidateQuery implements DB interfaace
func (f DBStub) ValidateQuery(ctx context.Context, query string) error {
	return nil
}

// Transaction implements DB interfaace
func (f DBStub) Transaction(ctx context.Context, action string, options *driver.TransactionOptions) (interface{}, error) {
	return nil, nil
}

// Collection implements DB interfaace
func (f DBStub) Collection(ctx context.Context, name string) (driver.Collection, error) {
	return nil, nil
}

// CollectionExists implements DB interfaace
func (f DBStub) CollectionExists(ctx context.Context, name string) (bool, error) {
	return false, nil
}

// Collections implements DB interfaace
func (f DBStub) Collections(ctx context.Context) ([]driver.Collection, error) {
	return []driver.Collection{}, nil
}

// CreateCollection implements DB interfaace
func (f DBStub) CreateCollection(ctx context.Context, name string, options *driver.CreateCollectionOptions) (driver.Collection, error) {
	return nil, nil
}

// Graph implements DB interfaace
func (f DBStub) Graph(ctx context.Context, name string) (driver.Graph, error) {
	return nil, nil
}

// GraphExists implements DB interfaace
func (f DBStub) GraphExists(ctx context.Context, name string) (bool, error) {
	return false, nil
}

// Graphs implements DB interfaace
func (f DBStub) Graphs(ctx context.Context) ([]driver.Graph, error) {
	return []driver.Graph{}, nil
}

// CreateGraph implements DB interfaace
func (f DBStub) CreateGraph(ctx context.Context, name string, options *driver.CreateGraphOptions) (driver.Graph, error) {
	return nil, nil
}

// View implements DB interfaace
func (f DBStub) View(ctx context.Context, name string) (driver.View, error) {
	return nil, nil
}

// ViewExists implements DB interfaace
func (f DBStub) ViewExists(ctx context.Context, name string) (bool, error) {
	return false, nil
}

// Views implements DB interfaace
func (f DBStub) Views(ctx context.Context) ([]driver.View, error) {
	return []driver.View{}, nil
}

// CreateArangoSearchView implements DB interfaace
func (f DBStub) CreateArangoSearchView(ctx context.Context, name string, options *driver.ArangoSearchViewProperties) (driver.ArangoSearchView, error) {
	return nil, nil
}
