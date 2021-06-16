// Package resolver_tests contains integration tests without REST-based functions
package resolvertests

import (
	"context"
	"testing"
	"time"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/resolver"
	"bitbucket.org/cerealia/apps/go-lib/resolver/testutil"
	"bitbucket.org/cerealia/apps/go-lib/setup/arangodb"
	"bitbucket.org/cerealia/apps/go-lib/stellar"
	"bitbucket.org/cerealia/apps/go-lib/stellar/txsource"
	"bitbucket.org/cerealia/apps/go-lib/stellar/txsource/txsourceimpl"
	driver "github.com/arangodb/go-driver"
	"github.com/robert-zaremba/flag"
	"github.com/robert-zaremba/log15"
	. "gopkg.in/check.v1"
)

type TradeIntegrationSuite struct {
	noopDriver                      *stellar.Driver
	testnetDriver                   *stellar.Driver
	noopResolver                    resolver.Resolver
	testnetResolver                 resolver.Resolver
	db                              driver.Database
	txSourceDriver                  txsource.Driver
	trade                           *model.Trade
	buyer, seller, third, moderator *testutil.Credentials
}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

var logger = log15.Root()
var testctx = context.Background()
var _ = Suite(&TradeIntegrationSuite{})

func init() {
	_ = flag.String("check.f", "", "Testing selector flag")
	flag.Parse()
}

func makeResolver(c *C, db driver.Database, driverName string, txSourceDriver txsource.Driver) (*stellar.Driver, resolver.Resolver) {
	driver, err := stellar.NewDriver(driverName)
	c.Assert(err, IsNil)
	return driver, resolver.NewResolver(db, driver, txSourceDriver)
}

func (s *TradeIntegrationSuite) SetUpSuite(c *C) {
	ctx := context.Background()
	db, erre := arangodb.GetDb(ctx)
	c.Assert(erre, IsNil)
	s.db = db
	s.txSourceDriver = txsourceimpl.NewDriver(db, time.Minute*4)
	s.noopDriver, s.noopResolver = makeResolver(c, db, "noop", s.txSourceDriver)
	s.testnetDriver, s.testnetResolver = makeResolver(c, db, "horizon-test", s.txSourceDriver)
}

func (s *TradeIntegrationSuite) SetUpTest(c *C) {
	err := testutil.CleanSourceAccs(testctx, s.db)
	c.Assert(err, IsNil)
	s.buyer, s.seller, s.third, s.moderator, s.trade, err = testutil.CreateTrade(s.noopResolver, &sampleDesc)
	c.Assert(err, IsNil)
}
