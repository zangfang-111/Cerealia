package txsource

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/robert-zaremba/flag"
	"github.com/stellar/go/keypair"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dal"
	"bitbucket.org/cerealia/apps/go-lib/model/dbconst"
	dbs "bitbucket.org/cerealia/apps/go-lib/setup/arangodb"

	driver "github.com/arangodb/go-driver"
	. "github.com/robert-zaremba/checkers"
	. "gopkg.in/check.v1"
)

var scAddr1 = model.SCAddr("abc123")
var scAddr2 = model.SCAddr("qwe456")

const lockDuration = (4 * time.Minute)

type TxSourceAccsSuite struct {
	ctx               context.Context
	randomString      string
	locker            SourceLocker
	db                driver.Database
	tradeID           string
	tradeAcc, poolAcc *keypair.Full
}

func init() {
	_ = flag.String("check.f", "", "Testing selector flag")
	flag.Parse()
}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&TxSourceAccsSuite{})

func (s *TxSourceAccsSuite) SetUpSuite(c *C) {
	var err error
	s.ctx = context.Background()
	s.db, err = dbs.GetDb(s.ctx)
	s.locker = sourceLocker{db: s.db}
	c.Assert(err, IsNil)
}

func (s *TxSourceAccsSuite) TestFindByKeyUnknown(c *C) {
	full, err := keypair.Random()
	c.Assert(err, IsNil)
	_, err = s.locker.FindByKey(s.ctx, model.SCAddr(full.Address()))
	c.Assert(err, ErrorContains, "not found")
}

func (s *TxSourceAccsSuite) TestFindByKeyUnclaimed(c *C) {
	acc, err := s.locker.FindByKey(s.ctx, model.SCAddr(s.tradeAcc.Address()))
	c.Assert(err, IsNil)
	c.Assert(acc, NotNil)
	c.Check(acc.PubKey, Equals, model.SCAddr(s.tradeAcc.Address()))
	c.Check(acc.Type, Equals, model.TxSourceAccTypeTrade)
	c.Check(acc.LockExpiresAt, WithinDuration, time.Now(), lockDuration)
	c.Check(acc.LockUnlockedAt.Before(time.Now()), Equals, true)
}

func (s *TxSourceAccsSuite) createNewLock(keyPair keypair.Full, tradeID, userID string, lockType model.TXSourceAccType) (*model.ParsedTXSourceAcc, error) {
	return s.locker.Create(s.ctx, keyPair, tradeID, userID, lockType, lockDuration)
}

func (s *TxSourceAccsSuite) deleteAllPoolLocks() error {
	query := fmt.Sprintf(`
for i in %s
filter i.type == "%s"
remove i._key in %s
`, dbconst.ColTxSourceAccs, model.TxSourceAccTypePool, dbconst.ColTxSourceAccs)
	return dal.DBExec(s.ctx, query, map[string]interface{}{}, s.db)
}

func (s *TxSourceAccsSuite) findPoolSourceAccs() ([]model.TXSourceAcc, error) {
	var poolAccs []model.TXSourceAcc
	err := dal.DBQueryMany(s.ctx, &poolAccs, `
for i in tx_source_accounts
filter i.type == "pool"
return i
`, map[string]interface{}{}, s.db)
	return poolAccs, err
}

func findLocksOfUser(poolAccs []model.TXSourceAcc, userID string) (int, *model.TXSourceAcc) {
	foundCount := 0
	var found model.TXSourceAcc
	for i := 0; i < len(poolAccs); i++ {
		if poolAccs[i].LockUserID == userID {
			foundCount++
			found = poolAccs[i]
		}
	}
	return foundCount, &found
}

func (s *TxSourceAccsSuite) SetUpTest(c *C) {
	err := s.deleteAllPoolLocks()
	if err != nil {
		c.Assert(err, Not(Contains), "not found")
	}
	s.tradeID = "trade_" + strconv.Itoa(time.Now().Nanosecond())
	s.tradeAcc, err = s.createUnlockedTradeLock("Test", s.tradeID)
	c.Assert(err, IsNil)
}

