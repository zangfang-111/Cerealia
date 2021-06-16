package txsource

import (
	"context"

	"bitbucket.org/cerealia/apps/go-lib/model"
)

// NoopDriver is a no operation tx locker
type NoopDriver struct {
}

// Create implements an interface noopsourceacc.Driver
func (n NoopDriver) Create(ctx context.Context, source model.TXSourceAcc) (*model.TXSourceAcc, error) {
	return nil, nil
}

// IsAcquiredFn implements an interface noopsourceacc.Driver
func (n NoopDriver) IsAcquiredFn(ctx context.Context, userID string, scAddr model.SCAddr) func() error {
	return func() error {
		return nil
	}
}

// ReleaseFn implements an interface noopsourceacc.Driver
func (n NoopDriver) ReleaseFn(ctx context.Context, userID string, scAddr model.SCAddr) func() error {
	return func() error {
		return nil
	}
}

// Acquire implements an interface noopsourceacc.Driver
func (n NoopDriver) Acquire(ctx context.Context, userID string, scAddr model.SCAddr) (*model.TXSourceAcc, error) {
	return &model.TXSourceAcc{}, nil
}

// Find implements an interface noopsourceacc.Driver
func (n NoopDriver) Find(ctx context.Context, scAddr model.SCAddr) (*model.TXSourceAcc, error) {
	return nil, nil
}
