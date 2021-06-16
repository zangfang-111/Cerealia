package txsourceimpl

import (
	"context"
	"sync"
	"time"

	driver "github.com/arangodb/go-driver"
	"github.com/robert-zaremba/errstack"
	"github.com/robert-zaremba/log15"
	"github.com/stellar/go/keypair"

	"bitbucket.org/cerealia/apps/go-lib/model"
	txsourcedal "bitbucket.org/cerealia/apps/go-lib/model/dal/txsource"
	"bitbucket.org/cerealia/apps/go-lib/stellar/txsource"
)

var logger = log15.Root()

type locker struct {
	locker       txsourcedal.SourceLocker
	mutex        *sync.Mutex
	lockDuration time.Duration
}

// NewDriver creates a new instance of driver
func NewDriver(db driver.Database, lockDuration time.Duration) txsource.Driver {
	return locker{
		txsourcedal.NewSourceLocker(db),
		&sync.Mutex{},
		lockDuration,
	}
}

// findMainTradeAcc returns main trade account credentials
func (l locker) findMainTradeAcc(ctx context.Context, expectedMainAddr model.SCAddr, poolAcc *model.ParsedTXSourceAcc) (*model.ParsedTXSourceAcc, errstack.E) {
	if poolAcc.PubKey == expectedMainAddr {
		return poolAcc, nil
	}
	return l.locker.FindByKey(ctx, expectedMainAddr)
}

// Find finds a source acc for specific address
func (l locker) Find(ctx context.Context, scAddr model.SCAddr, tradeID, userID string) (*txsource.SourceAccs, error) {
	poolAcc, err := l.locker.Find(ctx, tradeID, userID)
	if err != nil {
		return nil, errstack.WrapAsReqF(err, "Can't find any locked sources for trade")
	}
	tradeMainAcc, err := l.findMainTradeAcc(ctx, scAddr, poolAcc)
	if err != nil {
		return nil, errstack.WrapAsReqF(err, "Can't find main trade account")
	}
	return createSourceAccs(tradeMainAcc, poolAcc), nil
}

func createSourceAccs(tradeAcc, poolAcc *model.ParsedTXSourceAcc) *txsource.SourceAccs {
	return &txsource.SourceAccs{
		TradeKeyPair: tradeAcc.KeyPair,
		PoolAcc:      *poolAcc,
	}
}

// acquireUnsafe locks scAddr and a pool account without using a mutex
func (l locker) acquireUnsafe(ctx context.Context, scAddr model.SCAddr, tradeID, userID string) (*txsource.SourceAccs, error) {
	poolAcc, err := l.locker.Acquire(ctx, scAddr, tradeID, userID, l.lockDuration)
	if err != nil {
		return nil, errstack.WrapAsReqF(err, "Can't get a lock for trade account")
	}
	tradeMainAcc, err := l.findMainTradeAcc(ctx, scAddr, poolAcc)
	if err != nil {
		return nil, errstack.WrapAsReqF(err, "Can't find main trade account")
	}
	return createSourceAccs(tradeMainAcc, poolAcc), nil
}

// Create creates and acquires a new source account
func (l locker) Create(ctx context.Context, key keypair.Full, tradeID, userID string, accType model.TXSourceAccType) (*txsource.SourceAccs, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	poolAcc, err := l.locker.Acquire(ctx, model.SCAddr(key.Address()), tradeID, userID, l.lockDuration)
	if err != nil {
		return nil, errstack.WrapAsReqF(err, "Can't get a lock for trade account")
	}
	freshTradeAcc, err := l.locker.Create(ctx, key, tradeID, userID, accType, l.lockDuration)
	if err != nil {
		return nil, errstack.WrapAsInfF(err, "Couldn't create a source account")
	}
	return createSourceAccs(freshTradeAcc, poolAcc), nil
}

// IsAcquiredFn implements the txsource.Driver interface
func (l locker) IsAcquiredFn(ctx context.Context, tradeID, userID string) func() error {
	return func() error {
		poolSouceAcc, err := l.locker.Find(ctx, tradeID, userID)
		if err != nil {
			return err
		}
		return poolSouceAcc.MustBeValidFor(tradeID, userID)
	}
}

// Acquire tries to acquire a lock for the scAddr
func (l locker) Acquire(ctx context.Context, scAddr model.SCAddr, tradeID, userID string) (*txsource.SourceAccs, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.acquireUnsafe(ctx, scAddr, tradeID, userID)
}

// release implements the txsource.Driver interface
func (l locker) release(ctx context.Context, tradeID, userID string) error {
	return errstack.WrapAsInfF(l.locker.Unlock(ctx, tradeID, userID), "Can't release a pool lock.")
}

// ReleaseFn implements the txsource.Driver interface and calls Release inside of it
func (l locker) ReleaseFn(ctx context.Context, tradeID, userID string) func() error {
	return func() error {
		l.mutex.Lock()
		defer l.mutex.Unlock()
		return l.release(ctx, tradeID, userID)
	}
}