func (s *TxSourceAccsSuite) TestAcquireLockReacquireShouldExtend(c *C) {
	userID := "user ID"
	tradeID := "trade ID"
	acquiredSource, err := s.locker.Acquire(s.ctx, model.SCAddr(s.tradeAcc.Address()), tradeID, userID, lockDuration)
	c.Assert(err, IsNil)
	foundSource, err := s.locker.Find(s.ctx, tradeID, userID)
	c.Assert(err, IsNil)
	c.Assert(acquiredSource, DeepEquals, foundSource)
	// Extend
	_, err = s.locker.Acquire(s.ctx, model.SCAddr(s.tradeAcc.Address()), tradeID, userID, lockDuration)
	c.Assert(err, IsNil)
	newSource, err := s.locker.FindByKey(s.ctx, model.SCAddr(s.tradeAcc.Address()))
	c.Assert(err, IsNil)
	c.Assert(newSource.PubKey, Equals, acquiredSource.PubKey)
	c.Assert(newSource, NotNil)
	c.Check(newSource.PubKey, Equals, model.SCAddr(s.tradeAcc.Address()))
	c.Check(newSource.Type, Equals, model.TxSourceAccTypeTrade)
	c.Check(newSource.LockTradeID, Equals, tradeID)
	c.Check(newSource.LockUserID, Equals, userID)
	c.Check(newSource.LockExpiresAt, WithinDuration, time.Now(), lockDuration)
	c.Check(newSource.LockExpiresAt.After(acquiredSource.LockExpiresAt), IsTrue)
	c.Check(newSource.LockUnlockedAt, IsNil)
}

func (s *TxSourceAccsSuite) TestShouldFailOnAcquireNoLock(c *C) {
	_, err := s.locker.Acquire(s.ctx, model.SCAddr("AcquireNoLock"), "AcquireNoLock trade ID", "user ID", lockDuration)
	c.Assert(err, NotNil)
}

func (s *TxSourceAccsSuite) TestUnlock(c *C) {
	userID := "user ID"
	tradeID := "trade ID"
	_, err := s.locker.Acquire(s.ctx, model.SCAddr(s.tradeAcc.Address()), tradeID, userID, lockDuration)
	c.Assert(err, IsNil)
	source, err := s.locker.FindByKey(s.ctx, model.SCAddr(s.tradeAcc.Address()))
	c.Assert(err, IsNil)
	c.Assert(source, NotNil)
	err = source.MustBeValidFor(tradeID, userID)
	c.Assert(err, IsNil)
	err = s.locker.Unlock(s.ctx, tradeID, userID)
	c.Assert(err, IsNil)
	c.Assert(source, NotNil)
	// Should not find the lock
	source, err = s.locker.FindByKey(s.ctx, model.SCAddr(s.tradeAcc.Address()))
	c.Assert(source.MustBeValidFor(tradeID, userID), NotNil)
}

func (s *TxSourceAccsSuite) TestAcquireReacquireInvalidated(c *C) {
	userID := "user ID"
	oldLock, err := s.locker.Acquire(s.ctx, model.SCAddr(s.tradeAcc.Address()), s.tradeID, userID, lockDuration)
	c.Assert(err, IsNil)
	c.Assert(oldLock, NotNil)
	err = s.locker.Unlock(s.ctx, s.tradeID, userID)
	c.Assert(err, IsNil)
	// Reacquire should permit the unlocking
	_, err = s.locker.Acquire(s.ctx, model.SCAddr(s.tradeAcc.Address()), s.tradeID, userID, lockDuration)
	c.Assert(err, IsNil)
	newLock, err := s.locker.FindByKey(s.ctx, model.SCAddr(s.tradeAcc.Address()))
	c.Assert(err, IsNil)
	c.Assert(newLock, NotNil)
	c.Check(newLock.PubKey, Equals, model.SCAddr(s.tradeAcc.Address()))
	c.Check(newLock.Type, Equals, model.TxSourceAccTypeTrade)
	c.Check(newLock.LockUserID, Equals, userID)
	c.Check(newLock.LockExpiresAt, WithinDuration, time.Now(), lockDuration)
	c.Check(newLock.LockExpiresAt.After(oldLock.LockExpiresAt), IsTrue)
	c.Check(newLock.LockUnlockedAt, IsNil)
}

