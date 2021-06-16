package txsource

import (
	"context"
	"fmt"
	"time"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dal"
	"bitbucket.org/cerealia/apps/go-lib/model/dbconst"
	driver "github.com/arangodb/go-driver"
	"github.com/robert-zaremba/errstack"
	"github.com/stellar/go/keypair"
)

// SourceLocker provides methods for contract lock checking
type SourceLocker interface {
	// Create creates a new source account
	Create(ctx context.Context, key keypair.Full, tradeID, userID string, accType model.TXSourceAccType, lockDuration time.Duration) (*model.ParsedTXSourceAcc, errstack.E)
	// Find finds a locked source account of the user
	Find(ctx context.Context, tradeID, userID string) (*model.ParsedTXSourceAcc, errstack.E)
	// FindByKey finds a source account's secret key using it's public key
	FindByKey(ctx context.Context, scAddr model.SCAddr) (*model.ParsedTXSourceAcc, errstack.E)
	// Acquire acquires a source account lock using specific account's public address
	Acquire(ctx context.Context, tradeSCAddr model.SCAddr, tradeID, userID string, duration time.Duration) (*model.ParsedTXSourceAcc, errstack.E)
	// Unlock unlocks a specific account
	Unlock(ctx context.Context, tradeID, userID string) errstack.E
}

// NewSourceLocker creates a new lock check model
func NewSourceLocker(db driver.Database) SourceLocker {
	return sourceLocker{db}
}

type sourceLocker struct {
	db driver.Database
}

func createUnlockTime(duration time.Duration) time.Time {
	return time.Now().Add(duration).UTC()
}

// FindByKey finds trade's secret key
func (l sourceLocker) FindByKey(ctx context.Context, scAddr model.SCAddr) (*model.ParsedTXSourceAcc, errstack.E) {
	var cl model.TXSourceAcc
	err := dal.DBGetOneFromColl(ctx, &cl, string(scAddr), dbconst.ColTxSourceAccs, l.db)
	if err != nil {
		return nil, err.WithMsg(fmt.Sprintf("Can't find credentials for address '%s'", scAddr))
	}
	return cl.Parse()
}

// Find finds user's lock by it's context (trade and user's ID)
func (l sourceLocker) Find(ctx context.Context, tradeID, userID string) (*model.ParsedTXSourceAcc, errstack.E) {
	acc := model.TXSourceAcc{}
	query := fmt.Sprintf(`
return FIRST(
	FOR l IN %s
			FILTER l.lockUserID == @lockUserID && l.lockTradeID == @lockTradeID && l.lockUnlockedAt == null && @now < l.lockExpiresAt
			RETURN l
)
`, dbconst.ColTxSourceAccs)
	bindVars := map[string]interface{}{
		"now":         time.Now().UTC(),
		"lockTradeID": tradeID,
		"lockUserID":  userID,
	}
	err := dal.DBQueryFirst(ctx, &acc, query, bindVars, l.db)
	if err != nil {
		return nil, err.WithMsg(fmt.Sprintf("Couldn't find user's '%s' lock for trade '%s'", userID, tradeID))
	}
	return acc.Parse()
}

// Create fetches doc from DB
func (l sourceLocker) Create(ctx context.Context, key keypair.Full, tradeID, userID string, accType model.TXSourceAccType, lockDuration time.Duration) (*model.ParsedTXSourceAcc, errstack.E) {
	sourceAcc := model.TXSourceAcc{
		PubKey:        model.SCAddr(key.Address()),
		Type:          accType,
		SCSecret:      model.SCSecret(key.Seed()),
		LockExpiresAt: createUnlockTime(lockDuration),
		LockTradeID:   tradeID,
		LockUserID:    userID,
	}
	_, err := dal.InsertAny(ctx, dbconst.ColTxSourceAccs, &sourceAcc, l.db)
	if err != nil {
		return nil, err.WithMsg(fmt.Sprintf("Couldn't create a new source account: '%s'", key.Address()))
	}
	return sourceAcc.Parse()
}

// Acquire acquires a lock for a specific account (if it's allowed)
func (l sourceLocker) Acquire(ctx context.Context, tradeSCAddr model.SCAddr, tradeID, userID string, duration time.Duration) (*model.ParsedTXSourceAcc, errstack.E) {
	lock := model.TXSourceAcc{}
	query := fmt.Sprintf(`
LET requestableLock = FIRST(
    FOR l IN %s
        FILTER (l._key == @pubKey || l.type == "pool") && (l.lockExpiresAt == null || l.lockExpiresAt < @now || (l.lockUnlockedAt != null && l.lockUnlockedAt < @now) || (l.lockUserID == @lockUserID && l.lockTradeID == @lockTradeID) || l.lockUserID == null || l.lockTradeID == null)
				SORT l.type != "trade"
        RETURN l
)
UPDATE requestableLock
WITH {lockExpiresAt: @lockExpiresAt, lockUserID: @lockUserID, lockTradeID: @lockTradeID, lockUnlockedAt: null}
IN %s
RETURN NEW
`, dbconst.ColTxSourceAccs, dbconst.ColTxSourceAccs)
	bindVars := map[string]interface{}{
		"pubKey":        string(tradeSCAddr),
		"now":           time.Now().UTC(),
		"lockExpiresAt": createUnlockTime(duration),
		"lockTradeID":   tradeID,
		"lockUserID":    userID,
	}
	err := dal.DBQueryFirst(ctx, &lock, query, bindVars, l.db)
	if err != nil {
		return nil, errstack.WrapAsInfF(err, "Couldn't acquire a lock for an account '%s'. Source account pool is exhausted.", tradeSCAddr)
	}
	return lock.Parse()
}

// Unlock will unlock all user's locks for the given trade
func (l sourceLocker) Unlock(ctx context.Context, tradeID, userID string) errstack.E {
	query := fmt.Sprintf(`
FOR a IN %s
  FILTER a.lockUserID == @lockUserID && a.lockTradeID == @lockTradeID
  UPDATE a
  WITH { lockUnlockedAt: @now }
  IN %s
  RETURN NEW
`, dbconst.ColTxSourceAccs, dbconst.ColTxSourceAccs)
	bindVars := map[string]interface{}{
		"now":         time.Now().UTC(),
		"lockTradeID": tradeID,
		"lockUserID":  userID,
	}
	return dal.DBExec(ctx, query, bindVars, l.db)
}
