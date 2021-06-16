// Package trades contains integration tests of whole app
package trades

import (
	"context"
	"testing"
	"time"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/resolver"
	"bitbucket.org/cerealia/apps/go-lib/resolver/testutil"
	dbs "bitbucket.org/cerealia/apps/go-lib/setup/arangodb"
	"bitbucket.org/cerealia/apps/go-lib/stellar"
	"bitbucket.org/cerealia/apps/go-lib/stellar/txsource"
	"bitbucket.org/cerealia/apps/go-lib/stellar/txsource/txsourceimpl"
	driver "github.com/arangodb/go-driver"
	"github.com/robert-zaremba/flag"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

func init() {
	// Used when running selective test: 'go test -check.f ValidationSuite.*'
	_ = flag.String("check.f", "", "Testing selector flag")
	flag.Parse()
}

type ValidationSuite struct{}

var _ = Suite(&ValidationSuite{})

var _ = Suite(&TradeIntegrationSuite{})

type TradeIntegrationSuite struct {
	noopDriver        *stellar.Driver
	testnetDriver     *stellar.Driver
	noopResolver      resolver.Resolver
	testnetResolver   resolver.Resolver
	db                driver.Database
	noopDocHandler    DocHandler
	testnetDocHandler DocHandler
	txSourceDriver    txsource.Driver

	sampleExpireTimeStr string
	sampleExpireTime    time.Time

	trade                           *model.Trade
	buyer, seller, third, moderator *testutil.Credentials
}

func makeResolver(c *C, db driver.Database, driverName string, txSourceDriver txsource.Driver) (*stellar.Driver, resolver.Resolver) {
	driver, err := stellar.NewDriver(driverName)
	c.Assert(err, IsNil)
	return driver, resolver.NewResolver(db, driver, txSourceDriver)
}

func (s *TradeIntegrationSuite) SetUpSuite(c *C) {
	ctx := context.Background()
	var err error
	db, err := dbs.GetDb(ctx)
	c.Assert(err, IsNil)
	s.db = db
	s.txSourceDriver = txsourceimpl.NewDriver(db, time.Minute*4)
	s.noopDriver, s.noopResolver = makeResolver(c, s.db, "noop", s.txSourceDriver)
	s.testnetDriver, s.testnetResolver = makeResolver(c, s.db, "horizon-test", s.txSourceDriver)
	s.noopDocHandler = DocHandler{s.noopDriver, s.txSourceDriver}
	s.testnetDocHandler = DocHandler{s.testnetDriver, s.txSourceDriver}
	s.sampleExpireTimeStr = "2020-12-31T00:00:00+00:00"
	s.sampleExpireTime, err = time.Parse(time.RFC3339, s.sampleExpireTimeStr)
	c.Assert(err, IsNil)
}

func (s *TradeIntegrationSuite) SetUpTest(c *C) {
	ctx := context.Background()
	err := testutil.CleanSourceAccs(ctx, s.db)
	c.Assert(err, IsNil)
	s.buyer, s.seller, s.third, s.moderator, s.trade, err = testutil.CreateTrade(s.noopResolver, &sampleDesc)
	c.Assert(err, IsNil)
}