func (s *TxSourceAccsSuite) createUnlockedLock(prefix string, tradeID *string, lockType model.TXSourceAccType) (*keypair.Full, error) {
	userID := "unlock setup user ID"
	key1, err := keypair.Random()
	if err != nil {
		return nil, err
	}
	_, err = s.createNewLock(*key1, s.tradeID, userID, lockType)
	if err != nil {
		return nil, err
	}
	return key1, s.locker.Unlock(s.ctx, s.tradeID, userID)
}

func (s *TxSourceAccsSuite) createUnlockedTradeLock(prefix string, tradeID string) (*keypair.Full, error) {
	return s.createUnlockedLock(prefix, &tradeID, model.TxSourceAccTypeTrade)
}

func (s *TxSourceAccsSuite) createUnlockedPoolLock(prefix string) (*keypair.Full, error) {
	return s.createUnlockedLock(prefix, nil, model.TxSourceAccTypePool)
}

func (s *TxSourceAccsSuite) TestReacquireLock(c *C) {
	firstUserID := "user ID"
	secondUserID := "user ID 2"
	poolAcc, err := s.createUnlockedPoolLock("AcquireLockReacquireShouldExtend")
	c.Assert(err, IsNil)
	olderLock, err := s.locker.Acquire(s.ctx, model.SCAddr(s.tradeAcc.Address()), s.tradeID, firstUserID, lockDuration)
	c.Assert(err, IsNil)
	c.Assert(olderLock, NotNil)
	// Reacquire should allow for same user
	newerLock, err := s.locker.Acquire(s.ctx, model.SCAddr(s.tradeAcc.Address()), s.tradeID, secondUserID, lockDuration)
	c.Check(newerLock.LockExpiresAt.After(olderLock.LockExpiresAt), IsTrue)
	c.Assert(err, IsNil)
	// Trade account should belong to first user
	tradeLock, err := s.locker.FindByKey(s.ctx, model.SCAddr(s.tradeAcc.Address()))
	c.Assert(err, IsNil)
	c.Check(tradeLock.PubKey, Equals, model.SCAddr(s.tradeAcc.Address()))
	c.Check(tradeLock.Type, Equals, model.TxSourceAccTypeTrade)
	c.Check(tradeLock.LockUserID, Equals, firstUserID)
	c.Check(tradeLock.LockExpiresAt, WithinDuration, time.Now(), lockDuration)
	c.Check(tradeLock.LockUnlockedAt, IsNil)
	// Pool account should belong to second user
	poolLock, err := s.locker.FindByKey(s.ctx, model.SCAddr(poolAcc.Address()))
	c.Assert(err, IsNil)
	c.Assert(poolLock, NotNil)
	c.Check(poolLock.PubKey, Equals, model.SCAddr(poolAcc.Address()))
	c.Check(poolLock.Type, Equals, model.TxSourceAccTypePool)
	c.Check(poolLock.LockUserID, Equals, secondUserID)
	c.Check(poolLock.LockExpiresAt, WithinDuration, time.Now(), lockDuration)
	c.Check(poolLock.LockUnlockedAt, IsNil)
	// There are no more locks. Third user should get an error
	_, err = s.locker.Acquire(s.ctx, model.SCAddr(s.tradeAcc.Address()), s.tradeID, "third user ID", lockDuration)
	c.Assert(err, NotNil)
	// There are no more locks. User should get an error for a different trade
	_, err = s.locker.Acquire(s.ctx, model.SCAddr(s.tradeAcc.Address()), "Other trade ID", "third user ID", lockDuration)
	c.Assert(err, NotNil)
	// There are no more locks. Even the users that already own locks should get an error for a different trade
	_, err = s.locker.Acquire(s.ctx, model.SCAddr(s.tradeAcc.Address()), "Other trade ID", firstUserID, lockDuration)
	c.Assert(err, NotNil)
	// Reacquiring an invalidated lock should allow for anyone
	err = s.locker.Unlock(s.ctx, s.tradeID, firstUserID)
	c.Assert(err, IsNil)
	_, err = s.locker.Acquire(s.ctx, model.SCAddr(s.tradeAcc.Address()), "Other trade ID", "third user ID", lockDuration)
	c.Assert(err, IsNil)
}
