package txsource

import (
	"context"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"github.com/stellar/go/keypair"
)

// Driver is an interface that hides DB object of the source account
type Driver interface {
	Create(ctx context.Context, key keypair.Full, tradeID, userID string, accType model.TXSourceAccType) (*SourceAccs, error)
	Acquire(ctx context.Context, scAddr model.SCAddr, tradeID, userID string) (*SourceAccs, error)
	// ReleaseFn creates a no-arg function that checks if the lock is acquired
	IsAcquiredFn(ctx context.Context, tradeID, userID string) func() error
	Find(ctx context.Context, scAddr model.SCAddr, tradeID, userID string) (*SourceAccs, error)
	// ReleaseFn creates a no-arg function that releases the lock
	ReleaseFn(ctx context.Context, tradeID, userID string) func() error
}

// SourceAccs represents a lock of trade or pool source accounts
type SourceAccs struct {
	TradeKeyPair keypair.Full
	PoolAcc      model.ParsedTXSourceAcc // Can be a pool or a trade account. Determined by vacancy
}
