package dal

import (
	"context"
	"testing"

	dbs "bitbucket.org/cerealia/apps/go-lib/setup/arangodb"
	driver "github.com/arangodb/go-driver"
	. "github.com/robert-zaremba/checkers"
	"github.com/robert-zaremba/flag"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

var testctx = context.Background()

func init() {
	_ = flag.String("check.f", "", "Testing selector flag")
	flag.Parse()
}

type DalSuite struct {
	db driver.Database
}

var _ = Suite(&DalSuite{})

func (s *DalSuite) SetUpSuite(c *C) {
	var err error
	s.db, err = dbs.GetDb(testctx)
	c.Assert(err, IsNil, Comment("Failed to connect db"))
}
